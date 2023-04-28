package lit

import (
	"fmt"
	"strings"

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
	pattern = strings.Trim(pattern, "/")
	method = strings.ToUpper(method)

	if r.graph.Exists(pattern, method) {
		panic(fmt.Sprintf("%s /%s has been already defined", method, pattern))
	}

	r.graph.Add(pattern, method)
}
