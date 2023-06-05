package lit

import (
	"net/http"

	"github.com/jvcoutinho/lambda/maps"
)

// Request is the input of a HandlerFunc.
type Request struct {
	httpRequest *http.Request

	uriArguments map[string]string
}

func newRequest(httpRequest *http.Request) *Request {
	return &Request{
		httpRequest,
		make(map[string]string),
	}
}

// URIArguments returns a copy of the arguments matched with the pattern parameters.
func (r *Request) URIArguments() map[string]string {
	return maps.Copy(r.uriArguments)
}

func (r *Request) setURIArguments(uriArguments map[string]string) {
	r.uriArguments = uriArguments
}
