package lit

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/jvcoutinho/lambda/maps"
)

// Request is the input of a Lit handler function.
type Request struct {
	httpRequest  *http.Request
	uriArguments map[string]string

	// Method of this request.
	Method string
	// URI of this request.
	URI *url.URL
	// Header of this request.
	Header http.Header
	// Body of this request.
	Body io.ReadCloser
}

// NewRequest creates a new Request instance based on an incoming http.Request.
func NewRequest(httpRequest *http.Request) *Request {
	return &Request{
		httpRequest:  httpRequest,
		uriArguments: nil,
		Method:       httpRequest.Method,
		URI:          httpRequest.URL,
		Header:       httpRequest.Header,
		Body:         httpRequest.Body,
	}
}

// Context of this request. It is always non-nil.
//
// It is cancelled when the client's connection closes, the request is cancelled (with HTTP/2)
// or when the ServeHTTP method returns.
func (r *Request) Context() context.Context {
	return r.httpRequest.Context()
}

// URIArguments returns a copy of the dictionary of pattern parameters and their corresponding substitutions.
//
// For example, if the request URL path is /users/123 and the matching pattern is /users/:id,
// then URIArguments will return { ":id": "123" }.
func (r *Request) URIArguments() map[string]string {
	return maps.Copy(r.uriArguments)
}

func (r *Request) setURIArguments(uriArguments map[string]string) {
	r.uriArguments = uriArguments
}
