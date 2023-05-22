package render

import "net/http"

// HTTPResponse is a stardard HTTP response consisting of a status code, a body and a header.
type HTTPResponse struct {
	statusCode int
	body       []byte
	header     http.Header
}

func NewHTTPResponse(statusCode int, body []byte) *HTTPResponse {
	return &HTTPResponse{
		statusCode,
		body,
		make(http.Header),
	}
}

func (r *HTTPResponse) Write(writer http.ResponseWriter) error {
	header := writer.Header()
	for key, values := range r.header {
		header[key] = values
	}

	writer.WriteHeader(r.statusCode)

	_, err := writer.Write(r.body)
	if err != nil {
		return err
	}

	return nil
}

// Header of this response.
func (r *HTTPResponse) Header() http.Header {
	return r.header
}

// SetBody sets this response's body.
func (r *HTTPResponse) SetBody(body []byte) {
	r.body = body
}
