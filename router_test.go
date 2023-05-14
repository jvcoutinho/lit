package lit_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jvcoutinho/lit"
	"github.com/jvcoutinho/lit/internal/routes"
	"github.com/stretchr/testify/require"
)

func TestRouter_Handle(t *testing.T) {
	t.Parallel()

	defaultHandle := func(ctx *lit.Context) lit.Result { return nil }

	type Route struct {
		Pattern string
		Method  string
		Handle  lit.HandlerFunc
	}

	type TestCase struct {
		description string

		existingRoutes []Route

		pattern string
		method  string
		handle  lit.HandlerFunc

		panics     bool
		panicValue any
	}

	tests := []TestCase{
		{
			description:    "GivenHandleIsNil_ShouldPanic",
			existingRoutes: nil,
			pattern:        "/books",
			method:         http.MethodGet,
			handle:         nil,
			panics:         true,
			panicValue:     lit.ErrNilHandler,
		},
		{
			description:    "GivenPatternContainsDoubleSlashes_ShouldPanic",
			existingRoutes: nil,
			pattern:        "/books//",
			method:         http.MethodGet,
			handle:         defaultHandle,
			panics:         true,
			panicValue:     routes.ErrPatternContainsDoubleSlash,
		},
		{
			description:    "GivenMethodIsEmpty_ShouldPanic",
			existingRoutes: nil,
			pattern:        "/books",
			method:         "",
			handle:         defaultHandle,
			panics:         true,
			panicValue:     routes.ErrMethodIsEmpty,
		},
		{
			description:    "GivenRouterIsHandlingNoRoutes_ShouldNotPanic",
			existingRoutes: []Route{},
			pattern:        "/books",
			method:         http.MethodGet,
			handle:         defaultHandle,
			panics:         false,
			panicValue:     nil,
		},
		{
			description: "GivenRouteHasTrailingSlash_AndSameRouteWithoutTrailingExists_ShouldPanic",
			existingRoutes: []Route{
				{
					Pattern: "/books",
					Method:  http.MethodGet,
					Handle:  defaultHandle,
				},
			},
			pattern:    "/books/",
			method:     http.MethodGet,
			handle:     defaultHandle,
			panics:     true,
			panicValue: routes.ErrPatternHasBeenDefinedAlready,
		},
		{
			description: "GivenRouteHasOnlyStaticSegments_AndItExists_ShouldPanic",
			existingRoutes: []Route{
				{
					Pattern: "/books",
					Method:  http.MethodGet,
					Handle:  defaultHandle,
				},
			},
			pattern:    "/books",
			method:     http.MethodGet,
			handle:     defaultHandle,
			panics:     true,
			panicValue: routes.ErrPatternHasBeenDefinedAlready,
		},
		{
			description: "GivenRouteHasOnlyStaticSegments_AndItDoesNotExists_PatternDiffers_ShouldNotPanic",
			existingRoutes: []Route{
				{
					Pattern: "/books",
					Method:  http.MethodGet,
					Handle:  defaultHandle,
				},
			},
			pattern:    "/users",
			method:     http.MethodGet,
			handle:     defaultHandle,
			panics:     false,
			panicValue: nil,
		},
		{
			description: "GivenRouteHasOnlyStaticSegments_AndItDoesNotExists_MethodDiffers_ShouldNotPanic",
			existingRoutes: []Route{
				{
					Pattern: "/books",
					Method:  http.MethodGet,
					Handle:  defaultHandle,
				},
			},
			pattern:    "/books",
			method:     http.MethodPost,
			handle:     defaultHandle,
			panics:     false,
			panicValue: nil,
		},
		{
			description: "GivenRouteHasOnlyStaticSegments_AndItDoesNotExists_IsSubpattern_ShouldNotPanic",
			existingRoutes: []Route{
				{
					Pattern: "/books/book",
					Method:  http.MethodGet,
					Handle:  defaultHandle,
				},
			},
			pattern:    "/books",
			method:     http.MethodGet,
			handle:     defaultHandle,
			panics:     false,
			panicValue: nil,
		},
		{
			description: "GivenRouteHasOnlyStaticSegments_AndItDoesNotExists_IsSuperpattern_ShouldNotPanic",
			existingRoutes: []Route{
				{
					Pattern: "/books",
					Method:  http.MethodGet,
					Handle:  defaultHandle,
				},
			},
			pattern:    "/books/book",
			method:     http.MethodGet,
			handle:     defaultHandle,
			panics:     false,
			panicValue: nil,
		},
		{
			description: "GivenRouteHasDynamicSegments_AndItExists_ParametersAreEqual_ShouldPanic",
			existingRoutes: []Route{
				{
					Pattern: "/books/:id",
					Method:  http.MethodGet,
					Handle:  defaultHandle,
				},
			},
			pattern:    "/books/:id",
			method:     http.MethodGet,
			handle:     defaultHandle,
			panics:     true,
			panicValue: routes.ErrPatternHasBeenDefinedAlready,
		},
		{
			description: "GivenRouteHasDynamicSegments_AndItExists_ParametersAreDifferent_ShouldPanic",
			existingRoutes: []Route{
				{
					Pattern: "/books/:id",
					Method:  http.MethodGet,
					Handle:  defaultHandle,
				},
			},
			pattern:    "/books/:user_id",
			method:     http.MethodGet,
			handle:     defaultHandle,
			panics:     true,
			panicValue: routes.ErrPatternHasConflictingParameters,
		},
		{
			description: "GivenRouteHasDynamicSegments_AndItDoesNotExist_ParametersAreDifferent_MethodDiffers_ShouldNotPanic",
			existingRoutes: []Route{
				{
					Pattern: "/books/:id",
					Method:  http.MethodGet,
					Handle:  defaultHandle,
				},
			},
			pattern:    "/books/:user_id",
			method:     http.MethodPost,
			handle:     defaultHandle,
			panics:     false,
			panicValue: nil,
		},
		{
			description: "GivenRouteHasDynamicSegments_AndItDoesNotExist_ParametersAreEqual_MethodDiffers_ShouldNotPanic",
			existingRoutes: []Route{
				{
					Pattern: "/books/:id",
					Method:  http.MethodGet,
					Handle:  defaultHandle,
				},
			},
			pattern:    "/books/:id",
			method:     http.MethodPost,
			handle:     defaultHandle,
			panics:     false,
			panicValue: nil,
		},
		{
			description: "GivenRouteHasDynamicSegments_AndItDoesNotExist_ParametersAreDifferent_IsSubpattern_ShouldNotPanic",
			existingRoutes: []Route{
				{
					Pattern: "/books/:id/readers",
					Method:  http.MethodGet,
					Handle:  defaultHandle,
				},
			},
			pattern:    "/books/:user_id",
			method:     http.MethodGet,
			handle:     defaultHandle,
			panics:     false,
			panicValue: nil,
		},
		{
			description: "GivenRouteHasDynamicSegments_AndItDoesNotExist_ParametersAreEqual_IsSubpattern_ShouldNotPanic",
			existingRoutes: []Route{
				{
					Pattern: "/books/:id/readers",
					Method:  http.MethodGet,
					Handle:  defaultHandle,
				},
			},
			pattern:    "/books/:id",
			method:     http.MethodGet,
			handle:     defaultHandle,
			panics:     false,
			panicValue: nil,
		},
		{
			description: "GivenRouteHasDynamicSegments_AndItDoesNotExist_ParametersAreDifferent_IsSuperpattern_ShouldNotPanic",
			existingRoutes: []Route{
				{
					Pattern: "/books/:id",
					Method:  http.MethodGet,
					Handle:  defaultHandle,
				},
			},
			pattern:    "/books/:user_id/readers",
			method:     http.MethodGet,
			handle:     defaultHandle,
			panics:     false,
			panicValue: nil,
		},
		{
			description: "GivenRouteHasDynamicSegments_AndItDoesNotExist_ParametersAreEqual_IsSuperpattern_ShouldNotPanic",
			existingRoutes: []Route{
				{
					Pattern: "/books/:id",
					Method:  http.MethodGet,
					Handle:  defaultHandle,
				},
			},
			pattern:    "/books/:id/readers",
			method:     http.MethodGet,
			handle:     defaultHandle,
			panics:     false,
			panicValue: nil,
		},
		{
			description: "GivenRouteHasOnlyDynamicSegments_AndItDoesNotExist_StaticSegmentsExist_ShouldNotPanic",
			existingRoutes: []Route{
				{
					Pattern: "/books/:id",
					Method:  http.MethodGet,
					Handle:  defaultHandle,
				},
			},
			pattern:    "/:user_id/:book_id",
			method:     http.MethodGet,
			handle:     defaultHandle,
			panics:     false,
			panicValue: nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Arrange
			router := lit.NewRouter()

			for _, route := range test.existingRoutes {
				router.Handle(route.Pattern, route.Method, route.Handle)
			}

			// Act
			// Assert
			if test.panics {
				require.PanicsWithValue(t, test.panicValue, func() {
					router.Handle(test.pattern, test.method, test.handle)
				})
			} else {
				require.NotPanics(t, func() {
					router.Handle(test.pattern, test.method, test.handle)
				})
			}
		})
	}
}

