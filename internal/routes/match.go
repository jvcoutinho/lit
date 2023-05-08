package routes

import (
	"container/list"
	"strings"
)

// Match is a correspondence between an incoming pattern and method to a pre-defined Route.
//
// As the route may have parameters, they are collected throughout the progress of a match.
type Match struct {
	method string
	path   *list.List

	// The matched arguments from the pattern.
	Parameters map[string]string
}

func newMatch() *Match {
	return &Match{
		path:       list.New(),
		Parameters: make(map[string]string),
	}
}

func (m *Match) addMethod(method string) {
	m.method = method
}

func (m *Match) addPathFragmentAtBeginning(fragment string) {
	m.path.PushFront(fragment)
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

	for iterator := m.path.Front(); iterator != nil; iterator = iterator.Next() {
		path, ok := iterator.Value.(string)
		if !ok {
			continue
		}

		builder.WriteRune('/')
		builder.WriteString(path)
	}

	return NewRoute(builder.String(), m.method)
}
