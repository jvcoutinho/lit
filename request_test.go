package lit_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jvcoutinho/lit"
	"github.com/stretchr/testify/require"
)

func TestRequest(t *testing.T) {
	t.Parallel()

	type contextKey string

	tests := []struct {
		description string
		request     *http.Request
		panicValue  string
	}{
		{
			description: "WhenRequestIsNil_ShouldPanic",
			request:     nil,
			panicValue:  "request must not be nil",
		},
		{
			description: "WhenRequestIsNotNil_ShouldReturnRequest",
			request: httptest.
				NewRequest(http.MethodPost, "/users", strings.NewReader(`{"id": 1}`)).
				WithContext(context.WithValue(context.Background(), contextKey("key"), "value")),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Act
			if test.panicValue != "" {
				require.PanicsWithValue(t, test.panicValue, func() {
					_ = lit.NewRequest(test.request)
				})
				return
			}

			r := lit.NewRequest(test.request)

			// Assert
			require.Equal(t, test.request, r.Base())
			require.Equal(t, test.request.Method, r.Method())
			require.Equal(t, test.request.URL, r.URL())
			require.Equal(t, test.request.Body, r.Body())
			require.Equal(t, test.request.Header, r.Header())
			require.Equal(t, test.request.Context(), r.Context())
		})
	}
}

func TestRequest_WithContext(t *testing.T) {
	t.Parallel()

	type contextKey string

	tests := []struct {
		description string
		context     context.Context
		panicValue  string
	}{
		{
			description: "WhenContextIsNil_ShouldPanic",
			context:     nil,
			panicValue:  "nil context",
		},
		{
			description: "WhenContextIsNotNil_ShouldSetContext",
			context:     context.WithValue(context.Background(), contextKey("key"), "value"),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Arrange
			r := lit.NewRequest(
				httptest.NewRequest(http.MethodGet, "/users", nil),
			)

			// Act
			if test.panicValue != "" {
				require.PanicsWithValue(t, test.panicValue, func() {
					r.WithContext(test.context)
				})
				return
			}

			r.WithContext(test.context)

			// Assert
			require.Equal(t, test.context, r.Context())
		})
	}
}

func TestRequest_WithURIParameters(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		parameters  map[string]string
	}{
		{
			description: "WhenParametersIsNil_ShouldSet",
			parameters:  nil,
		},
		{
			description: "WhenParametersIsNotNil_ShouldSet",
			parameters:  map[string]string{"user_id": "123"},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Arrange
			r := lit.NewRequest(
				httptest.NewRequest(http.MethodGet, "/users", nil),
			)

			// Act
			r.WithURIParameters(test.parameters)

			// Assert
			require.Equal(t, test.parameters, r.URIParameters())
		})
	}
}
