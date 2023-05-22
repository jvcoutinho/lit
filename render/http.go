package render

import "net/http"

type httpResponse struct {
	header http.Header
}

func newHTTPResponse() *httpResponse {
	return &httpResponse{
		make(http.Header),
	}
}

func (r *httpResponse) Write(writer http.ResponseWriter) error {
	header := writer.Header()
	for key, values := range r.header {
		header[key] = values
	}

	return nil
}

// Header of this response.
func (r *httpResponse) Header() http.Header {
	return r.header
}
