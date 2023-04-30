package routes_test

import (
	"net/http"
	"testing"

	"github.com/jvcoutinho/lit/internal/routes"
	"github.com/stretchr/testify/require"
)

func TestGraph_CanBeInserted(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name string

		currentRoutes []routes.Route
		routeToCheck  routes.Route

		expectedError error
		expectedOk    bool
	}

	tests := []testCase{
		{
			name:          "EmptyGraph",
			currentRoutes: []routes.Route{},
			routeToCheck:  routes.NewRoute("/users", http.MethodGet),
			expectedError: nil,
			expectedOk:    true,
		},
		{
			name: "UndefinedRoute_DifferentPattern",
			currentRoutes: []routes.Route{
				routes.NewRoute("/books", http.MethodGet),
			},
			routeToCheck:  routes.NewRoute("/users", http.MethodGet),
			expectedError: nil,
			expectedOk:    true,
		},
		{
			name: "UndefinedRoute_DifferentMethod",
			currentRoutes: []routes.Route{
				routes.NewRoute("/users", http.MethodGet),
			},
			routeToCheck:  routes.NewRoute("/users", http.MethodPost),
			expectedError: nil,
			expectedOk:    true,
		},
		{
			name: "DefinedRoute_DifferentParameters",
			currentRoutes: []routes.Route{
				routes.NewRoute("/users/:id", http.MethodGet),
			},
			routeToCheck: routes.NewRoute("/users/:user_id", http.MethodGet),
			expectedError: routes.ErrRouteAlreadyDefined{
				Route: routes.NewRoute("/users/:user_id", http.MethodGet),
			},
			expectedOk: false,
		},
		{
			name:          "UndefinedRoute_DuplicateParameters",
			currentRoutes: []routes.Route{},
			routeToCheck:  routes.NewRoute("/users/:id/books/:id", http.MethodGet),
			expectedError: routes.ErrDuplicateArguments{Duplicate: ":id"},
			expectedOk:    false,
		},
		{
			name: "DefinedRoute_DifferentParameters",
			currentRoutes: []routes.Route{
				routes.NewRoute("/users/:id", http.MethodGet),
			},
			routeToCheck: routes.NewRoute("/users/:user_id", http.MethodGet),
			expectedError: routes.ErrRouteAlreadyDefined{
				Route: routes.NewRoute("/users/:user_id", http.MethodGet),
			},
			expectedOk: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			graph := routes.NewGraph()
			for _, route := range test.currentRoutes {
				graph.Add(route)
			}

			// Act
			actualError, actualOk := graph.CanBeInserted(test.routeToCheck)

			// Assert
			require.ErrorIs(t, test.expectedError, actualError)
			require.Equal(t, test.expectedOk, actualOk)
		})
	}
}

func TestGraph_MatchRoute(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name string

		currentRoutes []routes.Route
		routeToMatch  routes.Route

		expectedMatchedRoute routes.Route
		expectedOk           bool
	}

	tests := []testCase{
		{
			name:                 "EmptyGraph",
			currentRoutes:        []routes.Route{},
			routeToMatch:         routes.NewRoute("/users", http.MethodGet),
			expectedMatchedRoute: routes.Route{},
			expectedOk:           false,
		},
		{
			name: "UnmatchedRoute_DifferentPattern",
			currentRoutes: []routes.Route{
				routes.NewRoute("/books", http.MethodGet),
			},
			routeToMatch:         routes.NewRoute("/users", http.MethodGet),
			expectedMatchedRoute: routes.Route{},
			expectedOk:           false,
		},
		{
			name: "UnmatchedRoute_DifferentMethod",
			currentRoutes: []routes.Route{
				routes.NewRoute("/users", http.MethodGet),
			},
			routeToMatch:         routes.NewRoute("/users", http.MethodPost),
			expectedMatchedRoute: routes.Route{},
			expectedOk:           false,
		},
		{
			name: "UnmatchedRoute_Subpattern",
			currentRoutes: []routes.Route{
				routes.NewRoute("/users/:user_id", http.MethodGet),
			},
			routeToMatch:         routes.NewRoute("/users", http.MethodPost),
			expectedMatchedRoute: routes.Route{},
			expectedOk:           false,
		},
		{
			name: "UnmatchedRoute_Superpattern",
			currentRoutes: []routes.Route{
				routes.NewRoute("/users", http.MethodGet),
			},
			routeToMatch:         routes.NewRoute("/users/:user_id", http.MethodPost),
			expectedMatchedRoute: routes.Route{},
			expectedOk:           false,
		},
		{
			name: "MatchedRoute",
			currentRoutes: []routes.Route{
				routes.NewRoute("/users", http.MethodGet),
			},
			routeToMatch:         routes.NewRoute("/users", http.MethodGet),
			expectedMatchedRoute: routes.NewRoute("/users", http.MethodGet),
			expectedOk:           true,
		},
		{
			name: "MatchedRoute_Arguments",
			currentRoutes: []routes.Route{
				routes.NewRoute("/users/:user_id", http.MethodGet),
			},
			routeToMatch:         routes.NewRoute("/users/123", http.MethodGet),
			expectedMatchedRoute: routes.NewRoute("/users/:user_id", http.MethodGet),
			expectedOk:           true,
		},
		{
			name: "MatchedRoute_ArgumentInMiddle",
			currentRoutes: []routes.Route{
				routes.NewRoute("/users/:user_id/books", http.MethodGet),
			},
			routeToMatch:         routes.NewRoute("/users/123/books", http.MethodGet),
			expectedMatchedRoute: routes.NewRoute("/users/:user_id/books", http.MethodGet),
			expectedOk:           true,
		},
		{
			name: "UnmatchedRoute_ArgumentInMiddle",
			currentRoutes: []routes.Route{
				routes.NewRoute("/users/:user_id/books", http.MethodGet),
			},
			routeToMatch:         routes.NewRoute("/users/123", http.MethodGet),
			expectedMatchedRoute: routes.Route{},
			expectedOk:           false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			graph := routes.NewGraph()
			for _, route := range test.currentRoutes {
				graph.Add(route)
			}

			// Act
			match, actualOk := graph.MatchRoute(test.routeToMatch)

			// Assert
			require.Equal(t, test.expectedMatchedRoute, match.MatchedRoute())
			require.Equal(t, test.expectedOk, actualOk)
		})
	}
}
