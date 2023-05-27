package render

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// JSONResponse sets Content-Type header to application/json, marshals a given object to JSON
// and sets the product as the response body.
type JSONResponse struct {
	*HTTPResponse

	body any
}

func (r *JSONResponse) Write(writer http.ResponseWriter) error {
	r.header.Set("Content-Type", "application/json")

	objectBytes, err := json.Marshal(r.body)
	if err != nil {
		return fmt.Errorf("rendering JSON: %w", err)
	}

	r.setBody(objectBytes)

	return r.HTTPResponse.Write(writer)
}

// JSON sets Content-Type header to application/json, marshals obj to a JSON representation
// and sets the product as the response body.
func JSON(statusCode int, obj any) *JSONResponse {
	return &JSONResponse{
		NewHTTPResponse(statusCode, nil),
		obj,
	}
}

// OkJSON responds the request with Status Code 200 (OK) and a body marshalled as JSON.
func OkJSON(obj any) *JSONResponse {
	return JSON(http.StatusOK, obj)
}

// BadRequestJSON responds the request with Status Code 400 (Bad Request) and a body marshalled as JSON.
func BadRequestJSON(obj any) *JSONResponse {
	return JSON(http.StatusBadRequest, obj)
}

// UnauthorizedJSON responds the request with Status Code 401 (Unauthorized) and a body marshalled as JSON.
func UnauthorizedJSON(obj any) *JSONResponse {
	return JSON(http.StatusUnauthorized, obj)
}

// NotFoundJSON responds the request with Status Code 404 (Not Found) and a body marshalled as JSON.
func NotFoundJSON(obj any) *JSONResponse {
	return JSON(http.StatusNotFound, obj)
}

// ConflictJSON responds the request with Status Code 409 (Conflict) and a body marshalled as JSON.
func ConflictJSON(obj any) *JSONResponse {
	return JSON(http.StatusConflict, obj)
}

// UnprocessableEntityJSON responds the request with Status Code 422 (Unprocessable Entity) and
// an optional body marshalled as JSON.
func UnprocessableEntityJSON(obj any) *JSONResponse {
	return JSON(http.StatusUnprocessableEntity, obj)
}

// InternalServerErrorJSON responds the request with Status Code 500 (Internal Server Error) and
// an optional body marshalled as JSON.
func InternalServerErrorJSON(obj any) *JSONResponse {
	return JSON(http.StatusInternalServerError, obj)
}
