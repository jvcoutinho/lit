package routes

import (
	"fmt"
	"strings"
)

// Route is a representation of an HTTP endpoint.
type Route struct {
	Pattern string
	Method  string

	path []string
}

// NewRoute creates a new route instance.
func NewRoute(pattern, method string) Route {
	pattern = strings.Trim(pattern, "/")
	path := strings.Split(pattern, "/")
	method = strings.ToUpper(method)

	return Route{pattern, method, path}
}

// Path returns each part of route's pattern.
func (r Route) Path() []string {
	return r.path
}

// String returns "r.Method /r.Pattern".
func (r Route) String() string {
	return fmt.Sprintf("%s /%s", r.Method, r.Pattern)
}

func isArgument(path string) bool {
	return strings.HasPrefix(path, ":")
}
