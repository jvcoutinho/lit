package routes_test

import (
	"net/http"
	"testing"

	"github.com/jvcoutinho/lit/internal/routes"
	"github.com/stretchr/testify/require"
)

func TestGraph_CanBeInserted(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		description   string
		initialRoutes []routes.Route
		route         routes.Route
		expectedOk    bool
		expectedError error
	}

	tests := []TestCase{
		{
			description:   "GivenRouteContainsDuplicateParameters_ShouldReturnDuplicateArgumentsError",
			initialRoutes: []routes.Route{},
			route:         routes.NewRoute("/users/:id/books/:id", http.MethodGet),
			expectedOk:    false,
			expectedError: routes.DuplicateArgumentsError{Duplicate: ":id"},
		},
		{
			description:   "GivenGraphIsEmpty_ShouldReturnTrue",
			initialRoutes: []routes.Route{},
			route:         routes.NewRoute("/users", http.MethodGet),
			expectedOk:    true,
			expectedError: nil,
		},
		{
			description: "GivenRoutePatternDoesNotExistInGraph_AndMethodAlsoDoesNotExistInGraph_ShouldReturnTrue",
			initialRoutes: []routes.Route{
				routes.NewRoute("/books", http.MethodGet),
			},
			route:         routes.NewRoute("/users", http.MethodPost),
			expectedOk:    true,
			expectedError: nil,
		},
		{
			description: "GivenRoutePatternDoesNotExistInGraph_ButMethodDoes_ShouldReturnTrue",
			initialRoutes: []routes.Route{
				routes.NewRoute("/books", http.MethodGet),
			},
			route:         routes.NewRoute("/users", http.MethodGet),
			expectedOk:    true,
			expectedError: nil,
		},
		{
			description: "GivenRouteMethodDoesNotExistInGraph_ButPatternDoes_ShouldReturnTrue",
			initialRoutes: []routes.Route{
				routes.NewRoute("/users", http.MethodGet),
			},
			route:         routes.NewRoute("/users", http.MethodPost),
			expectedOk:    true,
			expectedError: nil,
		},
		{
			description: "GivenRouteMethodDoesNotExistInGraph_ButPatternDoes_AndPatternHasArguments_ShouldReturnTrue",
			initialRoutes: []routes.Route{
				routes.NewRoute("/users/:user_id", http.MethodGet),
			},
			route:         routes.NewRoute("/users/:id", http.MethodPost),
			expectedOk:    true,
			expectedError: nil,
		},
		{
			description: "GivenRouteParametersConflictWithExistingRouteInGraph_ShouldReturnRouteAlreadyDefinedError",
			initialRoutes: []routes.Route{
				routes.NewRoute("/users/:id", http.MethodGet),
			},
			route:         routes.NewRoute("/users/:user_id", http.MethodGet),
			expectedOk:    false,
			expectedError: routes.RouteAlreadyDefinedError{Route: routes.NewRoute("/users/:user_id", http.MethodGet)},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Arrange
			graph := routes.NewGraph(test.initialRoutes...)

			// Act
			actualOk, actualError := graph.CanBeInserted(test.route)

			// Assert
			require.Equal(t, test.expectedOk, actualOk)
			require.Equal(t, test.expectedError, actualError)
		})
	}
}

