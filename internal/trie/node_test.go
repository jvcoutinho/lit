package trie_test

import (
	"testing"

	"github.com/jvcoutinho/lit/internal/trie"
	"github.com/stretchr/testify/require"
)

func TestNode_IsTerminal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description    string
		staticChildren map[string]*trie.Node
		dynamicChild   *trie.Node
		expectedResult bool
	}{
		{
			description:    "GivenNodeHasNeitherStaticChildrenNorDynamicChild_ShouldReturnTrue",
			staticChildren: nil,
			dynamicChild:   nil,
			expectedResult: true,
		},
		{
			description:    "GivenNodeHasStaticChildren_ShouldReturnFalse",
			staticChildren: map[string]*trie.Node{"child": trie.NewNode()},
			dynamicChild:   nil,
			expectedResult: false,
		},
		{
			description:    "GivenNodeHasDynamicChild_ShouldReturnFalse",
			staticChildren: nil,
			dynamicChild:   trie.NewNode(),
			expectedResult: false,
		},
		{
			description:    "GivenNodeHasStaticChildrenAndDynamicChild_ShouldReturnFalse",
			staticChildren: map[string]*trie.Node{"child": trie.NewNode()},
			dynamicChild:   trie.NewNode(),
			expectedResult: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Arrange
			node := trie.NewNode()
			node.StaticChildren = test.staticChildren
			node.DynamicChild = test.dynamicChild

			// Act
			actualResult := node.IsTerminal()

			// Assert
			require.Equal(t, test.expectedResult, actualResult)
		})
	}
}
