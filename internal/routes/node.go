package routes

import "strings"

// Node is a representation of either a route path or an HTTP method.
type Node string

// IsArgument returns true if the current node is an argument (i.e. has the format ":name").
func (n Node) IsArgument() bool {
	return strings.HasPrefix(string(n), ":")
}
