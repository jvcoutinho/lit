package routes

import "strings"

type Match struct {
	method string
	path   *strings.Builder

	Parameters map[string]string
}

func NewMatch() *Match {
	return &Match{
		path:       &strings.Builder{},
		Parameters: make(map[string]string),
	}
}

func (m *Match) AddMethod(method string) {
	m.method = method
}

func (m *Match) AddPathFragment(fragment string) {
	m.path.WriteRune('/')
	m.path.WriteString(fragment)
}

func (m *Match) Len() int {
	return m.path.Len()
}

func (m *Match) AddPathArgument(parameter string, argument string) {
	m.AddPathFragment(parameter)
	m.Parameters[parameter] = argument
}

func (m *Match) MatchedRoute() Route {
	return NewRoute(m.path.String(), m.method)
}
