package lit

import (
	"io"
	"net/http"
	"net/url"
)

// Request is the input of a [Handler].
type Request struct {
	request    *http.Request
	parameters map[string]string
}

// NewRequest creates a new [Request] instance with an underlying [*http.Request].
func NewRequest(request *http.Request, parameters map[string]string) *Request {
	return &Request{
		request,
		parameters,
	}
}

// URLParameters returns this request's URL path parameters and their values.
//
// Use [bind.URLParameters] for standard model binding and validation features.
//
// The keys from this map don't start with ":" prefix.
func (r *Request) URLParameters() map[string]string {
	return r.parameters
}

// URL of this request.
func (r *Request) URL() *url.URL {
	return r.request.URL
}

// Method of this request.
func (r *Request) Method() string {
	return r.request.Method
}

// Body of this request.
func (r *Request) Body() io.ReadCloser {
	return r.request.Body
}

// Header of this request.
func (r *Request) Header() http.Header {
	return r.request.Header
}

// Base returns the underlying http.Request.
func (r *Request) Base() *http.Request {
	return r.request
}
