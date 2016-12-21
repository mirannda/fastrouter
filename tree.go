package fastrouter

import (
	"bytes"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"unsafe"
)

const Threshold = 8

const (
	OpenTag  = '{'
	CloseTag = '}'
	Slash    = '/'
	Colon    = ':'
	And      = '&'
	Equal    = '='
	Dot      = '.'
)

/************************************************************************
* radixTreeNode
*************************************************************************/

type radixTreeNode struct {
	chunk    string
	isParam  bool
	regex    *regexp.Regexp
	handler  http.HandlerFunc
	indices  string
	children []*radixTreeNode
}

func parseParam(path string) (name, value, newPath string) {
	// path[0] must be OpenTag.
	for i := 0; i < len(path); i++ {
		if path[i] == CloseTag && (i == len(path)-1 || path[i+1] == Slash) {
			j := strings.IndexByte(path[:i], Colon)
			if j == -1 {
				name = path[1:i]
			} else {
				name = path[1:j]
				value = path[j+1 : i]
			}
			newPath = path[i+1:]
			return
		}
	}
	panic("tag not closed in " + path)
}

func (node *radixTreeNode) childIndex(target byte) int {
	if len(node.indices) < Threshold {
		return strings.IndexByte(node.indices, target)
	}

	var mid int

	low, high, indices := 0, len(node.indices)-1, node.indices
	for low <= high {
		mid = low + ((high - low) >> 1)
		if indices[mid] < target {
			low = mid + 1
		} else if indices[mid] > target {
			high = mid - 1
		} else {
			return mid
		}
	}
	return -1
}

// add a none param child
func (node *radixTreeNode) addChild(child *radixTreeNode) {
	var i int

	for i = 0; i < len(node.indices); i++ {
		if node.indices[i] > child.chunk[0] {
			break
		}
	}

	node.indices = node.indices[:i] + string(child.chunk[0]) + node.indices[i:]
	node.children = append(
		node.children[:i],
		append([]*radixTreeNode{child}, node.children[i:]...)...,
	)
}

func (node *radixTreeNode) insert(
	path string, handler http.HandlerFunc) *radixTreeNode {

	var index int

	// new path
	if node == nil {
		index = strings.IndexByte(path, OpenTag)

		// no param in the path
		if index == -1 {
			return &radixTreeNode{
				chunk:    path,
				handler:  handler,
				children: make([]*radixTreeNode, 0),
			}
		}

		// deal with the part which is not param
		if index != 0 {
			node = &radixTreeNode{
				chunk:    path[:index],
				children: make([]*radixTreeNode, 1),
			}
			node.children[0] = (*radixTreeNode)(nil).insert(
				path[index:], handler,
			)
			return node
		}

		// parse the param, path[0] is OpenTag.
		name, value, path := parseParam(path)

		node = &radixTreeNode{
			chunk:    name,
			isParam:  true,
			children: make([]*radixTreeNode, 0),
		}

		if value != "" {
			regex, err := regexp.Compile(value)
			if err != nil {
				panic(err)
			}
			node.regex = regex
		}

		// inserting finished
		if path != "" {
			node.indices = string(path[0])
			node.children = append(
				node.children, (*radixTreeNode)(nil).insert(path, handler),
			)
		} else {
			node.handler = handler
		}

		return node
	}

	// node is not nil and node is param
	if node.isParam {
		if path[0] != OpenTag {
			panic(path + " conflicts with " + node.chunk)
		}

		// validate the param
		name, value, newPath := parseParam(path)
		if name != node.chunk ||
			node.regex != nil && value == "" ||
			node.regex == nil && value != "" ||
			node.regex != nil && value != "" && node.regex.String() != value {
			panic(path + ": conflicts with existed path " + node.regex.String())
		}

		// inserting finished
		if newPath == "" {
			if node.handler != nil {
				panic("overwrite handler: " + path)
			}
			node.handler = handler
			return node
		}

		path = newPath

		// deal with the left path
		j := strings.IndexByte(node.indices, path[0])
		if j != -1 {
			node.children[j] = node.children[j].insert(path, handler)
		} else {
			node.addChild((*radixTreeNode)(nil).insert(path, handler))
		}
		return node
	}

	// node is not nil and node is common path
	if path[0] == OpenTag {
		panic(path + " conflicts with " + node.chunk)
	}

	i, j, m, n := 0, 0, len(node.chunk), len(path)
	for i < n && j < m && path[i] == node.chunk[j] {
		i++
		j++
	}

	// need to split node.chunk
	if j < m {
		child := &radixTreeNode{
			chunk:    node.chunk[j:],
			handler:  node.handler,
			indices:  node.indices,
			children: node.children,
		}

		node.chunk = node.chunk[:j]
		node.handler = nil
		node.indices = string(child.chunk[0])
		node.children = []*radixTreeNode{child}
	}

	path = path[i:]

	// inserting finished
	if path == "" {
		node.handler = handler
		return node
	}

	// deal with the left path
	var child *radixTreeNode

	j = node.childIndex(path[0])
	if j != -1 {
		child = node.children[j]
	}
	child = child.insert(path, handler)

	// since node already has a none param child
	if child.isParam {
		panic(path + " conflicts with existed paths")
	}

	if j == -1 {
		node.addChild(child)
	}

	return node
}

