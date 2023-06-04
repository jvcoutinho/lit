package trie_test

import (
	"net/http"
	"testing"

	"github.com/jvcoutinho/lit/internal/trie"
	"github.com/stretchr/testify/require"
)

func TestTrie_Insert(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description    string
		trie           *trie.Trie
		pattern        string
		method         string
		expectedResult *trie.Node
		expectedError  error
	}{
		{
			description:    "WhenMethodIsEmpty_ShouldReturnError",
			trie:           trie.New(),
			pattern:        "/users",
			method:         "",
			expectedResult: nil,
			expectedError:  trie.ErrMethodIsEmpty,
		},
		{
			description:    "WhenPatternDoesNotStartWithSlash_ShouldReturnError",
			trie:           trie.New(),
			pattern:        "users",
			method:         http.MethodGet,
			expectedResult: nil,
			expectedError:  trie.ErrPatternDoesNotStartWithSlash,
		},
		{
			description:    "WhenPatternContainsDoubleSlashes_ShouldReturnError",
			trie:           trie.New(),
			pattern:        "/users//",
			method:         http.MethodGet,
			expectedResult: nil,
			expectedError:  trie.ErrPatternContainsDoubleSlash,
		},
		{
			description: "GivenPatternAndMethodExists_ShouldReturnError",
			trie: func() *trie.Trie {
				trie := trie.New()
				trie.Insert("/users", http.MethodGet)
				return trie
			}(),
			pattern:        "/users",
			method:         http.MethodGet,
			expectedResult: nil,
			expectedError:  trie.ErrPatternHasBeenDefinedAlready,
		},
		{
			description: "GivenPatternAndMethodExistsWithDifferentParameters_ShouldReturnError",
			trie: func() *trie.Trie {
				trie := trie.New()
				trie.Insert("/users/:id", http.MethodGet)
				return trie
			}(),
			pattern:        "/users/:user_id",
			method:         http.MethodGet,
			expectedResult: nil,
			expectedError:  trie.ErrPatternHasConflictingParameters,
		},
		{
			description: "GivenMethodDoesNotExist_ShouldReturnNode",
			trie: func() *trie.Trie {
				trie := trie.New()
				trie.Insert("/users", http.MethodGet)
				return trie
			}(),
			pattern:        "/users",
			method:         http.MethodPost,
			expectedResult: trie.NewNode(),
			expectedError:  nil,
		},
		{
			description: "GivenPatternDoesNotExist_ShouldReturnNode",
			trie: func() *trie.Trie {
				trie := trie.New()
				trie.Insert("/books", http.MethodGet)
				return trie
			}(),
			pattern:        "/users",
			method:         http.MethodGet,
			expectedResult: trie.NewNode(),
			expectedError:  nil,
		},
		{
			description: "GivenPatternDoesNotExistWithTrailingSlash_ShouldReturnNode",
			trie: func() *trie.Trie {
				trie := trie.New()
				trie.Insert("/users", http.MethodGet)
				return trie
			}(),
			pattern:        "/users/",
			method:         http.MethodGet,
			expectedResult: trie.NewNode(),
			expectedError:  nil,
		},
		{
			description: "GivenPatternDoesNotExistWithParameter_ShouldReturnNode",
			trie: func() *trie.Trie {
				trie := trie.New()
				trie.Insert("/users/books", http.MethodGet)
				return trie
			}(),
			pattern:        "/users/:id/books",
			method:         http.MethodGet,
			expectedResult: trie.NewNode(),
			expectedError:  nil,
		},
		{
			description: "GivenSubpatternExistsButNotPattern_ShouldReturnNode",
			trie: func() *trie.Trie {
				trie := trie.New()
				trie.Insert("/users/:id", http.MethodGet)
				return trie
			}(),
			pattern:        "/users/:id/books",
			method:         http.MethodGet,
			expectedResult: trie.NewNode(),
			expectedError:  nil,
		},
		{
			description: "GivenSuperpatternExistsButNotPattern_ShouldReturnNode",
			trie: func() *trie.Trie {
				trie := trie.New()
				trie.Insert("/users/:id/books", http.MethodGet)
				return trie
			}(),
			pattern:        "/users/:id",
			method:         http.MethodGet,
			expectedResult: trie.NewNode(),
			expectedError:  nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Arrange
			trie := test.trie

			// Act
			actualResult, actualError := trie.Insert(test.pattern, test.method)

			// Assert
			require.Equal(t, test.expectedResult, actualResult)
			require.ErrorIs(t, actualError, test.expectedError)
		})
	}
}

