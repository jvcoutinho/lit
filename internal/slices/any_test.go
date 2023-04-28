package slices_test

import (
	"testing"

	"github.com/jvcoutinho/lit/internal/slices"
	"github.com/stretchr/testify/require"
)

func TestAny(t *testing.T) {
	t.Parallel()

	type testCase[T comparable] struct {
		name string

		currentSlice []T
		predicate    func(T) bool

		expectedResult bool
	}

	tests := []testCase[int]{
		{
			name:           "SomeElementMatchesPredicate",
			currentSlice:   []int{2, 3},
			predicate:      func(i int) bool { return i%2 == 0 },
			expectedResult: true,
		},
		{
			name:           "NoElementMatchesPredicate",
			currentSlice:   []int{3, 5},
			predicate:      func(i int) bool { return i%2 == 0 },
			expectedResult: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			actualResult := slices.Any(test.currentSlice, test.predicate)

			require.Equal(t, test.expectedResult, actualResult)
		})
	}
}
