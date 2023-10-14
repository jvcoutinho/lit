package lit

import (
	"net/http"
)

// Request is the input of a HandlerFunc.
type Request struct {
	httpRequest *http.Request
	arguments   map[string]string
}

func newRequest(httpRequest *http.Request, arguments map[string]string) *Request {
	return &Request{
		httpRequest,
		arguments,
	}
}

// Arguments returns this request's URL path arguments matched against the pattern parameters.
//
// For example, if the pattern is "/users/:id" and the URI is "/users/123",
// Arguments' result will contain the { "id": "123" } key-value pair.
func (r *Request) Arguments() map[string]string {
	return r.arguments
}
