package routes

import "strings"

type Match struct {
	method string
	path   []string

	Parameters map[string]string
}

func NewMatch() *Match {
	return &Match{
		path:       make([]string, 0),
		Parameters: make(map[string]string),
	}
}

func (m *Match) AddMethod(method string) {
	m.method = method
}

func (m *Match) AddPathFragment(fragment string) {
	m.path = append(m.path, "/"+fragment)
}

func (m *Match) AddPathArgument(parameter string, argument string) {
	m.AddPathFragment(parameter)
	m.Parameters[parameter] = argument
}

func (m *Match) MatchedRoute() Route {
	builder := strings.Builder{}
	for i := len(m.path) - 1; i >= 0; i-- {
		builder.WriteString(m.path[i])
	}

	return NewRoute(builder.String(), m.method)
}
