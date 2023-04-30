package structures_test

import (
	"testing"

	"github.com/jvcoutinho/lit/internal/structures"
	"github.com/stretchr/testify/require"
)

func TestList_InsertAtBeginning(t *testing.T) {
	t.Parallel()

	type testCase[T comparable] struct {
		name string

		initialElements []int
		elementToAdd    T

		expectedSize int
	}

	tests := []testCase[int]{
		{
			name:            "EmptyList",
			initialElements: []int{},
			elementToAdd:    2,
			expectedSize:    1,
		},
		{
			name:            "NotEmptyList",
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
			list := structures.NewList[int](test.initialElements...)

			// Act
			list.InsertAtBeginning(test.elementToAdd)

			// Assert
			elementAtStart, _ := list.ElementAt(0)

			require.Equal(t, test.expectedSize, list.Len())
			require.Equal(t, test.elementToAdd, elementAtStart)
		})
	}
}
