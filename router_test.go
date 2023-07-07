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
		setupRouter func(*lit.Router)
		pattern     string
		method      string
		handler     lit.HandlerFunc
		panicValue  string
	}{
		{
			description: "WhenHandlerIsNil_ShouldPanic",
			setupRouter: func(router *lit.Router) {},
			pattern:     "/users",
			method:      http.MethodGet,
			handler:     nil,
			panicValue:  "handler should not be nil",
		},
		{
			description: "WhenMethodIsEmpty_ShouldPanic",
			setupRouter: func(router *lit.Router) {},
			pattern:     "/users",
			method:      "",
			handler:     func(r *lit.Request) lit.Response { return nil },
			panicValue:  "method should not be empty",
		},
		{
			description: "WhenPatternDoesNotStartWithSlash_ShouldPanic",
			setupRouter: func(router *lit.Router) {},
			pattern:     "users",
			method:      http.MethodGet,
			handler:     func(r *lit.Request) lit.Response { return nil },
			panicValue:  "pattern should start with a slash (/)",
		},
		{
			description: "WhenPatternContainsDoubleSlashes_ShouldPanic",
			setupRouter: func(router *lit.Router) {},
			pattern:     "//users",
			method:      http.MethodGet,
			handler:     func(r *lit.Request) lit.Response { return nil },
			panicValue:  "pattern should not contain double slashes (//)",
		},
		{
			description: "GivenPatternAndMethodHaveBeenDefinedAlready_ShouldPanic",
			setupRouter: func(router *lit.Router) {
				router.Handle("/users", http.MethodGet, func(r *lit.Request) lit.Response { return nil })
			},
			pattern:    "/users",
			method:     http.MethodGet,
			handler:    func(r *lit.Request) lit.Response { return nil },
			panicValue: "route already exists",
		},
		{
			description: "GivenPatternWithDifferentParametersAndMethodHaveBeenDefinedAlready_ShouldPanic",
			setupRouter: func(router *lit.Router) {
				router.Handle("/users/:id", http.MethodGet, func(r *lit.Request) lit.Response { return nil })
			},
			pattern:    "/users/:user_id",
			method:     http.MethodGet,
			handler:    func(r *lit.Request) lit.Response { return nil },
			panicValue: "parameters are conflicting with defined ones in another route",
		},
		{
			description: "GivenPatternDoesNotExist_ShouldNotPanic",
			setupRouter: func(router *lit.Router) {
				router.Handle("/users", http.MethodGet, func(r *lit.Request) lit.Response { return nil })
			},
			pattern:    "/users/:user_id",
			method:     http.MethodGet,
			handler:    func(r *lit.Request) lit.Response { return nil },
			panicValue: "",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Arrange
			router := lit.NewRouter()
			test.setupRouter(router)

			// Act
			// Assert
			if test.panicValue != "" {
				require.PanicsWithError(t, test.panicValue, func() {
					router.Handle(test.pattern, test.method, test.handler)
				})
			} else {
				require.NotPanics(t, func() {
					router.Handle(test.pattern, test.method, test.handler)
				})
			}
		})
	}
}

func TestRouter_ServeHTTP(t *testing.T) {
	t.Parallel()

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
			description: "GivenRouteIsRegistered_ShouldReturnHandlerResponse_AndMatchArguments",
			setupRouter: func(router *lit.Router) {
				router.Handle("/users/:user_id/books/:book_id", http.MethodGet,
					func(r *lit.Request) lit.Response {
						return lit.ResponseFunc(func(writer http.ResponseWriter) error {
							_, _ = writer.Write([]byte("response"))

							args := r.URLArguments()
							require.Equal(t, "1", args[":user_id"])
							require.Equal(t, "2", args[":book_id"])

							return nil
						})
					})
			},
			uri:                "/users/1/books/2",
			method:             http.MethodGet,
			body:               nil,
			expectedResponse:   "response",
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
