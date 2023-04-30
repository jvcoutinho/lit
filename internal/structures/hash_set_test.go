package structures_test

import (
	"testing"

	"github.com/jvcoutinho/lit/internal/structures"

	"github.com/stretchr/testify/require"
)

func TestHashSet_Add(t *testing.T) {
	t.Parallel()

	type testCase[T comparable] struct {
		name string

		initialElements []int
		elementToAdd    T

		expectedSize int
	}

	tests := []testCase[int]{
		{
			name:            "ExistingElement",
			initialElements: []int{2},
			elementToAdd:    2,
			expectedSize:    1,
		},
		{
			name:            "NewElement",
			initialElements: []int{2, 3},
			elementToAdd:    4,
			expectedSize:    3,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			set := structures.NewHashSet[int](test.initialElements...)

			// Act
			set.Add(test.elementToAdd)

			// Assert
			require.Equal(t, test.expectedSize, set.Len())
			require.True(t, set.Contains(test.elementToAdd))
		})
	}
}