/************************************************************************
* radixTree
*************************************************************************/

const N = 1000000

type radixTree struct {
	root     *radixTreeNode
	syncPool *sync.Pool
	chanPool chan *bytes.Buffer
}

func newRadixTree() *radixTree {
	tree := &radixTree{
		syncPool: &sync.Pool{},
		chanPool: make(chan *bytes.Buffer, N),
	}

	for i := 0; i < N; i++ {
		tree.chanPool <- new(bytes.Buffer)
	}

	return tree
}

func (tree *radixTree) insert(path string, handler http.HandlerFunc) {
	tree.root = tree.root.insert(path, handler)
}

func (tree *radixTree) handler(req *http.Request) (handler http.HandlerFunc) {
	var buffer *bytes.Buffer

	v := tree.syncPool.Get()
	if v == nil {
		buffer = new(bytes.Buffer)
	} else {
		buffer = v.(*bytes.Buffer)
	}

	handler = tree.search(req, buffer)
	tree.syncPool.Put(buffer)

	return
}

type puppetBuffer struct {
	buf []byte
	off int
}

// for worker pool
func (tree *radixTree) searchWorker (p *packet, buffer *bytes.Buffer) {
	handler := tree.search(p.req, buffer)
	if handler == nil {
		http.NotFoundHandler().ServeHTTP(*p.w, p.req)
		return
	}
	handler.ServeHTTP(*p.w, p.req)
}

func (tree *radixTree) search(
	req *http.Request, buffer *bytes.Buffer) http.HandlerFunc {

	var (
		i, j, start, length int
		value               string
		path                = cleanPath(req.URL.Path, buffer)
		node                = tree.root
	)
	buffer.Reset()

	n := len(path)
	for i < n {
		if node.isParam {
			start = i

			for i < n && path[i] != Slash {
				i++
			}

			value = path[start:i]
			if node.regex != nil && !node.regex.MatchString(value) {
				log.Println(node.regex, value)
				return nil
			}

			buffer.WriteByte(And)
			buffer.WriteString(node.chunk)
			buffer.WriteByte(Equal)
			buffer.WriteString(value)

			if i == n {
				break
			}
		} else {
			length = len(node.chunk)

			if n-i < length || path[i:i+length] != node.chunk {
				return nil
			}

			i += length
			if i == n {
				break
			}

			switch len(node.children) {
			case 0:
				return nil
			case 1:
				node = node.children[0]
				continue
			}
		}

		j = node.childIndex(path[i])
		if j == -1 {
			return nil
		}
		node = node.children[j]
	}

	if node.handler != nil && buffer.Len() > 0 {
		b := (*puppetBuffer)(unsafe.Pointer(buffer))
		bs := b.buf[b.off:]
		req.URL.RawQuery += *(*string)(unsafe.Pointer(&bs))
	}
	return node.handler
}

// traverse uses BFS to traverse the tree.
// NOTE: THIS IS ONLY FOR DEBUG
func (tree *radixTree) traverse() {
	if tree == nil {
		log.Println("<nil>")
		return
	}
	log.Println(*(tree.root), "\n")

	var node *radixTreeNode
	queue := []*radixTreeNode{tree.root}

	for len(queue) > 0 {
		length := len(queue)
		for i := 0; i < length; i++ {
			node = queue[0]
			queue = queue[1:]

			for _, child := range node.children {
				log.Println(*child)
				queue = append(queue, child)
			}
		}
		log.Println()
	}
}