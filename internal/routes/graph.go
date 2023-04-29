package routes

import (
	"github.com/jvcoutinho/lit/internal/maps"
	"github.com/jvcoutinho/lit/internal/slices"
)

// Graph stores route definitions.
type Graph map[Node][]Node

// Exists check if route is defined in this graph.
func (g Graph) Exists(route Route) bool {
	patternPaths := route.Path()
	methodNode := Node(route.Method)

	if !maps.ContainsKey(g, methodNode) {
		return false
	}

	previousNode := methodNode
	for _, path := range patternPaths {
		pathNode := Node(path)
		adjacentNodes := g[previousNode]

		if pathNode.IsArgument() && slices.Any(adjacentNodes, Node.IsArgument) {
			pathNode, _ = slices.First(adjacentNodes, Node.IsArgument)
		}

		if !maps.ContainsKey(g, pathNode) || !slices.Contains(adjacentNodes, pathNode) {
			return false
		}

		previousNode = pathNode
	}

	return true
}

// Add adds the route to this graph.
func (g Graph) Add(route Route) {
	patternPaths := route.Path()
	methodNode := Node(route.Method)

	if !maps.ContainsKey(g, methodNode) {
		g[methodNode] = make([]Node, 0)
	}

	previousNode := methodNode
	for _, path := range patternPaths {
		pathNode := Node(path)

		if !maps.ContainsKey(g, pathNode) {
			g[pathNode] = make([]Node, 0)
		}

		g[previousNode] = append(g[previousNode], pathNode)
		previousNode = pathNode
	}
}
