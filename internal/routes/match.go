package routes

import (
	"strings"

	"github.com/jvcoutinho/lit/internal/structures"
)

// Match is a correspondence between an incoming pattern and method to a pre-defined Route.
//
// As the route may have parameters, they are collected throughout the progress of a match.
type Match struct {
	method string
	path   *structures.List[string]

	// The matched arguments from the pattern.
	Parameters map[string]string
}

func newMatch() *Match {
	return &Match{
		path:       structures.NewList[string](),
		Parameters: make(map[string]string),
	}
}

func (m *Match) addMethod(method string) {
	m.method = method
}

func (m *Match) addPathFragmentAtBeginning(fragment string) {
	m.path.InsertAtBeginning(fragment)
}

func (m *Match) addPathArgumentAtBeginning(parameter string, argument string) {
	m.addPathFragmentAtBeginning(parameter)
	m.Parameters[parameter] = argument
}

// MatchedRoute returns the predefined route this match corresponds.
func (m *Match) MatchedRoute() Route {
	if m.path == nil {
		return Route{}
	}

	builder := strings.Builder{}

	m.path.Traverse(func(s string) {
		builder.WriteRune('/')
		builder.WriteString(s)
	})

	return NewRoute(builder.String(), m.method)
}
