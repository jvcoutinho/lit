package lit_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jvcoutinho/lit"
	"github.com/stretchr/testify/require"
)

func TestRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description   string
		request       *http.Request
		expectedPanic string
	}{
		{
			description:   "WhenRequestIsNil_ShouldPanic",
			request:       nil,
			expectedPanic: "request must not be nil",
		},
		{
			description: "WhenRequestIsNotNil_ShouldCreateRequest",
			request:     httptest.NewRequest(http.MethodGet, "/path/sub-path", bytes.NewBufferString("body")),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Arrange
			parameters := map[string]string{"param_1": "value_1"}

			// Act
			if test.expectedPanic != "" {
				require.PanicsWithValue(t, test.expectedPanic, func() {
					_ = lit.NewRequest(test.request, parameters)
				})

				return
			}

			r := lit.NewRequest(test.request, parameters)

			// Assert
			require.Equal(t, test.request.Header, r.Header())
			require.Equal(t, test.request.URL, r.URL())
			require.Equal(t, test.request.Body, r.Body())
			require.Equal(t, test.request.Method, r.Method())
			require.Equal(t, parameters, r.URIParameters())
		})
	}
}
