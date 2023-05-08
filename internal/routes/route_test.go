package routes_test

import (
	"net/http"
	"testing"

	"github.com/jvcoutinho/lit/internal/routes"
	"github.com/stretchr/testify/require"
)

func TestRoute_Path(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		description    string
		method         string
		pattern        string
		expectedResult []string
	}

	tests := []TestCase{
		{
			description:    "GivenPatternIsEmpty_ShouldReturnSliceWithEmptyString",
			method:         http.MethodGet,
			pattern:        "",
			expectedResult: []string{""},
		},
		{
			description:    "GivenPatternIsOnlySlash_ShouldReturnSliceWithEmptyString",
			method:         http.MethodGet,
			pattern:        "/",
			expectedResult: []string{""},
		},
		{
			description:    "GivenPatternIsSingle_ShouldReturnPath",
			method:         http.MethodGet,
			pattern:        "users",
			expectedResult: []string{"users"},
		},
		{
			description:    "GivenPatternIsSingle_AndPatternBeginsWithSingleSlash_ShouldReturnPath",
			method:         http.MethodGet,
			pattern:        "/users",
			expectedResult: []string{"users"},
		},
		{
			description:    "GivenPatternIsSingle_AndPatternBeginsWithMultipleSlashes_ShouldReturnPath",
			method:         http.MethodGet,
			pattern:        "////users",
			expectedResult: []string{"users"},
		},
		{
			description:    "GivenPatternIsSingle_AndPatternEndsWithSingleSlash_ShouldReturnPath",
			method:         http.MethodGet,
			pattern:        "users/",
			expectedResult: []string{"users"},
		},
		{
			description:    "GivenPatternIsSingle_AndPatternEndsWithMultipleSlashes_ShouldReturnPath",
			method:         http.MethodGet,
			pattern:        "users////",
			expectedResult: []string{"users"},
		},
		{
			description:    "GivenPatternIsMultiplePaths_ShouldReturnPath",
			method:         http.MethodGet,
			pattern:        "users/:user_id/books",
			expectedResult: []string{"users", ":user_id", "books"},
		},
		{
			description:    "GivenPatternIsMultiplePaths_AndPatternBeginsWithSingleSlash_ShouldReturnPath",
			method:         http.MethodGet,
			pattern:        "/users/:user_id/books",
			expectedResult: []string{"users", ":user_id", "books"},
		},
		{
			description:    "GivenPatternIsMultiplePaths_AndPatternBeginsWithMultipleSlashes_ShouldReturnPath",
			method:         http.MethodGet,
			pattern:        "////users/:user_id/books",
			expectedResult: []string{"users", ":user_id", "books"},
		},
		{
			description:    "GivenPatternIsMultiplePaths_AndPatternEndsWithSingleSlash_ShouldReturnPath",
			method:         http.MethodGet,
			pattern:        "users/:user_id/books/",
			expectedResult: []string{"users", ":user_id", "books"},
		},
		{
			description:    "GivenPatternIsMultiplePaths_AndPatternEndsWithMultipleSlashes_ShouldReturnPath",
			method:         http.MethodGet,
			pattern:        "users/:user_id/books////",
			expectedResult: []string{"users", ":user_id", "books"},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Arrange
			route := routes.NewRoute(test.pattern, test.method)

			// Act
			actualResult := route.Path()

			// Assert
			require.Equal(t, test.expectedResult, actualResult)
		})
	}
}

func TestRoute_String(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		description    string
		method         string
		pattern        string
		expectedResult string
	}

	tests := []TestCase{
		{
			description:    "GivenPatternIsEmpty_ShouldReturnSingleSlash",
			method:         http.MethodGet,
			pattern:        "",
			expectedResult: "GET /",
		},
		{
			description:    "ShouldReturnMethodWithPattern",
			method:         http.MethodGet,
			pattern:        "users/:user_id/books",
			expectedResult: "GET /users/:user_id/books",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Arrange
			route := routes.NewRoute(test.pattern, test.method)

			// Act
			actualResult := route.String()

			// Assert
			require.Equal(t, test.expectedResult, actualResult)
		})
	}
}
