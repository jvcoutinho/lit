package routes

import (
	"strings"

	"github.com/jvcoutinho/lit/internal/maps"
	"github.com/jvcoutinho/lit/internal/slices"
)

// Graph stores route definitions.
type Graph map[string][]string

// Exists check if route is defined in this graph.
func (g Graph) Exists(route Route) bool {
	patternPaths := route.Path()

	currentNode := route.Method
	for _, path := range patternPaths {
		edges, ok := g[currentNode]
		if !ok {
			return false
		}

		if !slices.Any(edges, func(edge string) bool {
			if strings.HasPrefix(edge, ":") && strings.HasPrefix(path, ":") {
				return true
			}

			return path == edge
		}) {
			return false
		}

		currentNode = path
	}

	return true
}

// Add adds the route to this graph.
func (g Graph) Add(route Route) {
	patternPaths := route.Path()

	if !maps.ContainsKey(g, route.Method) {
		g[route.Method] = make([]string, 1)
	}

	previousNode := route.Method
	for _, path := range patternPaths {
		if !maps.ContainsKey(g, path) {
			g[path] = make([]string, 0)
		}

		g[previousNode] = append(g[previousNode], path)
		previousNode = path
	}
}
