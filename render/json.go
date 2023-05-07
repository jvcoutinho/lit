package render

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jvcoutinho/lit"
	"github.com/jvcoutinho/lit/internal/slices"
)

// JSONResult sets Content-Type header to application/json, marshals a given object to JSON
// and sets the product as the response body.
type JSONResult struct {
	// The status code of the response.
	StatusCode int
	// The body of the response (to be marshalled).
	Body any
}

// NewJSONResult creates a new JSONResult instance.
func NewJSONResult(statusCode int, obj any) *JSONResult {
	return &JSONResult{statusCode, obj}
}

func (r *JSONResult) Render(ctx *lit.Context) error {
	objectBytes, err := json.Marshal(r.Body)
	if err != nil {
		return fmt.Errorf("JSONResult.Render: %w", err)
	}

	ctx.SetStatusCode(r.StatusCode)
	ctx.SetHeader("Content-Type", "application/json")
	ctx.WriteBody(objectBytes)

	return nil
}

// Ok responds the request with Status Code 200 (OK) and an optional body marshalled as JSON.
//
// All elements of obj but the first are ignored in order to mimic an optional parameter.
func Ok(obj ...any) *JSONResult {
	return NewJSONResult(http.StatusOK, slices.ElementAtOrDefault(obj, 0))
}
