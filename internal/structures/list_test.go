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

		initialElements []T
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

func TestList_ElementAt(t *testing.T) {
	t.Parallel()

	type testCase[T comparable] struct {
		name string

		elements []T
		index    int

		expectedResult T
		expectedOk     bool
	}

	tests := []testCase[int]{
		{
			name:           "IndexLessThanZero",
			elements:       []int{2, 3, 4},
			index:          -1,
			expectedResult: 0,
			expectedOk:     false,
		},
		{
			name:           "IndexGreaterThanListLength_EmptyList",
			elements:       []int{},
			index:          0,
			expectedResult: 0,
			expectedOk:     false,
		},
		{
			name:           "IndexGreaterThanListLength_NonEmptyList",
			elements:       []int{2, 3, 4},
			index:          3,
			expectedResult: 0,
			expectedOk:     false,
		},
		{
			name:           "IndexInRange",
			elements:       []int{2, 3, 4},
			index:          1,
			expectedResult: 3,
			expectedOk:     true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			list := structures.NewList[int](test.elements...)

			// Act
			actualResult, actualOk := list.ElementAt(test.index)

			// Assert
			require.Equal(t, test.expectedResult, actualResult)
			require.Equal(t, test.expectedOk, actualOk)
		})
	}
}

func TestList_Traverse(t *testing.T) {
	t.Parallel()

	type testCase[T comparable] struct {
		name string

		elements []T

		expectedSum T
	}

	tests := []testCase[int]{
		{
			name:        "EmptyList",
			elements:    []int{},
			expectedSum: 0,
		},
		{
			name:        "NotEmptyList",
			elements:    []int{2, 3, 4},
			expectedSum: 9,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			list := structures.NewList[int](test.elements...)
			sum := 0

			// Act
			list.Traverse(func(n int) {
				sum += n
			})

			// Assert
			require.Equal(t, test.expectedSum, sum)
		})
	}
}
