package routes

import (
	"fmt"
	"strings"
)

// Route is a representation of an HTTP endpoint.
type Route struct {
	Pattern string
	Method  string
}

// NewRoute creates a new route instance.
func NewRoute(pattern, method string) Route {
	pattern = strings.Trim(pattern, "/")
	method = strings.ToUpper(method)

	return Route{pattern, method}
}

// String returns "r.Method /r.Pattern".
func (r Route) String() string {
	return fmt.Sprintf("%s /%s", r.Method, r.Pattern)
}
