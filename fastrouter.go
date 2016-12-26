package fastrouter

import (
	"fmt"
	"net/http"
	"path"
	"strings"
)

const (
	DELETE  = "DELETE"
	GET     = "GET"
	HEAD    = "HEAD"
	OPTIONS = "OPTIONS"
	PATCH   = "PATCH"
	POST    = "POST"
	PUT     = "PUT"
)

// Router represents the router for the request.
type Router struct {
	trees           []*radixTree
	methods         []string
	NotFoundHandler http.Handler
}

// NewRouter returns a new Router pointer.
func NewRouter(notFoundHandler http.Handler) *Router {
	if notFoundHandler == nil {
		notFoundHandler = http.NotFoundHandler()
	}

	return &Router{
		trees:           make([]*radixTree, 0),
		NotFoundHandler: notFoundHandler,
	}
}

// register add a handler to the router according the method and path.
func (router *Router) register(method, path string, handler http.HandlerFunc) {
	method = strings.ToUpper(method)

	for i, tree := range router.trees {
		if method == router.methods[i] {
			tree.insert(path, handler)
			return
		}
	}

	router.trees = append(router.trees, newRadixTree())
	router.methods = append(router.methods, method)

	tree := router.trees[len(router.trees)-1]
	tree.insert(path, handler)
}

// Handle insert path to some tree according to method.
func (router *Router) Handle(
	path string, handler http.HandlerFunc, methods ...string) {

	for _, mehtod := range methods {
		router.register(mehtod, path, handler)
	}
}

// DELETE is a wrapper of Handle with DELETE method.
func (router *Router) DELETE(path string, handler http.HandlerFunc) {
	router.register(DELETE, path, handler)
}

// GET is a wrapper of Handle with GET method.
func (router *Router) GET(path string, handler http.HandlerFunc) {
	router.register(GET, path, handler)
}

// HEAD is a wrapper of Handle with HEAD method.
func (router *Router) HEAD(path string, handler http.HandlerFunc) {
	router.register(HEAD, path, handler)
}

// OPTIONS is a wrapper of Handle with OPTIONS method.
func (router *Router) OPTIONS(path string, handler http.HandlerFunc) {
	router.register(OPTIONS, path, handler)
}

// PATCH is a wrapper of Handle with PATCH method.
func (router *Router) PATCH(path string, handler http.HandlerFunc) {
	router.register(PATCH, path, handler)
}

// POST is a wrapper of Handle with POST method.
func (router *Router) POST(path string, handler http.HandlerFunc) {
	router.register(POST, path, handler)
}

// PUT is a wrapper of Handle with PUT method.
func (router *Router) PUT(path string, handler http.HandlerFunc) {
	router.register(PUT, path, handler)
}

// ServeHTTP implements interface http.Handler.
func (router *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var handler http.HandlerFunc

	n := len(router.methods)
	for i := 0; i < n; i++ {
		if router.methods[i] == req.Method {
			handler = router.trees[i].handler(req)
			break
		}
	}

	if handler == nil {
		router.NotFoundHandler.ServeHTTP(w, req)
		return
	}

	handler(w, req)
}

// insertSubtreePath traverses the radix tree in dfs, finds all paths
// and then re-inserts them.
func (router *Router) insertSubtreePath(
	node *radixTreeNode, method, path string) {

	if node == nil {
		return
	}

	if node.isParam {
		if node.regex != nil {
			path += fmt.Sprintf(
				"%s%s:%s%s", OpenTag, node.chunk, node.regex, CloseTag,
			)
		} else {
			path += fmt.Sprintf("%s%s%s", OpenTag, node.chunk, CloseTag)
		}
	} else {
		path += node.chunk
	}

	if node.handler != nil {
		router.register(method, path, node.handler)
	}

	for _, child := range node.children {
		router.insertSubtreePath(child, method, path)
	}
}

// SubRouter adds a prefixed router.
func (router *Router) SubRouter(prefix string, subrouter *Router) {
	prefix = path.Clean(strings.TrimRight(prefix, "/"))
	if prefix == "" {
		return
	}

	if prefix[0] != Slash {
		prefix = fmt.Sprintf("/%s", prefix)
	}

	for i, tree := range subrouter.trees {
		router.insertSubtreePath(tree.root, subrouter.methods[i], prefix)
	}
}

// NOTE: THIS IS ONLY FOR DEBUG.
func (router *Router) traverse() {
	for _, tree := range router.trees {
		tree.traverse()
	}
}
