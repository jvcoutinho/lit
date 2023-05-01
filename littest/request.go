// Package littest contains helper functions for testing HTTP routes defined in lit.Router instances.
package littest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jvcoutinho/lit"
)

// Request requests against router and returns its response.
func Request(t *testing.T, router *lit.Router, request *http.Request) *httptest.ResponseRecorder {
	t.Helper()

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	return recorder
}
