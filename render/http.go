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

func (r *httpResponse) Write(writer http.ResponseWriter) error {
	header := writer.Header()
	for key, values := range r.Header {
		header[key] = values
	}

	return nil
}
