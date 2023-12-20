package validate_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jvcoutinho/lit/validate"
)

func TestGreater(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		target      *int
		value       int
		expected    bool
	}{
		{
			description: "NilPointer",
			target:      nil,
			value:       10,
			expected:    false,
		},
		{
			description: "Greater",
			target:      pointerOf(20),
			value:       10,
			expected:    true,
		},
		{
			description: "NotGreater",
			target:      pointerOf(5),
			value:       10,
			expected:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Act
			result := validate.Greater(test.target, test.value)

			// Assert
			expectedResult := validate.Field{
				Valid:   test.expected,
				Message: fmt.Sprintf("{0} should be greater than %v", test.value),
				Fields:  []any{test.target},
			}

			require.Equal(t, expectedResult, result)
		})
	}
}

func TestGreaterOrEqual(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		target      *int
		value       int
		expected    bool
	}{
		{
			description: "NilPointer",
			target:      nil,
			value:       10,
			expected:    false,
		},
		{
			description: "Greater",
			target:      pointerOf(20),
			value:       10,
			expected:    true,
		},
		{
			description: "Equal",
			target:      pointerOf(10),
			value:       10,
			expected:    true,
		},
		{
			description: "NotGreaterOrEqual",
			target:      pointerOf(5),
			value:       10,
			expected:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Act
			result := validate.GreaterOrEqual(test.target, test.value)

			// Assert
			expectedResult := validate.Field{
				Valid:   test.expected,
				Message: fmt.Sprintf("{0} should be greater or equal than %v", test.value),
				Fields:  []any{test.target},
			}

			require.Equal(t, expectedResult, result)
		})
	}
}

func TestLess(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		target      *int
		value       int
		expected    bool
	}{
		{
			description: "NilPointer",
			target:      nil,
			value:       10,
			expected:    false,
		},
		{
			description: "Less",
			target:      pointerOf(5),
			value:       10,
			expected:    true,
		},
		{
			description: "NotLess",
			target:      pointerOf(20),
			value:       10,
			expected:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Act
			result := validate.Less(test.target, test.value)

			// Assert
			expectedResult := validate.Field{
				Valid:   test.expected,
				Message: fmt.Sprintf("{0} should be less than %v", test.value),
				Fields:  []any{test.target},
			}

			require.Equal(t, expectedResult, result)
		})
	}
}

func TestLessOrEqual(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		target      *int
		value       int
		expected    bool
	}{
		{
			description: "NilPointer",
			target:      nil,
			value:       10,
			expected:    false,
		},
		{
			description: "Less",
			target:      pointerOf(5),
			value:       10,
			expected:    true,
		},
		{
			description: "Equal",
			target:      pointerOf(10),
			value:       10,
			expected:    true,
		},
		{
			description: "NotLess",
			target:      pointerOf(20),
			value:       10,
			expected:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Act
			result := validate.LessOrEqual(test.target, test.value)

			// Assert
			expectedResult := validate.Field{
				Valid:   test.expected,
				Message: fmt.Sprintf("{0} should be less or equal than %v", test.value),
				Fields:  []any{test.target},
			}

			require.Equal(t, expectedResult, result)
		})
	}
}

func TestBetween(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		target      *int
		min         int
		max         int
		expected    bool
	}{
		{
			description: "NilPointer",
			target:      nil,
			min:         10,
			max:         20,
			expected:    false,
		},
		{
			description: "LessThanMin",
			target:      pointerOf(5),
			min:         10,
			max:         20,
			expected:    false,
		},
		{
			description: "GreaterThanMax",
			target:      pointerOf(30),
			min:         10,
			max:         20,
			expected:    false,
		},
		{
			description: "Between",
			target:      pointerOf(15),
			min:         10,
			max:         20,
			expected:    true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Act
			result := validate.Between(test.target, test.min, test.max)

			// Assert
			expectedResult := validate.Field{
				Valid:   test.expected,
				Message: fmt.Sprintf("{0} should be between %d and %d", test.min, test.max),
				Fields:  []any{test.target},
			}

			require.Equal(t, expectedResult, result)
		})
	}
}
