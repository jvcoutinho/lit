package lit

import (
	"net/http"
)

// Request is the input of a HandlerFunc.
type Request struct {
	httpRequest *http.Request

	urlArguments map[string]string
}

func newRequest(httpRequest *http.Request) *Request {
	return &Request{
		httpRequest,
		make(map[string]string),
	}
}

// URLArguments returns this request's URL path arguments matched against the pattern parameters.
//
// For example, if the pattern is "/users/:id" and the URI is "/users/123",
// URLArguments' result will contain the { "id": "123" } key-value pair.
func (r *Request) URLArguments() map[string]string {
	return r.urlArguments
}

func (r *Request) setURLArguments(arguments map[string]string) {
	r.urlArguments = arguments
}
