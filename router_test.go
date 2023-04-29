package lit_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jvcoutinho/lit"
	"github.com/stretchr/testify/require"
)

func TestRouter_Handle(t *testing.T) {
	t.Parallel()

	type route struct {
		Pattern string
		Method  string
		Handler lit.HandleFunc
	}

	tests := []struct {
		name string

		currentRoutes []route
		routeToHandle route

		panics        bool
		expectedError string
	}{
		{
			name: "RouteDoesNotExist_SameMethod",
			currentRoutes: []route{
				{Pattern: "/users", Method: http.MethodGet, Handler: func(ctx *lit.Context) {}},
			},
			routeToHandle: route{Pattern: "/users/owner", Method: http.MethodGet, Handler: func(ctx *lit.Context) {}},
			panics:        false,
			expectedError: "",
		},
		{
			name: "RouteDoesNotExist_SamePattern",
			currentRoutes: []route{
				{Pattern: "/users", Method: http.MethodGet, Handler: func(ctx *lit.Context) {}},
			},
			routeToHandle: route{Pattern: "/users", Method: http.MethodPost, Handler: func(ctx *lit.Context) {}},
			panics:        false,
			expectedError: "",
		},
		{
			name: "RouteAlreadyExists",
			currentRoutes: []route{
				{Pattern: "/users", Method: http.MethodGet, Handler: func(ctx *lit.Context) {}},
			},
			routeToHandle: route{Pattern: "/users", Method: http.MethodGet, Handler: func(ctx *lit.Context) {}},
			panics:        true,
			expectedError: "GET /users has been already defined",
		},
		{
			name: "RouteAlreadyExists_DifferentMethodCase",
			currentRoutes: []route{
				{Pattern: "/users", Method: http.MethodGet, Handler: func(ctx *lit.Context) {}},
			},
			routeToHandle: route{Pattern: "/users", Method: "get", Handler: func(ctx *lit.Context) {}},
			panics:        true,
			expectedError: "GET /users has been already defined",
		},
		{
			name: "RouteAlreadyExists_MissingLeadingSlash",
			currentRoutes: []route{
				{Pattern: "/users", Method: http.MethodGet, Handler: func(ctx *lit.Context) {}},
			},
			routeToHandle: route{Pattern: "users", Method: http.MethodGet, Handler: func(ctx *lit.Context) {}},
			panics:        true,
			expectedError: "GET /users has been already defined",
		},
		{
			name: "RouteAlreadyExists_PresentTrailingSlash",
			currentRoutes: []route{
				{Pattern: "/users", Method: http.MethodGet, Handler: func(ctx *lit.Context) {}},
			},
			routeToHandle: route{Pattern: "users/", Method: http.MethodGet, Handler: func(ctx *lit.Context) {}},
			panics:        true,
			expectedError: "GET /users has been already defined",
		},
		{
			name: "RouteAlreadyExists_MultiplePaths",
			currentRoutes: []route{
				{Pattern: "/users/owner", Method: http.MethodGet, Handler: func(ctx *lit.Context) {}},
			},
			routeToHandle: route{Pattern: "/users/owner", Method: http.MethodGet, Handler: func(ctx *lit.Context) {}},
			panics:        true,
			expectedError: "GET /users/owner has been already defined",
		},
		{
			name: "RouteAlreadyExists_MultiplePaths_MissingLeadingSlash",
			currentRoutes: []route{
				{Pattern: "/users/owner", Method: http.MethodGet, Handler: func(ctx *lit.Context) {}},
			},
			routeToHandle: route{Pattern: "users/owner", Method: http.MethodGet, Handler: func(ctx *lit.Context) {}},
			panics:        true,
			expectedError: "GET /users/owner has been already defined",
		},
		{
			name: "RouteAlreadyExists_MultiplePaths_PresentTrailingSlash",
			currentRoutes: []route{
				{Pattern: "/users/owner", Method: http.MethodGet, Handler: func(ctx *lit.Context) {}},
			},
			routeToHandle: route{Pattern: "users/owner/", Method: http.MethodGet, Handler: func(ctx *lit.Context) {}},
			panics:        true,
			expectedError: "GET /users/owner has been already defined",
		},
		{
			name: "RouteAlreadyExists_MultiplePaths_DifferentArguments",
			currentRoutes: []route{
				{Pattern: "/users/:id", Method: http.MethodGet, Handler: func(ctx *lit.Context) {}},
			},
			routeToHandle: route{Pattern: "/users/:user_id", Method: http.MethodGet, Handler: func(ctx *lit.Context) {}},
			panics:        true,
			expectedError: "GET /users/:user_id has been already defined",
		},
		{
			name: "RouteAlreadyExists_MultiplePaths_DifferentArguments_ArgumentInMiddle",
			currentRoutes: []route{
				{Pattern: "/users/:user_id/items", Method: http.MethodGet, Handler: func(ctx *lit.Context) {}},
			},
			routeToHandle: route{Pattern: "/users/:id/items", Method: http.MethodGet, Handler: func(ctx *lit.Context) {}},
			panics:        true,
			expectedError: "GET /users/:id/items has been already defined",
		},
		{
			name:          "InvalidRoute_DuplicateArgument",
			currentRoutes: nil,
			routeToHandle: route{Pattern: "/users/:id/items/:id", Method: http.MethodGet, Handler: func(ctx *lit.Context) {}},
			panics:        true,
			expectedError: "a pattern can not contain two arguments with the same name (:id)",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			r := lit.NewRouter()

			for _, currentRoute := range test.currentRoutes {
				r.Handle(currentRoute.Pattern, currentRoute.Method, currentRoute.Handler)
			}

			// Act
			// Assert
			if test.panics {
				require.PanicsWithError(t, test.expectedError, func() {
					r.Handle(test.routeToHandle.Pattern, test.routeToHandle.Method, test.routeToHandle.Handler)
				})
			} else {
				require.NotPanics(t, func() {
					r.Handle(test.routeToHandle.Pattern, test.routeToHandle.Method, test.routeToHandle.Handler)
				})
			}
		})
	}
}

