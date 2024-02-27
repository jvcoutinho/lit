package validate_test

import (
	"fmt"
	"testing"

	"github.com/jvcoutinho/lit/validate"
	"github.com/stretchr/testify/require"
)

func TestEqual(t *testing.T) {
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
			description: "Equal",
			target:      pointerOf(10),
			value:       10,
			expected:    true,
		},
		{
			description: "NotEqual",
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
			result := validate.Equal(test.target, test.value)

			// Assert
			expectedResult := validate.Field{
				Valid:   test.expected,
				Message: fmt.Sprintf("{0} should be equal to %v", test.value),
				Fields:  []any{test.target},
			}

			require.Equal(t, expectedResult, result)
		})
	}
}

func TestNotEqual(t *testing.T) {
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
			description: "Equal",
			target:      pointerOf(10),
			value:       10,
			expected:    false,
		},
		{
			description: "NotEqual",
			target:      pointerOf(5),
			value:       10,
			expected:    true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Act
			result := validate.NotEqual(test.target, test.value)

			// Assert
			expectedResult := validate.Field{
				Valid:   test.expected,
				Message: fmt.Sprintf("{0} should not be equal to %v", test.value),
				Fields:  []any{test.target},
			}

			require.Equal(t, expectedResult, result)
		})
	}
}

func TestOneOf(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		target      *int
		values      []int
		expected    bool
	}{
		{
			description: "NilPointer",
			target:      nil,
			values:      []int{10},
			expected:    false,
		},
		{
			description: "OneOf",
			target:      pointerOf(10),
			values:      []int{10, 15, 20},
			expected:    true,
		},
		{
			description: "NotOneOf",
			target:      pointerOf(5),
			values:      []int{10, 15, 20},
			expected:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Act
			result := validate.OneOf(test.target, test.values...)

			// Assert
			expectedResult := validate.Field{
				Valid:   test.expected,
				Message: fmt.Sprintf("{0} should be one of %v", test.values),
				Fields:  []any{test.target},
			}

			require.Equal(t, expectedResult, result)
		})
	}
}

func TestEqualField(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		target      *int
		field       *int
		expected    bool
	}{
		{
			description: "BothAreNil",
			target:      nil,
			field:       nil,
			expected:    true,
		},
		{
			description: "TargetIsNil",
			target:      nil,
			field:       pointerOf(10),
			expected:    false,
		},
		{
			description: "FieldIsNil",
			target:      pointerOf(10),
			field:       nil,
			expected:    false,
		},
		{
			description: "Equal",
			target:      pointerOf(10),
			field:       pointerOf(10),
			expected:    true,
		},
		{
			description: "NotEqual",
			target:      pointerOf(5),
			field:       pointerOf(10),
			expected:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Act
			result := validate.EqualField(test.target, test.field)

			// Assert
			expectedResult := validate.Field{
				Valid:   test.expected,
				Message: "{0} should be equal to {1}",
				Fields:  []any{test.target, test.field},
			}

			require.Equal(t, expectedResult, result)
		})
	}
}

func TestNotEqualField(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		target      *int
		field       *int
		expected    bool
	}{
		{
			description: "BothAreNil",
			target:      nil,
			field:       nil,
			expected:    false,
		},
		{
			description: "Equal",
			target:      pointerOf(10),
			field:       pointerOf(10),
			expected:    false,
		},
		{
			description: "NotEqual",
			target:      pointerOf(5),
			field:       pointerOf(10),
			expected:    true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Act
			result := validate.NotEqualField(test.target, test.field)

			// Assert
			expectedResult := validate.Field{
				Valid:   test.expected,
				Message: "{0} should not be equal to {1}",
				Fields:  []any{test.target, test.field},
			}

			require.Equal(t, expectedResult, result)
		})
	}
}

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

