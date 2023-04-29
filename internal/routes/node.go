package routes

import "strings"

type Node string

func (n Node) IsArgument() bool {
	return strings.HasPrefix(string(n), ":")
}
