package lit

import "net/http"

// Response is the output of a HandlerFunc.
type Response interface {
	// Write responds the request.
	Write(writer http.ResponseWriter)
}

// ResponseFunc writes response data into http.ResponseWriter directly.
type ResponseFunc func(writer http.ResponseWriter)

func (r ResponseFunc) Write(writer http.ResponseWriter) {
	r(writer)
}
