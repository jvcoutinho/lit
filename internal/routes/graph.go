package routes

import (
	"strings"

	"github.com/jvcoutinho/lit/internal/slices"
)

type Graph map[string][]string

func (g Graph) Exists(pattern, method string) bool {
	patternPaths := strings.Split(pattern, "/")

	currentNode := method
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

func (g Graph) Add(pattern, method string) {
	patternPaths := strings.Split(pattern, "/")

	currentNode := method
	for _, path := range patternPaths {
		_, ok := g[currentNode]
		if !ok {
			g[currentNode] = make([]string, 0)
		}

		g[currentNode] = append(g[currentNode], path)
		currentNode = path
	}
}
