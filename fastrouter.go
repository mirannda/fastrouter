package fastrouter

import (
	"net/http"
	"strings"
	"fmt"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	PATCH  = "PATCH"
	DELETE = "DELETE"
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

// GET is a wrapper of Handle with GET method.
func (router *Router) GET(path string, handler http.HandlerFunc) {
	router.register(GET, path, handler)
}

// POST is a wrapper of Handle with POST method.
func (router *Router) POST(path string, handler http.HandlerFunc) {
	router.register(POST, path, handler)
}

// PUT is a wrapper of Handle with PUT method.
func (router *Router) PUT(path string, handler http.HandlerFunc) {
	router.register(PUT, path, handler)
}

// PATCH is a wrapper of Handle with PATCH method.
func (router *Router) PATCH(path string, handler http.HandlerFunc) {
	router.register(PATCH, path, handler)
}

// DELETE is a wrapper of Handle with DELETE method.
func (router *Router) DELETE(path string, handler http.HandlerFunc) {
	router.register(DELETE, path, handler)
}

// ServeHTTP implements interface http.Handler.
func (router *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var handler http.HandlerFunc

	for i, method := range router.methods {
		if method == req.Method {
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

func (router *Router) insertSubtreePath(
	node *radixTreeNode, method, path string) {

	if node == nil {
		return
	}

	if node.isParam {
		if node.regex != nil {
			path += fmt.Sprintf("{%s:%s}", node.chunk, node.regex.String())
		} else {
			path += fmt.Sprintf("{%s}", node.chunk)
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

func (router *Router) SubRouter(prefix string, subrouter *Router) {
	for i, tree := range subrouter.trees {
		router.insertSubtreePath(tree.root, subrouter.methods[i], prefix)
	}
}

// NOTE: THIS IS ONLY FOR DEBUG.
func (router *Router) Traverse() {
	for _, tree := range router.trees {
		tree.traverse()
	}
}
