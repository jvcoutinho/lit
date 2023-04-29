package lit

import "net/http"

type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

func NewContext(writer http.ResponseWriter, request *http.Request) *Context {
	return &Context{writer, request}
}
