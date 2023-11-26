package validate_test

import (
	"fmt"
	"testing"

	"github.com/jvcoutinho/lit/validate"
	"github.com/stretchr/testify/require"
)

func TestHasPrefix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		target      *string
		prefix      string
		expected    bool
	}{
		{
			description: "NilPointer",
			target:      nil,
			prefix:      "prefix",
			expected:    true,
		},
		{
			description: "HasPrefix",
			target:      pointerOf("prefix-123"),
			prefix:      "prefix",
			expected:    true,
		},
		{
			description: "DoesNotHavePrefix",
			target:      pointerOf("123"),
			prefix:      "prefix",
			expected:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Act
			result := validate.HasPrefix(test.target, test.prefix)

			// Assert
			expectedResult := validate.Field{
				Valid:   test.expected,
				Message: fmt.Sprintf(`{0} should start with "%s"`, test.prefix),
				Fields:  []any{test.target},
			}

			require.Equal(t, expectedResult, result)
		})
	}
}

func TestHasSuffix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		target      *string
		suffix      string
		expected    bool
	}{
		{
			description: "NilPointer",
			target:      nil,
			suffix:      "suffix",
			expected:    true,
		},
		{
			description: "HasSuffix",
			target:      pointerOf("123-suffix"),
			suffix:      "suffix",
			expected:    true,
		},
		{
			description: "DoesNotHaveSuffix",
			target:      pointerOf("123"),
			suffix:      "suffix",
			expected:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Act
			result := validate.HasSuffix(test.target, test.suffix)

			// Assert
			expectedResult := validate.Field{
				Valid:   test.expected,
				Message: fmt.Sprintf(`{0} should end with "%s"`, test.suffix),
				Fields:  []any{test.target},
			}

			require.Equal(t, expectedResult, result)
		})
	}
}
