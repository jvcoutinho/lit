package lit

import (
	"io"
	"net/http"
	"net/url"
)

// Request is the input of a [Handler].
type Request struct {
	base       *http.Request
	parameters map[string]string
}

// NewRequest creates a new [Request] instance from a [*http.Request].
//
// If request is nil, NewRequest panics.
func NewRequest(request *http.Request) *Request {
	if request == nil {
		panic("request must not be nil")
	}

	return &Request{
		request,
		nil,
	}
}

// WithURIParameters sets the URI parameters (associated with their values) of this request.
func (r *Request) WithURIParameters(parameters map[string]string) *Request {
	r.parameters = parameters
	return r
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
	return r.base.URL
}

// Method of this request.
func (r *Request) Method() string {
	return r.base.Method
}

// Body of this request.
func (r *Request) Body() io.ReadCloser {
	return r.base.Body
}

// Header fields of this request.
func (r *Request) Header() http.Header {
	return r.base.Header
}

// Base returns the equivalent [*http.Request] of this request.
func (r *Request) Base() *http.Request {
	return r.base
}
