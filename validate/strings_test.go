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
			expected:    false,
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
			expected:    false,
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

func TestSubstring(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		target      *string
		substring   string
		expected    bool
	}{
		{
			description: "NilPointer",
			target:      nil,
			substring:   "substring",
			expected:    false,
		},
		{
			description: "HasSubstring",
			target:      pointerOf("123"),
			substring:   "2",
			expected:    true,
		},
		{
			description: "DoesNotHaveSubstring",
			target:      pointerOf("123"),
			substring:   "4",
			expected:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Act
			result := validate.Substring(test.target, test.substring)

			// Assert
			expectedResult := validate.Field{
				Valid:   test.expected,
				Message: fmt.Sprintf(`{0} should contain "%s"`, test.substring),
				Fields:  []any{test.target},
			}

			require.Equal(t, expectedResult, result)
		})
	}
}

func TestLowercase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		target      *string
		expected    bool
	}{
		{
			description: "NilPointer",
			target:      nil,
			expected:    false,
		},
		{
			description: "IsLowercase",
			target:      pointerOf("string"),
			expected:    true,
		},
		{
			description: "IsNotLowercase",
			target:      pointerOf("String"),
			expected:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Act
			result := validate.Lowercase(test.target)

			// Assert
			expectedResult := validate.Field{
				Valid:   test.expected,
				Message: "{0} should contain only lowercase characters",
				Fields:  []any{test.target},
			}

			require.Equal(t, expectedResult, result)
		})
	}
}

func TestUppercase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		target      *string
		expected    bool
	}{
		{
			description: "NilPointer",
			target:      nil,
			expected:    false,
		},
		{
			description: "IsUppercase",
			target:      pointerOf("STRING"),
			expected:    true,
		},
		{
			description: "IsNotUppercase",
			target:      pointerOf("String"),
			expected:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Act
			result := validate.Uppercase(test.target)

			// Assert
			expectedResult := validate.Field{
				Valid:   test.expected,
				Message: "{0} should contain only uppercase characters",
				Fields:  []any{test.target},
			}

			require.Equal(t, expectedResult, result)
		})
	}
}
