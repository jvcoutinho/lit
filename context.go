package lit

import (
	"net/http"
)

// Context manages the input and output of incoming requests.
type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

func newContext(writer http.ResponseWriter, request *http.Request) *Context {
	return &Context{writer, request}
}
