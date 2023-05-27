package render_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jvcoutinho/lit/render"
	"github.com/stretchr/testify/require"
)

func TestHTTPResponse_Write(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		description string
		statusCode  int
		body        []byte
		headerToAdd map[string]string
	}

	tests := []TestCase{
		{
			description: "ShouldWriteBody_AndShouldWriteStatusCode_AndShouldWriteHeader",
			statusCode:  http.StatusOK,
			body:        []byte("body"),
			headerToAdd: map[string]string{"key": "value"},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Arrange
			response := render.NewHTTPResponse(test.statusCode, test.body)

			for key, value := range test.headerToAdd {
				response.Header().Add(key, value)
			}

			recorder := httptest.NewRecorder()

			// Act
			actualError := response.Write(recorder)

			// Assert
			require.Nil(t, actualError)
			require.Equal(t, recorder.Code, test.statusCode)
			require.Equal(t, recorder.Body.Bytes(), test.body)
			require.Equal(t, recorder.Header(), response.Header())
		})
	}
}
