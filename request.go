package lit

import "net/http"

// Request is the input of a HandlerFunc.
type Request struct {
	httpRequest *http.Request
}

func newRequest(httpRequest *http.Request) *Request {
	return &Request{httpRequest}
}
