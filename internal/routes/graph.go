package routes

import (
	"github.com/jvcoutinho/lit/internal/maps"
	"github.com/jvcoutinho/lit/internal/slices"
)

const terminalNode = "/"

// Graph stores route definitions.
type Graph map[string][]string

// CanBeInserted checks if route can be defined in this graph.
//
// If it can't be inserted (ok equals false), CanBeInserted returns the reason error.
func (g Graph) CanBeInserted(route Route) (reason error, ok bool) {
	routePath := append(route.Path(), terminalNode)

	if duplicate, has := hasDuplicateArguments(routePath); has {
		return ErrDuplicateArguments{duplicate}, false
	}

	if !maps.ContainsKey(g, route.Method) {
		return nil, true
	}

	previousNode := route.Method
	for _, path := range routePath {
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
	if !maps.ContainsKey(g, route.Method) {
		g[route.Method] = make([]string, 0)
	}

	previousNode := route.Method
	routePath := append(route.Path(), terminalNode)

	for _, path := range routePath {
		if !maps.ContainsKey(g, path) {
			g[path] = make([]string, 0)
		}

		g[previousNode] = append(g[previousNode], path)
		previousNode = path
	}
}

func (g Graph) Match(route Route) (Match, bool) {
	if !maps.ContainsKey(g, route.Method) {
		return Match{}, false
	}

	match := NewMatch()
	match.AddMethod(route.Method)

	previousNode := route.Method
	routePath := route.Path()

	for _, path := range routePath {
		adjacentNodes := g[previousNode]

		if slices.Contains(adjacentNodes, path) {
			match.AddPathFragment(path)
			previousNode = path

			continue
		}

		arguments := slices.Filter(adjacentNodes, isArgument)
		if len(arguments) == 0 {
			return Match{}, false
		}

		match.AddPathArgument(arguments[0], path)
		previousNode = arguments[0]
	}

	if slices.Contains(g[previousNode], terminalNode) {
		return *match, true
	}

	return Match{}, false
}

func (g Graph) matchArgumentNode(parent string, adjacent string, match *Match) bool {
	adjacentNodes := g[parent]

	if slices.Contains(adjacentNodes, adjacent) {
		match.AddPathFragment(adjacent)

	}

	return true
}
