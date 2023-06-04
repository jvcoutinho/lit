package lit

import (
	"net/http"

	"github.com/jvcoutinho/lit/internal/trie"
)

// HandlerFunc is the standard HTTP handler function in Lit ecossystem.
type HandlerFunc func(*Request) Response

// Router manages, listens and serves HTTP requests.
type Router struct {
	trie     *trie.Trie
	handlers map[*trie.Node]HandlerFunc
}

// NewRouter creates a new Router instance.
func NewRouter() *Router {
	return &Router{
		trie.New(),
		make(map[*trie.Node]HandlerFunc),
	}
}

// Handle registers the handler for the given pattern and method.
//
// If the route can't be registered, Handle panics.
func (r *Router) Handle(pattern string, method string, handler HandlerFunc) {
	if handler == nil {
		panic("handler should not be nil")
	}

	handlerNode, err := r.trie.Insert(pattern, method)
	if err != nil {
		panic(err.Error())
	}

	r.handlers[handlerNode] = handler
}

// ServeHTTP dispatches the request to the handler whose pattern most closely matches the request URL
// and whose method is the same as the request method.
func (r *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
}
