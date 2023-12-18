package render

import "net/http"

// NoContentResponse is a lit.Response without a body and status code 204 No Content.
type NoContentResponse struct {
	Header http.Header
}

func (r NoContentResponse) Write(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)

	for key := range r.Header {
		w.Header().Set(key, r.Header.Get(key))
	}
}

// WithHeader sets the response header entries associated with key to value.
func (r NoContentResponse) WithHeader(key, value string) NoContentResponse {
	r.Header.Set(key, value)
	return r
}

// NoContent responds the request with [204 No Content].
//
// [204 No Content]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/204
func NoContent() NoContentResponse {
	return NoContentResponse{make(http.Header)}
}
