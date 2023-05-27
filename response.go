package lit

import "net/http"

// Response is the output of a Lit handler function.
type Response interface {
	// Write writes response data to writer.
	Write(writer http.ResponseWriter) error
}

// ResponseFunc is a function where one can manipulate the http.ResponseWriter directly.
type ResponseFunc func(writer http.ResponseWriter) error

func (r ResponseFunc) Write(writer http.ResponseWriter) error {
	return r(writer)
}
