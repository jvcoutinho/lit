package render_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jvcoutinho/lit/render"
	"github.com/stretchr/testify/require"
)

func TestJSONResponse_Write(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		description   string
		statusCode    int
		object        any
		headerToAdd   map[string]string
		expectedBody  string
		expectedError error
	}

	tests := []TestCase{
		{
			description: "GivenBodyCanBeMarshalled_ShouldWriteBody_AndShouldWriteStatusCode_AndShouldWriteHeader",
			statusCode:  http.StatusOK,
			object: map[string]any{
				"Key1": 123,
				"Key2": "value2",
				"Key3": []int{2, 3},
			},
			headerToAdd:   map[string]string{"key": "value"},
			expectedBody:  `{"Key1":123,"Key2":"value2","Key3":[2,3]}`,
			expectedError: nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Arrange
			response := render.JSON(test.statusCode, test.object)

			for key, value := range test.headerToAdd {
				response.Header().Add(key, value)
			}

			recorder := httptest.NewRecorder()

			// Act
			actualError := response.Write(recorder)

			// Assert
			require.Equal(t, test.expectedError, actualError)
			require.Equal(t, recorder.Code, test.statusCode)
			require.Equal(t, test.expectedBody, recorder.Body.String())
			require.Equal(t, recorder.Header(), response.Header())
			require.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
		})
	}
}

func TestJSONResponseConstructors(t *testing.T) {
	t.Parallel()

}
