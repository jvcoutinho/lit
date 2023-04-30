package maps_test

import (
	"testing"

	"github.com/jvcoutinho/lit/internal/maps"
	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Parallel()

	type testCase[T comparable, V any] struct {
		name string

		currentMap map[T]V
	}

	tests := []testCase[string, int]{
		{
			name:       "NotNilMap",
			currentMap: map[string]int{"2": 2, "3": 3},
		},
		{
			name:       "NilMap",
			currentMap: nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			mapCopy := maps.Copy(test.currentMap)

			require.Equal(t, mapCopy, test.currentMap)
			require.NotSame(t, mapCopy, test.currentMap)
		})
	}
}
