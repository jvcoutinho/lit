package render

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jvcoutinho/lambda/slices"
)

// JSONResponse sets Content-Type header to application/json, marshals a given object to JSON
// and sets the product as the response body.
type JSONResponse struct {
	*httpResponse
	statusCode int
	body       any
}

func (r *JSONResponse) Write(writer http.ResponseWriter) error {
	objectBytes, err := json.Marshal(r.body)
	if err != nil {
		return fmt.Errorf("rendering JSON: %w", err)
	}

	header := writer.Header()
	for key, values := range r.Header {
		header[key] = values
	}

	writer.WriteHeader(r.statusCode)

	_, err = writer.Write(objectBytes)
	if err != nil {
		return err
	}

	return nil
}

// JSON sets Content-Type header to application/json, marshals obj to a JSON representation
// and sets the product as the response body.
func JSON(statusCode int, obj any) *JSONResponse {
	httpResponse := newHTTPResponse()
	httpResponse.Header.Set("Content-Type", "application/json")

	return &JSONResponse{
		httpResponse: httpResponse,
		statusCode:   statusCode,
		body:         obj,
	}
}

// Ok responds the request with Status Code 200 (OK) and an optional body marshalled as JSON.
func Ok(obj ...any) *JSONResponse {
	return JSON(http.StatusOK, slices.ElementAtOrDefault(obj, 0))
}

// BadRequest responds the request with Status Code 400 (Bad Request) and an optional body marshalled as JSON.
func BadRequest(obj ...any) *JSONResponse {
	return JSON(http.StatusBadRequest, slices.ElementAtOrDefault(obj, 0))
}

// Unauthorized responds the request with Status Code 401 (Unauthorized) and an optional body marshalled as JSON.
func Unauthorized(obj ...any) *JSONResponse {
	return JSON(http.StatusUnauthorized, slices.ElementAtOrDefault(obj, 0))
}

// NotFound responds the request with Status Code 404 (Not Found) and an optional body marshalled as JSON.
func NotFound(obj ...any) *JSONResponse {
	return JSON(http.StatusNotFound, slices.ElementAtOrDefault(obj, 0))
}

// Conflict responds the request with Status Code 409 (Conflict) and an optional body marshalled as JSON.
func Conflict(obj ...any) *JSONResponse {
	return JSON(http.StatusConflict, slices.ElementAtOrDefault(obj, 0))
}

// UnprocessableEntity responds the request with Status Code 422 (Unprocessable Entity) and
// an optional body marshalled as JSON.
func UnprocessableEntity(obj ...any) *JSONResponse {
	return JSON(http.StatusUnprocessableEntity, slices.ElementAtOrDefault(obj, 0))
}
