package lit_test

import (
	"net/http"
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

func TestRouter_WithNotFoundHandler(t *testing.T) {
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
					router.WithNotFoundHandler(test.notFoundHandler)
				})
			} else {
				require.NotPanics(t, func() {
					router.WithNotFoundHandler(test.notFoundHandler)
				})
			}
		})
	}
}
