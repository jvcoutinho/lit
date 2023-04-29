package routes

import (
	"github.com/jvcoutinho/lit/internal/maps"
	"github.com/jvcoutinho/lit/internal/slices"
)

// Graph stores route definitions.
type Graph map[string][]string

// CanBeInserted checks if route can be defined in this graph.
//
// If it can't be inserted (ok equals false), CanBeInserted returns the reason error.
func (g Graph) CanBeInserted(route Route) (reason error, ok bool) {
	patternPaths := route.Path()

	if duplicate, has := hasDuplicateArguments(patternPaths); has {
		return ErrDuplicateArguments{duplicate}, false
	}

	if !maps.ContainsKey(g, route.Method) {
		return nil, true
	}

	previousNode := route.Method
	for _, path := range patternPaths {
		adjacentNodes := g[previousNode]

		if isArgument(path) && slices.Any(adjacentNodes, isArgument) {
			path, _ = slices.First(adjacentNodes, isArgument)
		}

		if !maps.ContainsKey(g, path) || !slices.Contains(adjacentNodes, path) {
			return nil, true
		}

		previousNode = path
	}

	return ErrRouteAlreadyDefined{route}, false
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
