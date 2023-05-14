package lit

import "net/http"

// Response is the output of a Lit handler function.
type Response interface {
	Write(writer http.ResponseWriter) error
}

// CustomResponse defines a fully customizable response.
type CustomResponse func(writer http.ResponseWriter) error

func (r CustomResponse) Write(writer http.ResponseWriter) error {
	return r(writer)
}
