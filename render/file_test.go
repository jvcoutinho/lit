package render_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jvcoutinho/lit"
	"github.com/jvcoutinho/lit/render"
	"github.com/stretchr/testify/require"
)

func TestFileResponse(t *testing.T) {
	t.Parallel()

	readFileContent := func(t *testing.T, fileName string) string {
		t.Helper()

		content, err := os.ReadFile(fileName)
		require.NoError(t, err)

		return string(content)
	}

	tests := []struct {
		description           string
		filePath              string
		expectedStatusCode    int
		expectedBody          string
		expectedContentLength int
		expectedHeader        http.Header
	}{
		{
			description:        "WhenFileDoesNotExist_ShouldRespondNotFound",
			filePath:           "not-a-valid-file",
			expectedStatusCode: http.StatusNotFound,
			expectedBody:       "404 page not found\n",
			expectedHeader: http.Header{
				"Content-Type":           {"text/plain; charset=utf-8"},
				"X-Content-Type-Options": {"nosniff"},
			},
		},
		{
			description:        "WhenFileExists_ShouldRespondWithItsBytes",
			filePath:           "./file_test.go",
			expectedStatusCode: http.StatusOK,
			expectedBody:       readFileContent(t, "./file_test.go"),
			expectedHeader: http.Header{
				"Accept-Ranges": {"bytes"},
			},
		},
		{
			description:        "WhenFileIsDirectory_ShouldRedirect",
			filePath:           "./",
			expectedStatusCode: http.StatusMovedPermanently,
			expectedBody:       "",
			expectedHeader: http.Header{
				"Location": {"stream/"},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Arrange
			var (
				writer   = httptest.NewRecorder()
				request  = lit.NewRequest(httptest.NewRequest(http.MethodGet, "/stream", nil))
				response = render.File(request, test.filePath)
			)

			// Act
			response.Write(writer)

			// Assert
			require.Equal(t, test.expectedStatusCode, writer.Code)
			require.Equal(t, test.expectedBody, writer.Body.String())
			require.Subset(t, writer.Header(), test.expectedHeader)
		})
	}
}
