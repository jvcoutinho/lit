package render

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// JSONResponse sets Content-Type header to application/json, marshals a given object to JSON
// and sets the product as the response body.
type JSONResponse struct {
	statusCode int
	body       any
	header     http.Header
}

func (r *JSONResponse) Write(writer http.ResponseWriter) error {
	objectBytes, err := json.Marshal(r.body)
	if err != nil {
		return fmt.Errorf("rendering JSON: %w", err)
	}

	header := writer.Header()
	for key, values := range r.header {
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
	header := make(http.Header)
	header.Set("Content-Type", "application/json")

	return &JSONResponse{
		statusCode: statusCode,
		body:       obj,
		header:     header,
	}
}

// Ok responds the request with Status Code 200 (OK) and an optional body marshalled as JSON.
//
// All elements of obj but the first are ignored in order to mimic an optional parameter.
func Ok(obj ...any) *JSONResponse {
	return JSON(http.StatusOK, obj[0])
}
