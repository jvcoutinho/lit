package lit

import (
	"fmt"
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
	handlers map[string]HandleFunc
}

// NewRouter creates a new Router instance.
func NewRouter() *Router {
	return &Router{
		make(routes.Graph),
		make(map[string]HandleFunc),
	}
}

// Handle registers the handler for the given pattern and method.
// If a handler already exists for pattern, Handle panics.
func (r *Router) Handle(pattern string, method string, handler HandleFunc) {
	route := routes.NewRoute(pattern, method)

	if r.graph.Exists(route) {
		panic(fmt.Sprintf("%s has been already defined", route))
	}

	r.graph.Add(route)
	r.handlers[route.String()] = handler
}

// ServeHTTP dispatches the request to the handler whose pattern and method most closely matches one previously defined.
func (r *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	route := routes.NewRoute(request.URL.Path, request.Method)

	ctx := newContext(writer, request)
	handler, ok := r.handlers[route.String()]
	if !ok {
		http.NotFound(writer, request)
		return
	}

	handler(ctx)
}

// Serve listens on the TCP network address addr and then handle requests on incoming connections.
func (r *Router) Serve(addr string) error {
	return http.ListenAndServe(addr, r)
}
