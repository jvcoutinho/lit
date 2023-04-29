package routes

import (
	"github.com/jvcoutinho/lit/internal/maps"
	"github.com/jvcoutinho/lit/internal/slices"
)

// Graph stores route definitions.
type Graph map[string][]string

// Exists check if route is defined in this graph.
func (g Graph) Exists(route Route) bool {
	patternPaths := route.Path()

	if !maps.ContainsKey(g, route.Method) {
		return false
	}

	previousNode := route.Method
	for _, path := range patternPaths {
		adjacentNodes := g[previousNode]

		if isArgument(path) && slices.Any(adjacentNodes, isArgument) {
			path, _ = slices.First(adjacentNodes, isArgument)
		}

		if !maps.ContainsKey(g, path) || !slices.Contains(adjacentNodes, path) {
			return false
		}

		previousNode = path
	}

	return true
}

// Add adds the route to this graph.
func (g Graph) Add(route Route) {
	patternPaths := route.Path()

	if !maps.ContainsKey(g, route.Method) {
		g[route.Method] = make([]string, 0)
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
