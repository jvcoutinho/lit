package render

import (
	"encoding/json"
	"net/http"
)

type messageResponse struct {
	Message string `json:"message"`
}

// JSONResponse is a lit.Response that prints a JSON formatted-body as response. It sets
// the Content-Type header to "application/json".
//
// If the body is set but its marshalling fails, JSON responds an Internal Server Error
// with the error message as plain text.
type JSONResponse struct {
	StatusCode int
	Header     http.Header
	Body       any
}

// JSON responds the request with statusCode and an optional body (it can be nil) marshalled as JSON.
//
// If body is a string or an error, it's marshalled as the value of a "message" key.
// Otherwise, it's marshalled as is.
func JSON(statusCode int, body any) JSONResponse {
	switch cast := body.(type) {
	case string:
		return JSONResponse{statusCode, make(http.Header), messageResponse{cast}}
	case error:
		return JSONResponse{statusCode, make(http.Header), messageResponse{cast.Error()}}
	default:
		return JSONResponse{statusCode, make(http.Header), cast}
	}
}

// WithHeader sets the response header entries associated with key to value.
func (r JSONResponse) WithHeader(key, value string) JSONResponse {
	r.Header.Set(key, value)
	return r
}

func (r JSONResponse) Write(w http.ResponseWriter) {
	for key := range r.Header {
		w.Header().Set(key, r.Header.Get(key))
	}

	if r.Body == nil {
		w.WriteHeader(r.StatusCode)
		return
	}

	if err := r.writeBody(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}

func (r JSONResponse) writeBody(w http.ResponseWriter) error {
	bodyBytes, err := json.Marshal(r.Body)
	if err != nil {
		return err
	}

	w.WriteHeader(r.StatusCode)

	_, err = w.Write(bodyBytes)

	return err
}

// OK responds the request with [200 OK].
//
// [200 OK]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/200
func OK(body any) JSONResponse {
	return JSON(http.StatusOK, body)
}

// Created responds the request with [201 Created] and the URL of the created resource in the Location header.
//
// [201 Created]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/201
func Created(body any, locationURL string) JSONResponse {
	return JSON(http.StatusCreated, body).WithHeader("Location", locationURL)
}

// Accepted responds the request with [202 Accepted].
//
// [202 Accepted]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/202
func Accepted(body any) JSONResponse {
	return JSON(http.StatusAccepted, body)
}

// BadRequest responds the request with [400 Bad Request].
//
// [400 Bad Request]: https://developer.mozilla.org/en-US/docs/web/http/status/400
func BadRequest(body any) JSONResponse {
	return JSON(http.StatusBadRequest, body)
}

// Unauthorized responds the request with [401 Unauthorized].
//
// [401 Unauthorized]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/401
func Unauthorized(body any) JSONResponse {
	return JSON(http.StatusUnauthorized, body)
}

// Forbidden responds the request with [403 Forbidden].
//
// [403 Forbidden]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/403
func Forbidden(body any) JSONResponse {
	return JSON(http.StatusForbidden, body)
}

// NotFound responds the request with [404 Not Found].
//
// [404 Not Found]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/404
func NotFound(body any) JSONResponse {
	return JSON(http.StatusNotFound, body)
}

// Conflict responds the request with [409 Conflict].
//
// [409 Conflict]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/409
func Conflict(body any) JSONResponse {
	return JSON(http.StatusConflict, body)
}

// UnprocessableContent responds the request with [422 Unprocessable Content].
//
// [422 Unprocessable Content]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/422
func UnprocessableContent(body any) JSONResponse {
	return JSON(http.StatusUnprocessableEntity, body)
}

// InternalServerError responds the request with [500 Internal Server Error].
//
// [500 Internal Server Error]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/500
func InternalServerError(body any) JSONResponse {
	return JSON(http.StatusInternalServerError, body)
}
