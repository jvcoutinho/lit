package render

import (
	"fmt"
	"net/http"

	"gopkg.in/yaml.v3"
)

// YAMLResponse sets Content-Type header to text/yaml, marshals a given object to YAML
// and sets the product as the response body.
type YAMLResponse struct {
	*HTTPResponse

	body any
}

func (r *YAMLResponse) Write(writer http.ResponseWriter) error {
	r.header.Set("Content-Type", "text/yaml")

	objectBytes, err := yaml.Marshal(r.body)
	if err != nil {
		return fmt.Errorf("rendering YAML: %w", err)
	}

	r.SetBody(objectBytes)

	return r.HTTPResponse.Write(writer)
}

// YAML sets Content-Type header to text/yaml, marshals obj to a YAML representation
// and sets the product as the response body.
func YAML(statusCode int, obj any) *YAMLResponse {
	return &YAMLResponse{
		NewHTTPResponse(statusCode, nil),
		obj,
	}
}

// OkYAML responds the request with Status Code 200 (OK) and a body marshalled as YAML.
func OkYAML(obj any) *YAMLResponse {
	return YAML(http.StatusOK, obj)
}

// BadRequestYAML responds the request with Status Code 400 (Bad Request) and a body marshalled as YAML.
func BadRequestYAML(obj any) *YAMLResponse {
	return YAML(http.StatusBadRequest, obj)
}

// UnauthorizedYAML responds the request with Status Code 401 (Unauthorized) and a body marshalled as YAML.
func UnauthorizedYAML(obj any) *YAMLResponse {
	return YAML(http.StatusUnauthorized, obj)
}

// NotFoundYAML responds the request with Status Code 404 (Not Found) and a body marshalled as YAML.
func NotFoundYAML(obj any) *YAMLResponse {
	return YAML(http.StatusNotFound, obj)
}

// ConflictYAML responds the request with Status Code 409 (Conflict) and a body marshalled as YAML.
func ConflictYAML(obj any) *YAMLResponse {
	return YAML(http.StatusConflict, obj)
}

// UnprocessableEntityYAML responds the request with Status Code 422 (Unprocessable Entity) and
// an optional body marshalled as YAML.
func UnprocessableEntityYAML(obj any) *YAMLResponse {
	return YAML(http.StatusUnprocessableEntity, obj)
}

// InternalServerErrorYAML responds the request with Status Code 500 (Internal Server Error) and
// an optional body marshalled as YAML.
func InternalServerErrorYAML(obj any) *YAMLResponse {
	return YAML(http.StatusInternalServerError, obj)
}
