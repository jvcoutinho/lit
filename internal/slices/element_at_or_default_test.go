package slices_test

import (
	"testing"

	"github.com/jvcoutinho/lit/internal/slices"
	"github.com/stretchr/testify/require"
)

func TestElementAtOrDefault(t *testing.T) {
	t.Parallel()

	type testCase[T comparable] struct {
		name string

		currentSlice []T
		index        int

		expectedResult T
	}

	tests := []testCase[int]{
		{
			name:           "IndexInRange",
			currentSlice:   []int{2, 3},
			index:          0,
			expectedResult: 2,
		},
		{
			name:           "IndexLowerThanZero",
			currentSlice:   []int{3, 5},
			index:          -2,
			expectedResult: 0,
		},
		{
			name:           "IndexGreaterThanLengthOfTheSlice",
			currentSlice:   []int{3, 5},
			index:          5,
			expectedResult: 0,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			actualResult := slices.ElementAtOrDefault(test.currentSlice, test.index)

			require.Equal(t, test.expectedResult, actualResult)
		})
	}
}
