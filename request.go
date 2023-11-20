package lit

import (
	"io"
	"net/http"
	"net/url"
)

// Request is the input of a [Handler].
type Request struct {
	Request    *http.Request
	parameters map[string]string
}

// NewRequest creates a new [Request] instance with an underlying [*http.Request].
//
// If request is nil, NewRequest panics.
func NewRequest(request *http.Request, parameters map[string]string) *Request {
	if request == nil {
		panic("request must not be nil")
	}

	return &Request{
		request,
		parameters,
	}
}

// URIParameters returns this request's URL path parameters and their values.
//
// Use [bind.URIParameters] for standard model binding and validation features.
//
// The keys from this map don't start with ":" prefix.
func (r *Request) URIParameters() map[string]string {
	return r.parameters
}

// URL of this request.
func (r *Request) URL() *url.URL {
	return r.Request.URL
}

// Method of this request.
func (r *Request) Method() string {
	return r.Request.Method
}

// Body of this request.
func (r *Request) Body() io.ReadCloser {
	return r.Request.Body
}

// Header fields of this request.
func (r *Request) Header() http.Header {
	return r.Request.Header
}

// Base returns the underlying http.Request.
func (r *Request) Base() *http.Request {
	return r.Request
}
