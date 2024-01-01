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
		response        lit.Response
		expectedContent *regexp.Regexp
	}{
		{
			description:     "WhenResponseIsNil_ShouldNotLog",
			writer:          lit.NewRecorder(httptest.NewRecorder()),
			response:        nil,
			expectedContent: regexp.MustCompile("^$"),
		},
		{
			description: "WhenStatusCodeIs2xx_ShouldLogGreen",
			writer:      lit.NewRecorder(httptest.NewRecorder()),
			response: lit.ResponseFunc(func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("log"))
			}),
			expectedContent: regexp.MustCompile(
				"^\n\u001B\\[97;1;42m>> GET /users\u001B\\[0m\n> 200 OK\n> Start Time: .+\n> Remote Address: .+\n> Duration: .+\n> Content-Length: 3\n$",
			),
		},
		{
			description: "WhenStatusCodeIs3xx_ShouldLogBlue",
			writer:      lit.NewRecorder(httptest.NewRecorder()),
			response: lit.ResponseFunc(func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusPermanentRedirect)
				w.Write([]byte("log"))
			}),
			expectedContent: regexp.MustCompile(
				"^\n\u001B\\[97;1;104m>> GET /users\u001B\\[0m\n> 308 Permanent Redirect\n> Start Time: .+\n> Remote Address: .+\n> Duration: .+\n> Content-Length: 3\n$",
			),
		},
		{
			description: "WhenStatusCodeIs4xx_ShouldLogYellow",
			writer:      lit.NewRecorder(httptest.NewRecorder()),
			response: lit.ResponseFunc(func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("log"))
			}),
			expectedContent: regexp.MustCompile(
				"^\n\u001B\\[97;1;43m>> GET /users\u001B\\[0m\n> 404 Not Found\n> Start Time: .+\n> Remote Address: .+\n> Duration: .+\n> Content-Length: 3\n$",
			),
		},
		{
			description: "WhenStatusCodeIs5xx_ShouldLogRed",
			writer:      lit.NewRecorder(httptest.NewRecorder()),
			response: lit.ResponseFunc(func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("log"))
			}),
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
				return test.response
			})

			var output bytes.Buffer
			log.SetOutput(&output)
			t.Cleanup(func() {
				log.SetOutput(os.Stderr)
			})

			// Act
			response := lit.Log(handler)(r)
			if response != nil {
				response.Write(httptest.NewRecorder())
			}

			// Assert
			require.Regexp(t, test.expectedContent, output.String())
		})
	}
}
