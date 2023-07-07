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
		setupTrie      func(*trie.Trie)
		pattern        string
		method         string
		expectedResult *trie.Node
		expectedError  error
	}{
		{
			description: "GivenPatternAndMethodExists_ShouldReturnError",
			setupTrie: func(t *trie.Trie) {
				_, _ = t.Insert("/users", http.MethodGet)
			},
			pattern:        "/users",
			method:         http.MethodGet,
			expectedResult: nil,
			expectedError:  trie.ErrPatternHasBeenDefinedAlready,
		},
		{
			description: "GivenPatternAndMethodExistsWithDifferentParameters_ShouldReturnError",
			setupTrie: func(t *trie.Trie) {
				_, _ = t.Insert("/users/:id", http.MethodGet)
			},
			pattern:        "/users/:user_id",
			method:         http.MethodGet,
			expectedResult: nil,
			expectedError:  trie.ErrPatternHasConflictingParameters,
		},
		{
			description: "GivenMethodDoesNotExist_ShouldReturnNode",
			setupTrie: func(t *trie.Trie) {
				_, _ = t.Insert("/users", http.MethodGet)
			},
			pattern:        "/users",
			method:         http.MethodPost,
			expectedResult: trie.NewNode(),
			expectedError:  nil,
		},
		{
			description: "GivenPatternDoesNotExist_ShouldReturnNode",
			setupTrie: func(t *trie.Trie) {
				_, _ = t.Insert("/books", http.MethodGet)
			},
			pattern:        "/users",
			method:         http.MethodGet,
			expectedResult: trie.NewNode(),
			expectedError:  nil,
		},
		{
			description: "GivenPatternDoesNotExistWithTrailingSlash_ShouldReturnNode",
			setupTrie: func(t *trie.Trie) {
				_, _ = t.Insert("/users", http.MethodGet)
			},
			pattern:        "/users/",
			method:         http.MethodGet,
			expectedResult: trie.NewNode(),
			expectedError:  nil,
		},
		{
			description: "GivenPatternDoesNotExistWithParameter_ShouldReturnNode",
			setupTrie: func(t *trie.Trie) {
				_, _ = t.Insert("/users/books", http.MethodGet)
			},
			pattern:        "/users/:id/books",
			method:         http.MethodGet,
			expectedResult: trie.NewNode(),
			expectedError:  nil,
		},
		{
			description: "GivenSubpatternExistsButNotPattern_ShouldReturnNode",
			setupTrie: func(t *trie.Trie) {
				_, _ = t.Insert("/users/:id", http.MethodGet)
			},
			pattern:        "/users/:id/books",
			method:         http.MethodGet,
			expectedResult: trie.NewNode(),
			expectedError:  nil,
		},
		{
			description: "GivenSuperpatternExistsButNotPattern_ShouldReturnNode",
			setupTrie: func(t *trie.Trie) {
				_, _ = t.Insert("/users/:id/books", http.MethodGet)
			},
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
			tr := trie.New()
			test.setupTrie(tr)

			// Act
			actualResult, actualError := tr.Insert(test.pattern, test.method)

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
		setupTrie         func(*trie.Trie)
		pattern           string
		method            string
		expectedNode      *trie.Node
		expectedArguments map[string]string
		expectedError     error
	}{
		{
			description:       "GivenTrieHasNoRoutes_ShouldReturnError",
			setupTrie:         func(t *trie.Trie) {},
			pattern:           "/users",
			method:            http.MethodGet,
			expectedNode:      nil,
			expectedArguments: nil,
			expectedError:     trie.ErrMatchNotFound,
		},
		{
			description: "GivenPatternDoesNotExistInTrie_ShouldReturnError",
			setupTrie: func(t *trie.Trie) {
				_, _ = t.Insert("/users/books", http.MethodGet)
			},
			pattern:           "/users",
			method:            http.MethodGet,
			expectedNode:      nil,
			expectedArguments: nil,
			expectedError:     trie.ErrMatchNotFound,
		},
		{
			description: "GivenPatternWithoutTrailingSlashDoesNotExistInTrie_ShouldReturnError",
			setupTrie: func(t *trie.Trie) {
				_, _ = t.Insert("/users/", http.MethodGet)
			},
			pattern:           "/users",
			method:            http.MethodGet,
			expectedNode:      nil,
			expectedArguments: nil,
			expectedError:     trie.ErrMatchNotFound,
		},
		{
			description: "GivenPatternWithParametersDoesNotExistInTrie_ShouldReturnError",
			setupTrie: func(t *trie.Trie) {
				_, _ = t.Insert("/users/:id", http.MethodGet)
			},
			pattern:           "/users/user/123",
			method:            http.MethodGet,
			expectedNode:      nil,
			expectedArguments: nil,
			expectedError:     trie.ErrMatchNotFound,
		},
		{
			description: "GivenPatternIsSubpatternInTrie_ShouldReturnError",
			setupTrie: func(t *trie.Trie) {
				_, _ = t.Insert("/users/:id", http.MethodGet)
			},
			pattern:           "/users",
			method:            http.MethodGet,
			expectedNode:      nil,
			expectedArguments: nil,
			expectedError:     trie.ErrMatchNotFound,
		},
		{
			description: "GivenPatternIsSuperpatternInTrie_ShouldReturnError",
			setupTrie: func(t *trie.Trie) {
				_, _ = t.Insert("/users/:id", http.MethodGet)
			},
			pattern:           "/users/123/books",
			method:            http.MethodGet,
			expectedNode:      nil,
			expectedArguments: nil,
			expectedError:     trie.ErrMatchNotFound,
		},
		{
			description: "GivenMethodDoesNotExistInTrie_ShouldReturnError",
			setupTrie: func(t *trie.Trie) {
				_, _ = t.Insert("/users/books", http.MethodGet)
			},
			pattern:           "/user/books",
			method:            http.MethodPost,
			expectedNode:      nil,
			expectedArguments: nil,
			expectedError:     trie.ErrMatchNotFound,
		},
		{
			description: "GivenPatternAndMethodExistsInTrie_ShouldReturnNode",
			setupTrie: func(t *trie.Trie) {
				_, _ = t.Insert("/users/books", http.MethodGet)
			},
			pattern:           "/users/books",
			method:            http.MethodGet,
			expectedNode:      trie.NewNode(),
			expectedArguments: map[string]string{},
			expectedError:     nil,
		},
		{
			description: "GivenPatternWithParametersAndMethodExistsInTrie_ShouldReturnNode",
			setupTrie: func(t *trie.Trie) {
				_, _ = t.Insert("/users/:id", http.MethodGet)
			},
			pattern:           "/users/123",
			method:            http.MethodGet,
			expectedNode:      trie.NewNode(),
			expectedArguments: map[string]string{":id": "123"},
			expectedError:     nil,
		},
		{
			description: "GivenPatternWithDifferentMethodsAndMethodExistsInTrie_ShouldReturnNode",
			setupTrie: func(t *trie.Trie) {
				_, _ = t.Insert("/users/:id", http.MethodGet)
				_, _ = t.Insert("/users/:user_id", http.MethodPatch)
			},
			pattern:           "/users/123",
			method:            http.MethodPatch,
			expectedNode:      trie.NewNode(),
			expectedArguments: map[string]string{":user_id": "123"},
			expectedError:     nil,
		},
		{
			description: "GivenPatternAndSuperpatternWithParametersAndMethodExistsInTrie_ShouldReturnNode",
			setupTrie: func(t *trie.Trie) {
				_, _ = t.Insert("/users/:user_id", http.MethodGet)
				_, _ = t.Insert("/users/:id/books", http.MethodGet)
			},
			pattern:           "/users/123",
			method:            http.MethodGet,
			expectedNode:      trie.NewNode(),
			expectedArguments: map[string]string{":user_id": "123"},
			expectedError:     nil,
		},
		{
			description: "GivenPatternAndSubpatternWithParametersAndMethodExistsInTrie_ShouldReturnNode",
			setupTrie: func(t *trie.Trie) {
				_, _ = t.Insert("/users/:id/books", http.MethodGet)
				_, _ = t.Insert("/users/", http.MethodGet)
			},
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
			tr := trie.New()
			test.setupTrie(tr)

			// Act
			actualNode, actualArguments, actualError := tr.Match(test.pattern, test.method)

			// Assert
			require.Equal(t, test.expectedNode, actualNode)
			require.Equal(t, test.expectedArguments, actualArguments)
			require.Equal(t, test.expectedError, actualError)
		})
	}
}
