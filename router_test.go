package lit_test

import (
	"net/http"
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
