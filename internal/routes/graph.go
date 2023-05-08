package routes

import (
	"github.com/jvcoutinho/lambda/maps"
	"github.com/jvcoutinho/lambda/slices"
)

const terminalNode = "/"

// Graph stores route definitions.
type Graph map[string][]string

// NewGraph creates a new Graph instance.
func NewGraph() Graph {
	return make(map[string][]string)
}

// CanBeInserted checks if route can be defined in this graph.
//
// If it can't be inserted (ok equals false), CanBeInserted returns the reason error.
func (g Graph) CanBeInserted(route Route) (ok bool, reason error) {
	routePath := append(route.Path(), terminalNode)

	if duplicate, has := hasDuplicateArguments(routePath); has {
		return false, DuplicateArgumentsError{duplicate}
	}

	if !maps.ContainsKey(g, route.Method) {
		return true, nil
	}

	previousNode := route.Method
	for _, path := range routePath {
		adjacentNodes := g[previousNode]

		if isArgument(path) && slices.Any(adjacentNodes, isArgument) {
			path, _ = slices.FirstBy(adjacentNodes, isArgument)
		}

		if !maps.ContainsKey(g, path) || !slices.Contains(adjacentNodes, path) {
			return true, nil
		}

		previousNode = path
	}

	return false, RouteAlreadyDefinedError{route}
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

// MatchRoute finds a correspondence between route and a pre-defined one (in this graph).
//
// If a correspondence can not be found, ok equals false.
func (g Graph) MatchRoute(route Route) (Match, bool) {
	if !maps.ContainsKey(g, route.Method) {
		return Match{}, false
	}

	match := newMatch()
	match.addMethod(route.Method)

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
		if !g.matchAdjacentNode(child, path, pathIndex+1, match) {
			return false
		}

		match.addPathFragmentAtBeginning(child)

		return true
	}

	childrenArguments := slices.Filter(children, isArgument)
	if len(childrenArguments) == 0 {
		return false
	}

	for i := 0; i < len(childrenArguments); i++ {
		if !g.matchAdjacentNode(childrenArguments[i], path, pathIndex+1, match) {
			continue
		}

		match.addPathArgumentAtBeginning(childrenArguments[i], child)

		return true
	}

	return false
}
