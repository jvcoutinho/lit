package lit

import (
	"io"
	"net/http"
	"net/url"
)

// Request is the input of a [Handler].
type Request struct {
	request   *http.Request
	arguments map[string]string
}

// NewRequest creates a new [Request] instance with an underlying [*http.Request].
func NewRequest(request *http.Request, arguments map[string]string) *Request {
	return &Request{
		request,
		arguments,
	}
}

// Arguments returns this request's URL path arguments.
//
// This is intended for advanced usage. For regular usage,
// [bind] package is preferred due to model binding and validation features.
//
// The keys from this map don't start with ":" prefix.
func (r *Request) Arguments() map[string]string {
	return r.arguments
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