func TestTrie_Match(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description       string
		trie              *trie.Trie
		pattern           string
		method            string
		expectedNode      *trie.Node
		expectedArguments map[string]string
		expectedError     error
	}{
		{
			description:       "GivenTrieHasNoRoutes_ShouldReturnError",
			trie:              trie.New(),
			pattern:           "/users",
			method:            http.MethodGet,
			expectedNode:      nil,
			expectedArguments: nil,
			expectedError:     trie.ErrMatchNotFound,
		},
		{
			description: "GivenPatternDoesNotExistInTrie_ShouldReturnError",
			trie: func() *trie.Trie {
				trie := trie.New()
				trie.Insert("/users/books", http.MethodGet)
				return trie
			}(),
			pattern:           "/users",
			method:            http.MethodGet,
			expectedNode:      nil,
			expectedArguments: nil,
			expectedError:     trie.ErrMatchNotFound,
		},
		{
			description: "GivenPatternWithoutTrailingSlashDoesNotExistInTrie_ShouldReturnError",
			trie: func() *trie.Trie {
				trie := trie.New()
				trie.Insert("/users/", http.MethodGet)
				return trie
			}(),
			pattern:           "/users",
			method:            http.MethodGet,
			expectedNode:      nil,
			expectedArguments: nil,
			expectedError:     trie.ErrMatchNotFound,
		},
		{
			description: "GivenPatternWithParametersDoesNotExistInTrie_ShouldReturnError",
			trie: func() *trie.Trie {
				trie := trie.New()
				trie.Insert("/users/:id", http.MethodGet)
				return trie
			}(),
			pattern:           "/users/user/123",
			method:            http.MethodGet,
			expectedNode:      nil,
			expectedArguments: nil,
			expectedError:     trie.ErrMatchNotFound,
		},
		{
			description: "GivenPatternIsSubpatternInTrie_ShouldReturnError",
			trie: func() *trie.Trie {
				trie := trie.New()
				trie.Insert("/users/:id", http.MethodGet)
				return trie
			}(),
			pattern:           "/users",
			method:            http.MethodGet,
			expectedNode:      nil,
			expectedArguments: nil,
			expectedError:     trie.ErrMatchNotFound,
		},
		{
			description: "GivenPatternIsSuperpatternInTrie_ShouldReturnError",
			trie: func() *trie.Trie {
				trie := trie.New()
				trie.Insert("/users/:id", http.MethodGet)
				return trie
			}(),
			pattern:           "/users/123/books",
			method:            http.MethodGet,
			expectedNode:      nil,
			expectedArguments: nil,
			expectedError:     trie.ErrMatchNotFound,
		},
		{
			description: "GivenMethodDoesNotExistInTrie_ShouldReturnError",
			trie: func() *trie.Trie {
				trie := trie.New()
				trie.Insert("/users/books", http.MethodGet)
				return trie
			}(),
			pattern:           "/user/books",
			method:            http.MethodPost,
			expectedNode:      nil,
			expectedArguments: nil,
			expectedError:     trie.ErrMatchNotFound,
		},
		{
			description: "GivenPatternAndMethodExistsInTrie_ShouldReturnNode",
			trie: func() *trie.Trie {
				trie := trie.New()
				trie.Insert("/users/books", http.MethodGet)
				return trie
			}(),
			pattern:           "/users/books",
			method:            http.MethodGet,
			expectedNode:      trie.NewNode(),
			expectedArguments: map[string]string{},
			expectedError:     nil,
		},
		{
			description: "GivenPatternWithParametersAndMethodExistsInTrie_ShouldReturnNode",
			trie: func() *trie.Trie {
				trie := trie.New()
				trie.Insert("/users/:id", http.MethodGet)
				return trie
			}(),
			pattern:           "/users/123",
			method:            http.MethodGet,
			expectedNode:      trie.NewNode(),
			expectedArguments: map[string]string{":id": "123"},
			expectedError:     nil,
		},
		{
			description: "GivenPatternWithDifferentMethodsAndMethodExistsInTrie_ShouldReturnNode",
			trie: func() *trie.Trie {
				trie := trie.New()
				trie.Insert("/users/:id", http.MethodGet)
				trie.Insert("/users/:user_id", http.MethodPatch)
				return trie
			}(),
			pattern:           "/users/123",
			method:            http.MethodPatch,
			expectedNode:      trie.NewNode(),
			expectedArguments: map[string]string{":user_id": "123"},
			expectedError:     nil,
		},
		{
			description: "GivenPatternAndSuperpatternWithParametersAndMethodExistsInTrie_ShouldReturnNode",
			trie: func() *trie.Trie {
				trie := trie.New()
				trie.Insert("/users/:user_id", http.MethodGet)
				trie.Insert("/users/:id/books", http.MethodGet)
				return trie
			}(),
			pattern:           "/users/123",
			method:            http.MethodGet,
			expectedNode:      trie.NewNode(),
			expectedArguments: map[string]string{":user_id": "123"},
			expectedError:     nil,
		},
		{
			description: "GivenPatternAndSubpatternWithParametersAndMethodExistsInTrie_ShouldReturnNode",
			trie: func() *trie.Trie {
				trie := trie.New()
				trie.Insert("/users/:id/books", http.MethodGet)
				trie.Insert("/users/", http.MethodGet)
				return trie
			}(),
			pattern:           "/users/123/books",
			method:            http.MethodGet,
			expectedNode:      trie.NewNode(),
			expectedArguments: map[string]string{":id": "123"},
			expectedError:     nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Arrange
			trie := test.trie

			// Act
			actualNode, actualArguments, actualError := trie.Match(test.pattern, test.method)

			// Assert
			require.Equal(t, test.expectedNode, actualNode)
			require.Equal(t, test.expectedArguments, actualArguments)
			require.Equal(t, test.expectedError, actualError)
		})
	}
}
