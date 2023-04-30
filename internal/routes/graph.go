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

		if !slices.Contains(g[previousNode], path) {
			g[previousNode] = append(g[previousNode], path)
		}

		previousNode = path
	}
}

func (g Graph) Match(route Route) (Match, bool) {
	if !maps.ContainsKey(g, route.Method) {
		return Match{}, false
	}

	match := NewMatch()
	match.AddMethod(route.Method)

	if !g.matchAdjacentNode(route.Method, route.Path(), 0, match) {
		return Match{}, false
	}

	return *match, true
}

func (g Graph) matchAdjacentNode(parent string, path []string, pathIndex int, match *Match) bool {
	children := g[parent]

	if pathIndex == len(path) {
		return slices.Contains(children, terminalNode)
	}

	child := path[pathIndex]

	if slices.Contains(children, child) {
		if g.matchAdjacentNode(child, path, pathIndex+1, match) {
			match.AddPathFragment(child)
			return true
		}

		return false
	}

	childrenArguments := slices.Filter(children, isArgument)
	if len(childrenArguments) == 0 {
		return false
	}

	for i := 0; i < len(childrenArguments); i++ {
		if !g.matchAdjacentNode(childrenArguments[i], path, pathIndex+1, match) {
			continue
		}

		match.AddPathArgument(childrenArguments[i], child)
	}

	return true
}
