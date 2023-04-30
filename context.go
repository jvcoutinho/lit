package lit

import (
	"net/http"

	"github.com/jvcoutinho/lit/internal/maps"
)

// Context is an implementation of context.Context that manages the input and output of incoming requests.
type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request

	arguments map[string]string
}

func newContext(writer http.ResponseWriter, request *http.Request) *Context {
	return &Context{
		ResponseWriter: writer,
		Request:        request,
		arguments:      nil,
	}
}

func (c *Context) setArguments(arguments map[string]string) {
	c.arguments = arguments
}

// URIArguments returns a copy of the dictionary of pattern parameters and their corresponding substitutions.
//
// For example, if the request URL path is /users/123 and the matching pattern is /users/:id,
// then URIArguments will return { ":id": "123" }.
//
// This is best suited for custom manipulation of the arguments, but one should use bind.URI for
// regular usage.
func (c *Context) URIArguments() map[string]string {
	return maps.Copy(c.arguments)
}
