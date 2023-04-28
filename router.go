package lit

import (
	"fmt"

	"github.com/jvcoutinho/lit/internal/routes"
)

// Router manages API routes.
//
// It is the entrypoint of a Lit-based application.
type Router struct {
	graph    routes.Graph
	handlers map[string]func(*Context)
}

// NewRouter creates a new Router instance.
func NewRouter() *Router {
	return &Router{
		make(routes.Graph),
		make(map[string]func(*Context)),
	}
}

// Handle registers the handler for the given pattern and method.
// If a handler already exists for pattern, Handle panics.
func (r *Router) Handle(pattern string, method string, handler func(*Context)) {
	route := routes.NewRoute(pattern, method)

	if r.graph.Exists(route) {
		panic(fmt.Sprintf("%s has been already defined", route))
	}

	r.graph.Add(route)
	r.handlers[route.String()] = handler
}
