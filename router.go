package lit

import (
	"net/http"

	"github.com/jvcoutinho/lit/internal/routes"
)

// HandleFunc is a function that handles requests.
type HandleFunc func(ctx *Context)

// Router manages API routes.
//
// It is the entrypoint of a Lit-based application.
type Router struct {
	graph    routes.Graph
	handlers map[routes.Route]HandleFunc
}

// NewRouter creates a new Router instance.
func NewRouter() *Router {
	return &Router{
		routes.NewGraph(),
		make(map[routes.Route]HandleFunc),
	}
}

// Handle registers the handler for the given pattern and method.
// If a handler already exists for pattern, Handle panics.
func (r *Router) Handle(pattern string, method string, handler HandleFunc) {
	route := routes.NewRoute(pattern, method)

	if err, ok := r.graph.CanBeInserted(route); !ok {
		panic(err)
	}

	r.graph.Add(route)
	r.handlers[route] = handler
}

// ServeHTTP dispatches the request to the handler whose pattern and method most closely matches one previously defined.
func (r *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	route := routes.NewRoute(request.URL.Path, request.Method)

	match, ok := r.graph.MatchRoute(route)
	if !ok {
		http.NotFound(writer, request)
		return
	}

	handler := r.handlers[match.MatchedRoute()]

	ctx := newContext(writer, request)
	ctx.setArguments(match.Parameters)

	handler(ctx)
}

// ListenAndServe listens on the TCP network address addr and then handle requests on incoming connections.
func (r *Router) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, r)
}