func TestGraph_MatchRoute(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		description   string
		initialRoutes []routes.Route
		route         routes.Route
		expectedRoute routes.Route
		expectedOk    bool
	}

	tests := []TestCase{
		{
			description:   "GivenGraphIsEmpty_ShouldReturnFalse",
			initialRoutes: []routes.Route{},
			route:         routes.NewRoute("/users", http.MethodGet),
			expectedRoute: routes.Route{},
			expectedOk:    false,
		},
		{
			description: "GivenRoutePatternDoesNotExistInGraph_ShouldReturnFalse",
			initialRoutes: []routes.Route{
				routes.NewRoute("/books", http.MethodGet),
			},
			route:         routes.NewRoute("/users", http.MethodGet),
			expectedRoute: routes.Route{},
			expectedOk:    false,
		},
		{
			description: "GivenRouteMethodDoesNotExistInGraph_ShouldReturnFalse",
			initialRoutes: []routes.Route{
				routes.NewRoute("/users", http.MethodGet),
			},
			route:         routes.NewRoute("/users", http.MethodPost),
			expectedRoute: routes.Route{},
			expectedOk:    false,
		},
		{
			description: "GivenRoutePatternDoesNotExistInGraph_ButASuperpatternDoes_ShouldReturnFalse",
			initialRoutes: []routes.Route{
				routes.NewRoute("/users/:user_id", http.MethodGet),
			},
			route:         routes.NewRoute("/users", http.MethodGet),
			expectedRoute: routes.Route{},
			expectedOk:    false,
		},
		{
			description: "GivenRoutePatternDoesNotExistInGraph_ButASubpatternDoes_ShouldReturnFalse",
			initialRoutes: []routes.Route{
				routes.NewRoute("/users", http.MethodGet),
			},
			route:         routes.NewRoute("/users/:user_id", http.MethodGet),
			expectedRoute: routes.Route{},
			expectedOk:    false,
		},
		{
			description: "GivenRoutePatternExistsInGraph_AndRouteMethodAlsoDoes_ShouldReturnMatch",
			initialRoutes: []routes.Route{
				routes.NewRoute("/users", http.MethodGet),
			},
			route:         routes.NewRoute("/users", http.MethodGet),
			expectedRoute: routes.NewRoute("/users", http.MethodGet),
			expectedOk:    true,
		},
		{
			description: "GivenRoutePatternExistsInGraphWithParameter_AndParameterIsInTheEndOfPattern_ShouldReturnMatch",
			initialRoutes: []routes.Route{
				routes.NewRoute("/users/:user_id", http.MethodGet),
			},
			route:         routes.NewRoute("/users/123", http.MethodGet),
			expectedRoute: routes.NewRoute("/users/:user_id", http.MethodGet),
			expectedOk:    true,
		},
		{
			description: "GivenRoutePatternExistsInGraphWithParameter_AndParameterIsInTheMiddleOfPattern_ShouldReturnMatch",
			initialRoutes: []routes.Route{
				routes.NewRoute("/users/:user_id/books", http.MethodGet),
			},
			route:         routes.NewRoute("/users/123/books", http.MethodGet),
			expectedRoute: routes.NewRoute("/users/:user_id/books", http.MethodGet),
			expectedOk:    true,
		},
		{
			description: "GivenRoutePatternExistsInGraphWithParameter_AndParameterIsInTheBeginningOfPattern_ShouldReturnMatch",
			initialRoutes: []routes.Route{
				routes.NewRoute("/:user_id/books", http.MethodGet),
			},
			route:         routes.NewRoute("/123/books", http.MethodGet),
			expectedRoute: routes.NewRoute("/:user_id/books", http.MethodGet),
			expectedOk:    true,
		},
		{
			description: "GivenRoutePatternDoesNotExistInGraph_ButSuperpatternWithParameterDoes_ShouldReturnFalse",
			initialRoutes: []routes.Route{
				routes.NewRoute("/users/:user_id/books", http.MethodGet),
			},
			route:         routes.NewRoute("/users/123", http.MethodGet),
			expectedRoute: routes.Route{},
			expectedOk:    false,
		},
		{
			description: "GivenRoutePatternDoesNotExistInGraph_ButSubpatternWithParameterDoes_ShouldReturnFalse",
			initialRoutes: []routes.Route{
				routes.NewRoute("/users/:user_id", http.MethodGet),
			},
			route:         routes.NewRoute("/users/123/books", http.MethodGet),
			expectedRoute: routes.Route{},
			expectedOk:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Arrange
			graph := routes.NewGraph(test.initialRoutes...)

			// Act
			match, actualOk := graph.MatchRoute(test.route)

			// Assert
			require.Equal(t, test.expectedRoute, match.MatchedRoute())
			require.Equal(t, test.expectedOk, actualOk)
		})
	}
}
