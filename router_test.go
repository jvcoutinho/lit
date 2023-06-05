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
		description   string
		router        *lit.Router
		pattern       string
		method        string
		handler       lit.HandlerFunc
		expectedPanic any
	}{
		{
			description:   "WhenHandlerIsNil_ShouldPanic",
			router:        lit.NewRouter(),
			pattern:       "/users",
			method:        http.MethodGet,
			handler:       nil,
			expectedPanic: "handler should not be nil",
		},
		{
			description:   "WhenMethodIsEmpty_ShouldPanic",
			router:        lit.NewRouter(),
			pattern:       "/users",
			method:        "",
			handler:       func(r *lit.Request) lit.Response { return nil },
			expectedPanic: "method should not be empty",
		},
		{
			description:   "WhenPatternDoesNotStartWithSlash_ShouldPanic",
			router:        lit.NewRouter(),
			pattern:       "users",
			method:        http.MethodGet,
			handler:       func(r *lit.Request) lit.Response { return nil },
			expectedPanic: "pattern should start with a slash (/)",
		},
		{
			description:   "WhenPatternContainsDoubleSlashes_ShouldPanic",
			router:        lit.NewRouter(),
			pattern:       "//users",
			method:        http.MethodGet,
			handler:       func(r *lit.Request) lit.Response { return nil },
			expectedPanic: "pattern should not contain double slashes (//)",
		},
		{
			description: "GivenPatternAndMethodHaveBeenDefinedAlready_ShouldPanic",
			router: func() *lit.Router {
				router := lit.NewRouter()
				router.Handle("/users", http.MethodGet, func(r *lit.Request) lit.Response { return nil })
				return router
			}(),
			pattern:       "/users",
			method:        http.MethodGet,
			handler:       func(r *lit.Request) lit.Response { return nil },
			expectedPanic: "route already exists",
		},
		{
			description: "GivenPatternWithDifferentParametersAndMethodHaveBeenDefinedAlready_ShouldPanic",
			router: func() *lit.Router {
				router := lit.NewRouter()
				router.Handle("/users/:id", http.MethodGet, func(r *lit.Request) lit.Response { return nil })
				return router
			}(),
			pattern:       "/users/:user_id",
			method:        http.MethodGet,
			handler:       func(r *lit.Request) lit.Response { return nil },
			expectedPanic: "parameters are conflicting with defined ones in another route",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Arrange
			router := test.router

			// Act
			// Assert
			if test.expectedPanic != nil {
				require.PanicsWithValue(t, test.expectedPanic, func() {
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
		router             *lit.Router
		uri                string
		method             string
		body               io.Reader
		expectedResponse   string
		expectedStatusCode int
	}{
		{
			description:        "GivenRouteIsNotRegistered_ShouldReturnNotFoundResponse",
			router:             lit.NewRouter(),
			uri:                "/users",
			method:             http.MethodGet,
			body:               nil,
			expectedResponse:   "404 page not found\n",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			description: "GivenRouteIsRegistered_AndPatternDoesNotContainParameters_ShouldReturnHandlerResponse",
			router: func() *lit.Router {
				router := lit.NewRouter()
				router.Handle("/users", http.MethodGet, func(r *lit.Request) lit.Response {
					return lit.ResponseFunc(func(writer http.ResponseWriter) error {
						writer.Write([]byte("response"))
						return nil
					})
				})

				return router
			}(),
			uri:                "/users",
			method:             http.MethodGet,
			body:               nil,
			expectedResponse:   "response",
			expectedStatusCode: http.StatusOK,
		},
		{
			description: "GivenRouteIsRegistered_AndPatternContainsParameters_ShouldMatchArguments_AndReturnHandlerResponse",
			router: func() *lit.Router {
				router := lit.NewRouter()
				router.Handle("/users/:id", http.MethodGet, func(r *lit.Request) lit.Response {
					return lit.ResponseFunc(func(writer http.ResponseWriter) error {
						require.Equal(t, "123", r.URIArguments()[":id"])
						return nil
					})
				})

				return router
			}(),
			uri:                "/users/123",
			method:             http.MethodGet,
			body:               nil,
			expectedResponse:   "",
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Arrange
			router := test.router

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
