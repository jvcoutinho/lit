package render

import "net/http"

// HTTPResponse is a stardard HTTP response consisting of a status code, a body and a header.
type HTTPResponse struct {
	statusCode int
	body       []byte
	header     http.Header
}

// NewHTTPResponse creates a new HTTPResponse instance.
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

// setBody sets this response's body.
func (r *HTTPResponse) setBody(body []byte) {
	r.body = body
}

// Ok responds the request with Status Code 200 (OK).
func Ok() *HTTPResponse {
	return NewHTTPResponse(http.StatusOK, nil)
}

// NoContent responds the request with Status Code 204 (No Content).
func NoContent() *HTTPResponse {
	return NewHTTPResponse(http.StatusNoContent, nil)
}

// BadRequest responds the request with Status Code 400 (Bad Request).
func BadRequest() *HTTPResponse {
	return NewHTTPResponse(http.StatusBadRequest, nil)
}

// Unauthorized responds the request with Status Code 401 (Unauthorized).
func Unauthorized() *HTTPResponse {
	return NewHTTPResponse(http.StatusUnauthorized, nil)
}

// NotFound responds the request with Status Code 404 (Not Found).
func NotFound() *HTTPResponse {
	return NewHTTPResponse(http.StatusNotFound, nil)
}

// Conflict responds the request with Status Code 409 (Conflict).
func Conflict() *HTTPResponse {
	return NewHTTPResponse(http.StatusConflict, nil)
}

// UnprocessableEntity responds the request with Status Code 422 (Unprocessable Entity).
func UnprocessableEntity() *HTTPResponse {
	return NewHTTPResponse(http.StatusUnprocessableEntity, nil)
}

// InternalServerError responds the request with Status Code 500 (Internal Server Error) and
// an optional body marshalled as .
func InternalServerError() *HTTPResponse {
	return NewHTTPResponse(http.StatusInternalServerError, nil)
}