func TestRouter_ServeHTTP(t *testing.T) {
	t.Parallel()

	type Route struct {
		Pattern string
		Method  string
		Handle  lit.HandlerFunc
	}

	type TestCase struct {
		description          string
		setUpRouter          func(*lit.Router)
		existingRoutes       []Route
		path                 string
		method               string
		expectedArguments    map[string]string
		expectedResponseBody string
		expectedStatusCode   int
		expectedHeader       http.Header
	}

	tests := []TestCase{
		{
			description:          "GivenRouterHandlesNoRoutes_ShouldWrite404NotFoundResponse",
			setUpRouter:          func(r *lit.Router) {},
			existingRoutes:       []Route{},
			path:                 "/users",
			method:               http.MethodGet,
			expectedArguments:    map[string]string{},
			expectedResponseBody: "404 page not found\n",
			expectedStatusCode:   http.StatusNotFound,
			expectedHeader: http.Header{
				"Content-Type":           {"text/plain; charset=utf-8"},
				"X-Content-Type-Options": {"nosniff"},
			},
		},
		{
			description: "GivenRouteDoesNotExist_MethodDiffers_ShouldWrite404NotFoundResponse",
			setUpRouter: func(r *lit.Router) {},
			existingRoutes: []Route{
				{
					Pattern: "/users",
					Method:  http.MethodGet,
					Handle:  func(ctx *lit.Context) lit.Result { return nil },
				},
			},
			path:                 "/users",
			method:               http.MethodPost,
			expectedArguments:    map[string]string{},
			expectedResponseBody: "404 page not found\n",
			expectedStatusCode:   http.StatusNotFound,
			expectedHeader: http.Header{
				"Content-Type":           {"text/plain; charset=utf-8"},
				"X-Content-Type-Options": {"nosniff"},
			},
		},
		{
			description: "GivenRouteDoesNotExist_PatternDiffers_ShouldWrite404NotFoundResponse",
			setUpRouter: func(r *lit.Router) {},
			existingRoutes: []Route{
				{
					Pattern: "/users",
					Method:  http.MethodGet,
					Handle:  func(ctx *lit.Context) lit.Result { return nil },
				},
			},
			path:                 "/books",
			method:               http.MethodGet,
			expectedArguments:    map[string]string{},
			expectedResponseBody: "404 page not found\n",
			expectedStatusCode:   http.StatusNotFound,
			expectedHeader: http.Header{
				"Content-Type":           {"text/plain; charset=utf-8"},
				"X-Content-Type-Options": {"nosniff"},
			},
		},
		{
			description: "GivenRouteExists_AndRoutePatternHasStaticSegmentsOnly_ShouldWrite200OKResponse",
			setUpRouter: func(r *lit.Router) {},
			existingRoutes: []Route{
				{
					Pattern: "/users",
					Method:  http.MethodGet,
					Handle:  func(ctx *lit.Context) lit.Result { return nil },
				},
			},
			path:                 "/users",
			method:               http.MethodGet,
			expectedArguments:    map[string]string{},
			expectedResponseBody: "",
			expectedStatusCode:   http.StatusOK,
			expectedHeader:       http.Header{},
		},
		{
			description: "GivenRouteExists_AndRoutePatternHasDynamicSegments_ShouldWrite200OKResponse_AndBindArguments",
			setUpRouter: func(r *lit.Router) {},
			existingRoutes: []Route{
				{
					Pattern: "/users/:id",
					Method:  http.MethodGet,
					Handle:  func(ctx *lit.Context) lit.Result { return nil },
				},
			},
			path:                 "/users/Bob",
			method:               http.MethodGet,
			expectedArguments:    map[string]string{":id": "Bob"},
			expectedResponseBody: "",
			expectedStatusCode:   http.StatusOK,
			expectedHeader:       http.Header{},
		},
		{
			description: "GivenRouteExists_AndRoutePatternHasMultipleDynamicSegments_ShouldWrite200OKResponse_AndBindArguments",
			setUpRouter: func(r *lit.Router) {},
			existingRoutes: []Route{
				{
					Pattern: "/users/:user_id/books/:book_id",
					Method:  http.MethodGet,
					Handle:  func(ctx *lit.Context) lit.Result { return nil },
				},
			},
			path:                 "/users/Bob/books/123",
			method:               http.MethodGet,
			expectedArguments:    map[string]string{":user_id": "Bob", ":book_id": "123"},
			expectedResponseBody: "",
			expectedStatusCode:   http.StatusOK,
			expectedHeader:       http.Header{},
		},
		{
			description: "GivenRouteExists_AndPathHasTrailingSlash_ShouldHandleRouteWithoutTrailingSlash",
			setUpRouter: func(r *lit.Router) {},
			existingRoutes: []Route{
				{
					Pattern: "/users",
					Method:  http.MethodGet,
					Handle: func(ctx *lit.Context) lit.Result {
						return nil
					},
				},
				{
					Pattern: "/users/:id",
					Method:  http.MethodGet,
					Handle: func(ctx *lit.Context) lit.Result {
						ctx.SetStatusCode(http.StatusOK)

						return nil
					},
				},
			},
			path:                 "/users",
			method:               http.MethodGet,
			expectedArguments:    map[string]string{},
			expectedResponseBody: "",
			expectedStatusCode:   http.StatusOK,
			expectedHeader:       http.Header{},
		},
		{
			description: "GivenRouteWritesToContext_ShouldReflectInResponse",
			setUpRouter: func(r *lit.Router) {},
			existingRoutes: []Route{
				{
					Pattern: "/users",
					Method:  http.MethodGet,
					Handle: func(ctx *lit.Context) lit.Result {
						ctx.SetStatusCode(http.StatusBadRequest)
						ctx.WriteBody([]byte("body"))
						ctx.SetHeader("Content-Type", "application/json")

						return nil
					},
				},
			},
			path:                 "/users",
			method:               http.MethodGet,
			expectedArguments:    map[string]string{},
			expectedResponseBody: "body",
			expectedStatusCode:   http.StatusBadRequest,
			expectedHeader: http.Header{
				"Content-Type": {"application/json"},
			},
		},
		{
			description: "GivenRouteDoesNotExist_AndNotFoundHandlerIsSet_ShouldRunHandler",
			setUpRouter: func(r *lit.Router) {
				r.SetNotFoundHandler(func(ctx *lit.Context) lit.Result {
					ctx.SetStatusCode(http.StatusNotFound)
					ctx.WriteBody([]byte("not found"))

					return nil
				})
			},
			existingRoutes: []Route{
				{
					Pattern: "/users",
					Method:  http.MethodGet,
					Handle:  func(ctx *lit.Context) lit.Result { return nil },
				},
			},
			path:                 "/books",
			method:               http.MethodGet,
			expectedArguments:    map[string]string{},
			expectedResponseBody: "not found",
			expectedStatusCode:   http.StatusNotFound,
			expectedHeader:       http.Header{},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			testArgumentsHandler := func(route Route) func(ctx *lit.Context) lit.Result {
				return func(ctx *lit.Context) lit.Result {
					require.Equal(t, test.expectedArguments, ctx.URIArguments())

					return route.Handle(ctx)
				}
			}

			// Arrange
			router := lit.NewRouter()
			test.setUpRouter(router)

			for _, route := range test.existingRoutes {
				router.Handle(route.Pattern, route.Method, testArgumentsHandler(route))
			}

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(test.method, test.path, nil)

			// Act
			router.ServeHTTP(recorder, request)

			// Assert
			require.Equal(t, test.expectedResponseBody, recorder.Body.String())
			require.Equal(t, test.expectedStatusCode, recorder.Code)
			require.Equal(t, test.expectedHeader, recorder.Header())
		})
	}
}

func TestRouter_SetNotFoundHandler(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		description     string
		notFoundHandler lit.HandlerFunc
		panics          bool
		panicValue      any
	}

	tests := []TestCase{
		{
			description:     "GivenHandlerIsNil_ShouldPanic",
			notFoundHandler: nil,
			panics:          true,
			panicValue:      lit.ErrNilHandler,
		},
		{
			description:     "GivenHandlerIsNotNil_ShouldNotPanic",
			notFoundHandler: func(ctx *lit.Context) lit.Result { return nil },
			panics:          false,
			panicValue:      nil,
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
			if test.panics {
				require.PanicsWithValue(t, test.panicValue, func() {
					router.SetNotFoundHandler(test.notFoundHandler)
				})
			} else {
				require.NotPanics(t, func() {
					router.SetNotFoundHandler(test.notFoundHandler)
				})
			}
		})
	}
}
