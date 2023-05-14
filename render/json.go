package render

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jvcoutinho/lit"
)

// JSONResult sets Content-Type header to application/json, marshals a given object to JSON
// and sets the product as the response body.
type JSONResult struct {
	// The status code of the response.
	StatusCode int
	// The body of the response (to be marshalled).
	Body any
}

// JSON sets Content-Type header to application/json, marshals obj to a JSON representation
// and sets the product as the response body.
func JSON(statusCode int, obj ...any) lit.Response {
	return func(writer http.ResponseWriter) {
		objectBytes, err := json.Marshal(obj[0])
		if err != nil {
			panic(fmt.Errorf("rendering JSON: %w", err))
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(statusCode)
		_, _ = writer.Write(objectBytes)
	}
}

// // Ok responds the request with Status Code 200 (OK) and an optional body marshalled as JSON.
// //
// // All elements of obj but the first are ignored in order to mimic an optional parameter.
// func Ok(obj ...any) *JSONResult {
// 	return NewJSONResult(http.StatusOK, obj[0])
// }
