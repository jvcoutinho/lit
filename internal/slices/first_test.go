package slices_test

import (
	"testing"

	"github.com/jvcoutinho/lit/internal/slices"
	"github.com/stretchr/testify/require"
)

func TestFirst(t *testing.T) {
	t.Parallel()

	type testCase[T comparable] struct {
		name string

		currentSlice []T
		predicate    func(T) bool

		expectedResult T
		expectedOk     bool
	}

	tests := []testCase[int]{
		{
			name:           "SomeElementMatchesPredicate",
			currentSlice:   []int{2, 3},
			predicate:      func(i int) bool { return i%2 == 0 },
			expectedResult: 2,
			expectedOk:     true,
		},
		{
			name:           "NoElementMatchesPredicate",
			currentSlice:   []int{3, 5},
			predicate:      func(i int) bool { return i%2 == 0 },
			expectedResult: 0,
			expectedOk:     false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			actualResult, actualOk := slices.First(test.currentSlice, test.predicate)

			require.Equal(t, test.expectedResult, actualResult)
			require.Equal(t, test.expectedOk, actualOk)
		})
	}
}
