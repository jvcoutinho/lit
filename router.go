package lit

import (
	"net/http"
	"time"

	"github.com/jvcoutinho/lit/internal/routes"
)

const DefaultReadHeaderTimeout = 3 * time.Second

// Result of an HTTP request.
//
// See the lit/render package.
type Result interface {
	// Render writes this into the HTTP response managed by ctx.
	Render(ctx *Context) error
}

// HandleFunc is a function that handles requests.
type HandleFunc func(ctx *Context) Result

// Router manages API routes.
//
// It is the entrypoint of a Lit-based application.
type Router struct {
	trie     *routes.Trie
	handlers map[*routes.Node]HandleFunc
	server   *http.Server
}

// NewRouter creates a new Router instance.
func NewRouter() *Router {
	return &Router{
		trie:     routes.NewTrie(),
		handlers: make(map[*routes.Node]HandleFunc),
		server: &http.Server{
			ReadHeaderTimeout: DefaultReadHeaderTimeout,
		},
	}
}

// Handle registers the handler for the given pattern and method.
// If a handler already exists for pattern, Handle panics.
func (r *Router) Handle(pattern string, method string, handle HandleFunc) {
	if handle == nil {
		panic("handle should not be nil")
	}

	node, err := r.trie.Insert(pattern, method)
	if err != nil {
		panic(err)
	}

	r.handlers[node] = handle
}

// ServeHTTP dispatches the request to the handler whose pattern and method most closely matches one previously defined.
func (r *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	node, parameters := r.trie.Match(request.URL.Path, request.Method)
	if node == nil {
		http.NotFound(writer, request)

		return
	}

	context := NewContext(writer, request)
	context.setArguments(parameters)

	handle := r.handlers[node]

	result := handle(context)
	result.Render(context)
}

// Server this router uses for listening and serving requests.
//
// By default, it has ReadHeaderTimeout = DefaultReadHeaderTimeout, but one can change it at will.
func (r *Router) Server() *http.Server {
	return r.server
}

// ListenAndServe listens on the TCP network address addr and then handle requests on incoming connections.
func (r *Router) ListenAndServe(addr string) error {
	r.server.Addr = addr
	r.server.Handler = r

	return r.server.ListenAndServe()
}
