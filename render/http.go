package render

import "net/http"

type httpResponse struct {
	// Header of this response.
	Header http.Header
}

func newHTTPResponse() *httpResponse {
	return &httpResponse{
		make(http.Header),
	}
}
