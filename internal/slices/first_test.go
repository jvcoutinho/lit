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
		expectedError  error
	}

	tests := []testCase[int]{
		{
			name:           "SomeElementMatchesPredicate",
			currentSlice:   []int{2, 3},
			predicate:      func(i int) bool { return i%2 == 0 },
			expectedResult: 2,
			expectedError:  nil,
		},
		{
			name:           "NoElementMatchesPredicate",
			currentSlice:   []int{3, 5},
			predicate:      func(i int) bool { return i%2 == 0 },
			expectedResult: 0,
			expectedError:  slices.ErrNoElementFound,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			actualResult, actualError := slices.First(test.currentSlice, test.predicate)

			require.Equal(t, test.expectedResult, actualResult)
			require.Equal(t, test.expectedError, actualError)
		})
	}
}