func TestRouter_ServeHTTP(t *testing.T) {
	t.Parallel()

	type route struct {
		Pattern string
		Method  string
		Handler lit.HandleFunc
	}

	tests := []struct {
		name string

		currentRoutes   []route
		incomingMethod  string
		incomingPattern string

		expectedResponse   string
		expectedStatusCode int
	}{
		{
			name: "RouteNotDefined_DifferentMethod",
			currentRoutes: []route{
				{Pattern: "/users", Method: http.MethodGet, Handler: func(ctx *lit.Context) {}},
			},
			incomingMethod:     http.MethodPost,
			incomingPattern:    "/users",
			expectedResponse:   "404 page not found\n",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "RouteNotDefined_DifferentPattern",
			currentRoutes: []route{
				{Pattern: "/users", Method: http.MethodGet, Handler: func(ctx *lit.Context) {}},
			},
			incomingMethod:     http.MethodGet,
			incomingPattern:    "/books",
			expectedResponse:   "404 page not found\n",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "RouteNotDefined_Subpattern",
			currentRoutes: []route{
				{Pattern: "/users", Method: http.MethodGet, Handler: func(ctx *lit.Context) {}},
			},
			incomingMethod:     http.MethodGet,
			incomingPattern:    "/users/:user_id",
			expectedResponse:   "404 page not found\n",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "RouteNotDefined_Superpattern",
			currentRoutes: []route{
				{Pattern: "/users/:user_id", Method: http.MethodGet, Handler: func(ctx *lit.Context) {}},
			},
			incomingMethod:     http.MethodGet,
			incomingPattern:    "/users",
			expectedResponse:   "404 page not found\n",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "RouteDefined_EmptyHandler",
			currentRoutes: []route{
				{Pattern: "/users", Method: http.MethodGet, Handler: func(ctx *lit.Context) {}},
			},
			incomingMethod:     http.MethodGet,
			incomingPattern:    "/users",
			expectedResponse:   "",
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "RouteDefined_ArgumentSubstitution",
			currentRoutes: []route{
				{Pattern: "/users/:user_id", Method: http.MethodGet, Handler: func(ctx *lit.Context) {}},
			},
			incomingMethod:     http.MethodGet,
			incomingPattern:    "/users/123",
			expectedResponse:   "",
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			r := lit.NewRouter()

			for _, currentRoute := range test.currentRoutes {
				r.Handle(currentRoute.Pattern, currentRoute.Method, currentRoute.Handler)
			}

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(test.incomingMethod, test.incomingPattern, nil)

			// Act
			r.ServeHTTP(recorder, request)

			// Assert
			require.Equal(t, test.expectedResponse, recorder.Body.String())
			require.Equal(t, test.expectedStatusCode, recorder.Code)
		})
	}
}
