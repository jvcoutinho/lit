package lit

import (
	"log"
	"net/http"

	"github.com/jvcoutinho/lit/internal/maps"
)

// Context is an implementation of context.Context that manages the input and output of incoming requests.
type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request

	arguments map[string]string
}

// NewContext creates a new Context instance.
func NewContext(writer http.ResponseWriter, request *http.Request) *Context {
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

// SetStatusCode sets the response' status code to statusCode.
func (c *Context) SetStatusCode(statusCode int) {
	c.ResponseWriter.WriteHeader(statusCode)
}

// WriteBody writes bytes to the response.
func (c *Context) WriteBody(bytes []byte) {
	_, err := c.ResponseWriter.Write(bytes)
	if err != nil {
		log.Println(err)
	}
}

// SetHeader sets the header entries associated with key to the single element value.
// It replaces any existing values associated with key.
func (c *Context) SetHeader(key, value string) {
	header := c.ResponseWriter.Header()
	header.Set(key, value)
}
