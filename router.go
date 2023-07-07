package lit

import (
	"errors"
	"net/http"
	"strings"

	"github.com/jvcoutinho/lit/internal/trie"
)

var (
	ErrNilHandler                   = errors.New("handler should not be nil")
	ErrMethodIsEmpty                = errors.New("method should not be empty")
	ErrPatternDoesNotStartWithSlash = errors.New("pattern should start with a slash (/)")
	ErrPatternContainsDoubleSlash   = errors.New("pattern should not contain double slashes (//)")
)

// HandlerFunc is the standard HTTP handler function in Lit ecosystem.
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
		panic(ErrNilHandler)
	}

	if method == "" {
		panic(ErrMethodIsEmpty)
	}

	if !strings.HasPrefix(pattern, "/") {
		panic(ErrPatternDoesNotStartWithSlash)
	}

	if strings.Contains(pattern, "//") {
		panic(ErrPatternContainsDoubleSlash)
	}

	handlerNode, err := r.trie.Insert(pattern, method)
	if err != nil {
		panic(err)
	}

	r.handlers[handlerNode] = handler
}

// ServeHTTP dispatches the request to the handler whose pattern most closely matches the request URL
// and whose method is the same as the request method.
func (r *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	handlerNode, arguments, err := r.trie.Match(request.URL.Path, request.Method)
	if err != nil {
		http.NotFound(writer, request)
		return
	}

	handler := r.handlers[handlerNode]

	req := newRequest(request)
	req.setURLArguments(arguments)

	res := handler(req)
	res.Write(writer)
}
