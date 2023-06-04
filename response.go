package lit

import "net/http"

// Response is the output of a HandlerFunc.
type Response interface {
	// Write responds the request.
	Write(writer http.ResponseWriter) error
}

// ResponseFunc writes response data into http.ResponseWriter directly.
type ResponseFunc func(writer http.ResponseWriter) error

func (r ResponseFunc) Write(writer http.ResponseWriter) error {
	return r(writer)
}
