package slices_test

import (
	"testing"

	"github.com/jvcoutinho/lit/internal/slices"
	"github.com/stretchr/testify/require"
)

func TestContains(t *testing.T) {
	t.Parallel()

	type testCase[T comparable] struct {
		name string

		currentSlice   []T
		elementToCheck T

		expectedResult bool
	}

	tests := []testCase[int]{
		{
			name:           "ElementIsContained",
			currentSlice:   []int{2, 3},
			elementToCheck: 2,
			expectedResult: true,
		},
		{
			name:           "ElementIsNotContained",
			currentSlice:   []int{2, 3},
			elementToCheck: 4,
			expectedResult: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			actualResult := slices.Contains(test.currentSlice, test.elementToCheck)

			require.Equal(t, test.expectedResult, actualResult)
		})
	}
}
