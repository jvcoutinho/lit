// Package result contains functions that return responses for incoming requests.
package result

import (
	"net/http"
)

// Ok responds the request with Status Code 200 (OK).
func Ok() *TerminalResponse {
	return newTerminalResponse(http.StatusOK, nil)
}

// BadRequest responds the request with Status Code 400 (OK).
func BadRequest() *TerminalResponse {
	return newTerminalResponse(http.StatusBadRequest, nil)
}
