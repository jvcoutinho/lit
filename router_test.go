package lit_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jvcoutinho/lit"
	"github.com/stretchr/testify/require"
)

func TestRouter_Handle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		pattern     string
		method      string
		handler     lit.Handler
		panicValue  any
	}{
		{
			description: "WhenHandlerIsNil_ShouldPanic",
			pattern:     "/users",
			method:      http.MethodGet,
			handler:     nil,
			panicValue:  lit.ErrNilHandler,
		},
		{
			description: "WhenMethodIsEmpty_ShouldPanic",
			pattern:     "/users",
			method:      "",
			handler:     func(r *lit.Request) lit.Response { return nil },
			panicValue:  lit.ErrMethodIsEmpty,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Arrange
			router := lit.NewRouter()

			// Act
			// Assert
			require.PanicsWithValue(t, test.panicValue, func() {
				router.Handle(test.pattern, test.method, test.handler)
			})
		})
	}
}

func TestRouter_ServeHTTP(t *testing.T) {
	t.Parallel()

	testHandler := func(res string, arguments map[string]string) lit.Handler {
		return func(r *lit.Request) lit.Response {
			return lit.ResponseFunc(func(w http.ResponseWriter) {
				_, _ = w.Write([]byte(res))

				require.Equal(t, arguments, r.URIParameters())
			})
		}
	}

	tests := []struct {
		description        string
		setupRouter        func(*lit.Router)
		uri                string
		method             string
		body               io.Reader
		expectedResponse   string
		expectedStatusCode int
	}{
		{
			description:        "GivenRouteIsNotRegistered_ShouldReturnNotFoundResponse",
			setupRouter:        func(*lit.Router) {},
			uri:                "/users",
			method:             http.MethodGet,
			body:               nil,
			expectedResponse:   "404 page not found\n",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			description: "GivenRouteIsRegistered_ShouldReturnHandlerResponse",
			setupRouter: func(router *lit.Router) {
				router.Handle("/users/:user_id/books/:book_id", http.MethodGet,
					testHandler("test", map[string]string{"user_id": "1", "book_id": "2"}))
			},
			uri:                "/users/1/books/2",
			method:             http.MethodGet,
			body:               nil,
			expectedResponse:   "test",
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Arrange
			router := lit.NewRouter()
			test.setupRouter(router)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(test.method, test.uri, test.body)

			// Act
			router.ServeHTTP(recorder, request)

			// Assert
			require.Equal(t, test.expectedResponse, recorder.Body.String())
			require.Equal(t, test.expectedStatusCode, recorder.Code)
		})
	}
}
