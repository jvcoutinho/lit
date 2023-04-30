package routes

import (
	"strings"

	"github.com/jvcoutinho/lit/internal/structures"
)

type Match struct {
	method string
	path   *structures.List[string]

	Parameters map[string]string
}

func NewMatch() *Match {
	return &Match{
		path:       structures.NewList[string](),
		Parameters: make(map[string]string),
	}
}

func (m *Match) AddMethod(method string) {
	m.method = method
}

func (m *Match) AddPathFragmentAtBeginning(fragment string) {
	m.path.InsertAtBeginning(fragment)
}

func (m *Match) AddPathArgumentAtBeginning(parameter string, argument string) {
	m.AddPathFragmentAtBeginning(parameter)
	m.Parameters[parameter] = argument
}

func (m *Match) MatchedRoute() Route {
	builder := strings.Builder{}

	m.path.Traverse(func(s string) {
		builder.WriteRune('/')
		builder.WriteString(s)
	})

	return NewRoute(builder.String(), m.method)
}
