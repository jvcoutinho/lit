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
}

// NewRequest creates a new Request instance based on an incoming http.Request.
func NewRequest(httpRequest *http.Request) *Request {
	return &Request{
		httpRequest:  httpRequest,
		uriArguments: nil,
	}
}

// Context of this request. It is always non-nil.
//
// It is cancelled when the client's connection closes, the request is cancelled (with HTTP/2)
// or when the ServeHTTP method returns.
func (r *Request) Context() context.Context {
	return r.httpRequest.Context()
}

// Method specifies the HTTP method of this request.
func (r *Request) Method() string {
	return r.httpRequest.Method
}

// URI being requested by this request.
func (r *Request) URI() *url.URL {
	return r.httpRequest.URL
}

// Header of this request.
func (r *Request) Header() http.Header {
	return r.httpRequest.Header
}

// Body of this request.
func (r *Request) Body() io.ReadCloser {
	return r.httpRequest.Body
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
