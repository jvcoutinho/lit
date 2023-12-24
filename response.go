package lit

import "net/http"

// Response is the output of a [Handler].
type Response interface {
	// Write responds the request.
	Write(w http.ResponseWriter)
}

// ResponseFunc writes response data into [http.ResponseWriter].
type ResponseFunc func(w http.ResponseWriter)

func (r ResponseFunc) Write(w http.ResponseWriter) {
	r(w)
}
