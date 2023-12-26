package lit_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"

	"github.com/jvcoutinho/lit"
	"github.com/stretchr/testify/require"
)

func TestRecover(t *testing.T) {
	t.Parallel()

	// Arrange
	r := lit.NewRequest(
		httptest.NewRequest(http.MethodGet, "/users", nil),
	)

	panicHandler := lit.Handler(func(r *lit.Request) lit.Response {
		panic("scary!")
	})

	recorder := httptest.NewRecorder()

	// Act
	response := lit.Recover(panicHandler)(r)
	response.Write(recorder)

	// Assert
	require.Equal(t, http.StatusInternalServerError, recorder.Code)
	require.Equal(t, "scary!\n", recorder.Body.String())
	require.Equal(t, http.Header{
		"Content-Type":           {"text/plain; charset=utf-8"},
		"X-Content-Type-Options": {"nosniff"},
	}, recorder.Header())
}

func TestLog(t *testing.T) {
	tests := []struct {
		description     string
		writer          http.ResponseWriter
		statusCode      int
		expectedContent *regexp.Regexp
	}{
		{
			description:     "WhenWriterIsNotRecorder_ShouldNotLog",
			writer:          httptest.NewRecorder(),
			statusCode:      http.StatusNotFound,
			expectedContent: regexp.MustCompile(""),
		},
		{
			description: "WhenWriterIsRecorder_AndStatusCodeIs2xx_ShouldLogGreen",
			writer:      lit.NewRecorder(httptest.NewRecorder()),
			statusCode:  http.StatusOK,
			expectedContent: regexp.MustCompile(
				"^\n\u001B\\[97;1;42m>> GET /users\u001B\\[0m\n> 200 OK\n> Start Time: .+\n> Remote Address: .+\n> Duration: .+\n> Content-Length: 3\n$",
			),
		},
		{
			description: "WhenWriterIsRecorder_AndStatusCodeIs3xx_ShouldLogBlue",
			writer:      lit.NewRecorder(httptest.NewRecorder()),
			statusCode:  http.StatusPermanentRedirect,
			expectedContent: regexp.MustCompile(
				"^\n\u001B\\[97;1;104m>> GET /users\u001B\\[0m\n> 308 Permanent Redirect\n> Start Time: .+\n> Remote Address: .+\n> Duration: .+\n> Content-Length: 3\n$",
			),
		},
		{
			description: "WhenWriterIsRecorder_AndStatusCodeIs4xx_ShouldLogYellow",
			writer:      lit.NewRecorder(httptest.NewRecorder()),
			statusCode:  http.StatusNotFound,
			expectedContent: regexp.MustCompile(
				"^\n\u001B\\[97;1;43m>> GET /users\u001B\\[0m\n> 404 Not Found\n> Start Time: .+\n> Remote Address: .+\n> Duration: .+\n> Content-Length: 3\n$",
			),
		},
		{
			description: "WhenWriterIsRecorder_AndStatusCodeIs5xx_ShouldLogRed",
			writer:      lit.NewRecorder(httptest.NewRecorder()),
			statusCode:  http.StatusInternalServerError,
			expectedContent: regexp.MustCompile(
				"^\n\u001B\\[97;1;41m>> GET /users\u001B\\[0m\n> 500 Internal Server Error\n> Start Time: .+\n> Remote Address: .+\n> Duration: .+\n> Content-Length: 3\n$",
			),
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			// Arrange
			r := lit.NewRequest(
				httptest.NewRequest(http.MethodGet, "/users", nil),
			)

			handler := lit.Handler(func(r *lit.Request) lit.Response {
				return lit.ResponseFunc(func(w http.ResponseWriter) {
					w.WriteHeader(test.statusCode)
					w.Write([]byte("log"))
				})
			})

			var output bytes.Buffer
			log.SetOutput(&output)
			t.Cleanup(func() {
				log.SetOutput(os.Stderr)
			})

			// Act
			response := lit.Log(handler)(r)
			response.Write(test.writer)

			// Assert
			require.Regexp(t, test.expectedContent, output.String())
		})
	}

}
