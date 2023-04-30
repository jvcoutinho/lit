package maps_test

import (
	"testing"

	"github.com/jvcoutinho/lit/internal/maps"
	"github.com/stretchr/testify/require"
)

func TestContainsKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string

		currentValues map[string]int
		keyToCheck    string

		expectedResult bool
	}{
		{
			name:           "MapContainsKey",
			currentValues:  map[string]int{"2": 2, "3": 3},
			keyToCheck:     "2",
			expectedResult: true,
		},
		{
			name:           "MapDoesNotContainKey",
			currentValues:  map[string]int{"2": 2, "3": 3},
			keyToCheck:     "5",
			expectedResult: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, test.expectedResult, maps.ContainsKey(test.currentValues, test.keyToCheck))
		})
	}
}
