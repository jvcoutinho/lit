package lit_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jvcoutinho/lit"
	"github.com/jvcoutinho/lit/littest"
	"github.com/jvcoutinho/lit/render"
	"github.com/stretchr/testify/require"
)

func TestRouter_Handle(t *testing.T) {
	t.Parallel()

	type route struct {
		Pattern string
		Method  string
		Handler lit.HandleFunc
	}

	emptyHandler := func(ctx *lit.Context) lit.Result { return nil }

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
				{Pattern: "/users", Method: http.MethodGet, Handler: emptyHandler},
			},
			routeToHandle: route{Pattern: "/users/owner", Method: http.MethodGet, Handler: emptyHandler},
			panics:        false,
			expectedError: "",
		},
		{
			name: "RouteDoesNotExist_SamePattern",
			currentRoutes: []route{
				{Pattern: "/users", Method: http.MethodGet, Handler: emptyHandler},
			},
			routeToHandle: route{Pattern: "/users", Method: http.MethodPost, Handler: emptyHandler},
			panics:        false,
			expectedError: "",
		},
		{
			name: "RouteDoesNotExist_Subpattern",
			currentRoutes: []route{
				{Pattern: "/users/:user_id", Method: http.MethodGet, Handler: emptyHandler},
			},
			routeToHandle: route{Pattern: "/users", Method: http.MethodGet, Handler: emptyHandler},
			panics:        false,
			expectedError: "",
		},
		{
			name: "RouteDoesNotExist_Superpattern",
			currentRoutes: []route{
				{Pattern: "/users", Method: http.MethodGet, Handler: emptyHandler},
			},
			routeToHandle: route{Pattern: "/users/:user_id", Method: http.MethodGet, Handler: emptyHandler},
			panics:        false,
			expectedError: "",
		},
		{
			name: "RouteAlreadyExists",
			currentRoutes: []route{
				{Pattern: "/users", Method: http.MethodGet, Handler: emptyHandler},
			},
			routeToHandle: route{Pattern: "/users", Method: http.MethodGet, Handler: emptyHandler},
			panics:        true,
			expectedError: "GET /users has been already defined",
		},
		{
			name: "RouteAlreadyExists_DifferentMethodCase",
			currentRoutes: []route{
				{Pattern: "/users", Method: http.MethodGet, Handler: emptyHandler},
			},
			routeToHandle: route{Pattern: "/users", Method: "get", Handler: emptyHandler},
			panics:        true,
			expectedError: "GET /users has been already defined",
		},
		{
			name: "RouteAlreadyExists_MissingLeadingSlash",
			currentRoutes: []route{
				{Pattern: "/users", Method: http.MethodGet, Handler: emptyHandler},
			},
			routeToHandle: route{Pattern: "users", Method: http.MethodGet, Handler: emptyHandler},
			panics:        true,
			expectedError: "GET /users has been already defined",
		},
		{
			name: "RouteAlreadyExists_PresentTrailingSlash",
			currentRoutes: []route{
				{Pattern: "/users", Method: http.MethodGet, Handler: emptyHandler},
			},
			routeToHandle: route{Pattern: "users/", Method: http.MethodGet, Handler: emptyHandler},
			panics:        true,
			expectedError: "GET /users has been already defined",
		},
		{
			name: "RouteAlreadyExists_MultiplePaths",
			currentRoutes: []route{
				{Pattern: "/users/owner", Method: http.MethodGet, Handler: emptyHandler},
			},
			routeToHandle: route{Pattern: "/users/owner", Method: http.MethodGet, Handler: emptyHandler},
			panics:        true,
			expectedError: "GET /users/owner has been already defined",
		},
		{
			name: "RouteAlreadyExists_MultiplePaths_MissingLeadingSlash",
			currentRoutes: []route{
				{Pattern: "/users/owner", Method: http.MethodGet, Handler: emptyHandler},
			},
			routeToHandle: route{Pattern: "users/owner", Method: http.MethodGet, Handler: emptyHandler},
			panics:        true,
			expectedError: "GET /users/owner has been already defined",
		},
		{
			name: "RouteAlreadyExists_MultiplePaths_PresentTrailingSlash",
			currentRoutes: []route{
				{Pattern: "/users/owner", Method: http.MethodGet, Handler: emptyHandler},
			},
			routeToHandle: route{Pattern: "users/owner/", Method: http.MethodGet, Handler: emptyHandler},
			panics:        true,
			expectedError: "GET /users/owner has been already defined",
		},
		{
			name: "RouteAlreadyExists_MultiplePaths_DifferentArguments",
			currentRoutes: []route{
				{Pattern: "/users/:id", Method: http.MethodGet, Handler: emptyHandler},
			},
			routeToHandle: route{Pattern: "/users/:user_id", Method: http.MethodGet, Handler: emptyHandler},
			panics:        true,
			expectedError: "GET /users/:user_id has been already defined",
		},
		{
			name: "RouteAlreadyExists_MultiplePaths_DifferentArguments_ArgumentInMiddle",
			currentRoutes: []route{
				{Pattern: "/users/:user_id/items", Method: http.MethodGet, Handler: emptyHandler},
			},
			routeToHandle: route{Pattern: "/users/:id/items", Method: http.MethodGet, Handler: emptyHandler},
			panics:        true,
			expectedError: "GET /users/:id/items has been already defined",
		},
		{
			name:          "InvalidRoute_DuplicateArgument",
			currentRoutes: nil,
			routeToHandle: route{Pattern: "/users/:id/items/:id", Method: http.MethodGet, Handler: emptyHandler},
			panics:        true,
			expectedError: "a pattern can not contain two arguments with the same name (:id)",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			router := lit.NewRouter()

			for _, currentRoute := range test.currentRoutes {
				router.Handle(currentRoute.Pattern, currentRoute.Method, currentRoute.Handler)
			}

			// Act
			// Assert
			if test.panics {
				require.PanicsWithError(t, test.expectedError, func() {
					router.Handle(test.routeToHandle.Pattern, test.routeToHandle.Method, test.routeToHandle.Handler)
				})
			} else {
				require.NotPanics(t, func() {
					router.Handle(test.routeToHandle.Pattern, test.routeToHandle.Method, test.routeToHandle.Handler)
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

	okHandler := func(_ *lit.Context) lit.Result { return nil }
	notFoundHandler := func(ctx *lit.Context) lit.Result {
		http.NotFound(ctx.ResponseWriter, ctx.Request)

		return nil
	}

	tests := []struct {
		name string

		currentRoutes   []route
		incomingMethod  string
		incomingPattern string

		expectedArguments  map[string]string
		expectedResponse   string
		expectedStatusCode int
	}{
		{
			name: "RouteNotDefined_DifferentMethod",
			currentRoutes: []route{
				{Pattern: "/users", Method: http.MethodGet, Handler: okHandler},
			},
			incomingMethod:     http.MethodPost,
			incomingPattern:    "/users",
			expectedArguments:  map[string]string{},
			expectedResponse:   "404 page not found\n",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "RouteNotDefined_DifferentPattern",
			currentRoutes: []route{
				{Pattern: "/users", Method: http.MethodGet, Handler: okHandler},
			},
			incomingMethod:     http.MethodGet,
			incomingPattern:    "/books",
			expectedArguments:  map[string]string{},
			expectedResponse:   "404 page not found\n",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "RouteDefined_Root",
			currentRoutes: []route{
				{Pattern: "/", Method: http.MethodGet, Handler: okHandler},
				{Pattern: "/:user_id", Method: http.MethodGet, Handler: okHandler},
			},
			incomingMethod:     http.MethodGet,
			incomingPattern:    "/",
			expectedArguments:  map[string]string{},
			expectedResponse:   "",
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "RouteDefined_ArgumentAtRoot",
			currentRoutes: []route{
				{Pattern: "/:user_id", Method: http.MethodGet, Handler: okHandler},
			},
			incomingMethod:  http.MethodGet,
			incomingPattern: "/",
			expectedArguments: map[string]string{
				":user_id": "",
			},
			expectedResponse:   "",
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "RouteNotDefined_Subpattern",
			currentRoutes: []route{
				{Pattern: "/users", Method: http.MethodGet, Handler: okHandler},
			},
			incomingMethod:     http.MethodGet,
			incomingPattern:    "/users/:user_id",
			expectedArguments:  map[string]string{},
			expectedResponse:   "404 page not found\n",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "RouteNotDefined_Superpattern",
			currentRoutes: []route{
				{Pattern: "/users/:user_id", Method: http.MethodGet, Handler: okHandler},
			},
			incomingMethod:     http.MethodGet,
			incomingPattern:    "/users",
			expectedArguments:  map[string]string{},
			expectedResponse:   "404 page not found\n",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "RouteDefined",
			currentRoutes: []route{
				{Pattern: "/users", Method: http.MethodGet, Handler: okHandler},
			},
			incomingMethod:     http.MethodGet,
			incomingPattern:    "/users",
			expectedArguments:  map[string]string{},
			expectedResponse:   "",
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "RouteDefined_ArgumentSubstitution",
			currentRoutes: []route{
				{Pattern: "/users/:user_id", Method: http.MethodGet, Handler: okHandler},
			},
			incomingMethod:  http.MethodGet,
			incomingPattern: "/users/123",
			expectedArguments: map[string]string{
				":user_id": "123",
			},
			expectedResponse:   "",
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "RouteDefined_ArgumentSubstitution_SameArguments",
			currentRoutes: []route{
				{Pattern: "/users/:user_id", Method: http.MethodGet, Handler: notFoundHandler},
				{Pattern: "/users/:user_id/books/:book_id", Method: http.MethodGet, Handler: okHandler},
			},
			incomingMethod:  http.MethodGet,
			incomingPattern: "/users/123/books/234",
			expectedArguments: map[string]string{
				":user_id": "123",
				":book_id": "234",
			},
			expectedResponse:   "",
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "RouteDefined_ArgumentSubstitution_Subpattern_DifferentArguments",
			currentRoutes: []route{
				{Pattern: "/users/:user_id", Method: http.MethodGet, Handler: notFoundHandler},
				{Pattern: "/users/:id/books/:book_id", Method: http.MethodGet, Handler: okHandler},
			},
			incomingMethod:  http.MethodGet,
			incomingPattern: "/users/123/books/234",
			expectedArguments: map[string]string{
				":id":      "123",
				":book_id": "234",
			},
			expectedResponse:   "",
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "RouteDefined_ArgumentSubstitution_Superpattern_DifferentArguments",
			currentRoutes: []route{
				{Pattern: "/users/:id/books/:book_id", Method: http.MethodGet, Handler: notFoundHandler},
				{Pattern: "/users/:user_id", Method: http.MethodGet, Handler: okHandler},
			},
			incomingMethod:  http.MethodGet,
			incomingPattern: "/users/123",
			expectedArguments: map[string]string{
				":user_id": "123",
			},
			expectedResponse:   "",
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "RouteDefined_RenderResponse",
			currentRoutes: []route{
				{
					Pattern: "/users",
					Method:  http.MethodGet,
					Handler: func(ctx *lit.Context) lit.Result {
						return render.Ok([]string{})
					},
				},
			},
			incomingMethod:     http.MethodGet,
			incomingPattern:    "/users",
			expectedArguments:  map[string]string{},
			expectedResponse:   "[]",
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "RouteDefined_RenderResponse_ResponseIsNotRenderable",
			currentRoutes: []route{
				{
					Pattern: "/users",
					Method:  http.MethodGet,
					Handler: func(ctx *lit.Context) lit.Result {
						return render.Ok(complex(1.0, 1.0))
					},
				},
			},
			incomingMethod:     http.MethodGet,
			incomingPattern:    "/users",
			expectedArguments:  map[string]string{},
			expectedResponse:   "rendering JSON: json: unsupported type: complex128\n",
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			router := lit.NewRouter()

			for _, currentRoute := range test.currentRoutes {
				router.Handle(currentRoute.Pattern, currentRoute.Method, func(ctx *lit.Context) lit.Result {
					require.Equal(t, test.expectedArguments, ctx.URIArguments())

					return currentRoute.Handler(ctx)
				})
			}

			// Act
			recorder := littest.Request(t, router, httptest.NewRequest(test.incomingMethod, test.incomingPattern, nil))

			// Assert
			require.Equal(t, test.expectedResponse, recorder.Body.String())
			require.Equal(t, test.expectedStatusCode, recorder.Code)
		})
	}
}