func TestGreaterField(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		target      *int
		field       *int
		expected    bool
	}{
		{
			description: "TargetIsNil",
			target:      nil,
			field:       pointerOf(10),
			expected:    false,
		},
		{
			description: "FieldIsNil",
			target:      pointerOf(10),
			field:       nil,
			expected:    false,
		},
		{
			description: "Greater",
			target:      pointerOf(20),
			field:       pointerOf(10),
			expected:    true,
		},
		{
			description: "NotGreater",
			target:      pointerOf(5),
			field:       pointerOf(10),
			expected:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Act
			result := validate.GreaterField(test.target, test.field)

			// Assert
			expectedResult := validate.Field{
				Valid:   test.expected,
				Message: "{0} should be greater than {1}",
				Fields:  []any{test.target, test.field},
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

func TestGreaterOrEqualField(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		target      *int
		field       *int
		expected    bool
	}{
		{
			description: "TargetIsNil",
			target:      nil,
			field:       pointerOf(10),
			expected:    false,
		},
		{
			description: "FieldIsNil",
			target:      pointerOf(10),
			field:       nil,
			expected:    false,
		},
		{
			description: "Greater",
			target:      pointerOf(20),
			field:       pointerOf(10),
			expected:    true,
		},
		{
			description: "Equal",
			target:      pointerOf(10),
			field:       pointerOf(10),
			expected:    true,
		},
		{
			description: "NotGreaterOrEqual",
			target:      pointerOf(5),
			field:       pointerOf(10),
			expected:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Act
			result := validate.GreaterOrEqualField(test.target, test.field)

			// Assert
			expectedResult := validate.Field{
				Valid:   test.expected,
				Message: "{0} should be greater or equal than {1}",
				Fields:  []any{test.target, test.field},
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

func TestLessField(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		target      *int
		field       *int
		expected    bool
	}{
		{
			description: "TargetIsNil",
			target:      nil,
			field:       pointerOf(10),
			expected:    false,
		},
		{
			description: "FieldIsNil",
			target:      pointerOf(10),
			field:       nil,
			expected:    false,
		},
		{
			description: "Less",
			target:      pointerOf(5),
			field:       pointerOf(10),
			expected:    true,
		},
		{
			description: "NotLess",
			target:      pointerOf(20),
			field:       pointerOf(10),
			expected:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Act
			result := validate.LessField(test.target, test.field)

			// Assert
			expectedResult := validate.Field{
				Valid:   test.expected,
				Message: "{0} should be less than {1}",
				Fields:  []any{test.target, test.field},
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

func TestLessOrEqualField(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		target      *int
		field       *int
		expected    bool
	}{
		{
			description: "TargetIsNil",
			target:      nil,
			field:       pointerOf(10),
			expected:    false,
		},
		{
			description: "FieldIsNil",
			target:      pointerOf(10),
			field:       nil,
			expected:    false,
		},
		{
			description: "Less",
			target:      pointerOf(5),
			field:       pointerOf(10),
			expected:    true,
		},
		{
			description: "Equal",
			target:      pointerOf(10),
			field:       pointerOf(10),
			expected:    true,
		},
		{
			description: "NotLessOrEqual",
			target:      pointerOf(20),
			field:       pointerOf(10),
			expected:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Act
			result := validate.LessOrEqualField(test.target, test.field)

			// Assert
			expectedResult := validate.Field{
				Valid:   test.expected,
				Message: "{0} should be less or equal than {1}",
				Fields:  []any{test.target, test.field},
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

func TestRequired(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		target      *int
		expected    bool
	}{
		{
			description: "NilPointer",
			target:      nil,
			expected:    false,
		},
		{
			description: "NotNilPointer",
			target:      pointerOf(10),
			expected:    true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Act
			result := validate.Required(test.target)

			// Assert
			expectedResult := validate.Field{
				Valid:   test.expected,
				Message: "{0} is required",
				Fields:  []any{test.target},
			}

			require.Equal(t, expectedResult, result)
		})
	}
}
