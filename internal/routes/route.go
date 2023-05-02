package routes

import (
	"fmt"
	"strings"

	"github.com/jvcoutinho/lit/internal/sets"
	"github.com/jvcoutinho/lit/internal/slices"
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

// Path returns each part of route's pattern.
func (r Route) Path() []string {
	return strings.Split(r.Pattern, "/")
}

// String returns "r.Method /r.Pattern".
func (r Route) String() string {
	return fmt.Sprintf("%s /%s", r.Method, r.Pattern)
}

func hasDuplicateArguments(path []string) (string, bool) {
	arguments := slices.Filter(path, isArgument)
	set := sets.NewHashSet[string]()

	for _, argument := range arguments {
		if set.Contains(argument) {
			return argument, true
		}

		set.Add(argument)
	}

	return "", false
}

func isArgument(path string) bool {
	return strings.HasPrefix(path, ":")
}
