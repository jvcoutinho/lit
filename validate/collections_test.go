package validate_test

import (
	"fmt"
	"testing"

	"github.com/jvcoutinho/lit/validate"
	"github.com/stretchr/testify/require"
)

func TestEmpty(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		target      *any
		expected    bool
		shouldPanic bool
	}{
		{
			description: "NilPointer",
			target:      nil,
			expected:    false,
		},
		{
			description: "NotLenType",
			target:      pointerOf[any](2),
			shouldPanic: true,
		},
		{
			description: "StringEmpty",
			target:      pointerOf[any](""),
			expected:    true,
		},
		{
			description: "StringNotEmpty",
			target:      pointerOf[any]("string"),
			expected:    false,
		},
		{
			description: "SliceEmpty",
			target:      pointerOf[any]([]int{}),
			expected:    true,
		},
		{
			description: "SliceNotEmpty",
			target:      pointerOf[any]([]int{1, 2, 3}),
			expected:    false,
		},
		{
			description: "ArrayEmpty",
			target:      pointerOf[any]([0]int{}),
			expected:    true,
		},
		{
			description: "ArrayNotEmpty",
			target:      pointerOf[any]([3]int{1, 2, 3}),
			expected:    false,
		},
		{
			description: "MapEmpty",
			target:      pointerOf[any](map[int]string{}),
			expected:    true,
		},
		{
			description: "MapNotEmpty",
			target:      pointerOf[any](map[int]string{1: "1", 2: "2", 3: "3"}),
			expected:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Act
			if test.shouldPanic {
				require.Panics(t, func() { validate.Empty(test.target) })
				return
			}

			result := validate.Empty(test.target)

			// Assert
			expectedResult := validate.Field{
				Valid:   test.expected,
				Message: "{0} should be empty",
				Fields:  []any{test.target},
			}

			require.Equal(t, expectedResult, result)
		})
	}
}

func TestNotEmpty(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		target      *any
		expected    bool
		shouldPanic bool
	}{
		{
			description: "NilPointer",
			target:      nil,
			expected:    false,
		},
		{
			description: "NotLenType",
			target:      pointerOf[any](2),
			shouldPanic: true,
		},
		{
			description: "StringEmpty",
			target:      pointerOf[any](""),
			expected:    false,
		},
		{
			description: "StringNotEmpty",
			target:      pointerOf[any]("string"),
			expected:    true,
		},
		{
			description: "SliceEmpty",
			target:      pointerOf[any]([]int{}),
			expected:    false,
		},
		{
			description: "SliceNotEmpty",
			target:      pointerOf[any]([]int{1, 2, 3}),
			expected:    true,
		},
		{
			description: "ArrayEmpty",
			target:      pointerOf[any]([0]int{}),
			expected:    false,
		},
		{
			description: "ArrayNotEmpty",
			target:      pointerOf[any]([3]int{1, 2, 3}),
			expected:    true,
		},
		{
			description: "MapEmpty",
			target:      pointerOf[any](map[int]string{}),
			expected:    false,
		},
		{
			description: "MapNotEmpty",
			target:      pointerOf[any](map[int]string{1: "1", 2: "2", 3: "3"}),
			expected:    true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Act
			if test.shouldPanic {
				require.Panics(t, func() { validate.NotEmpty(test.target) })
				return
			}

			result := validate.NotEmpty(test.target)

			// Assert
			expectedResult := validate.Field{
				Valid:   test.expected,
				Message: "{0} should not be empty",
				Fields:  []any{test.target},
			}

			require.Equal(t, expectedResult, result)
		})
	}
}

func TestLength(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		target      *any
		length      int
		expected    bool
		shouldPanic bool
	}{
		{
			description: "NilPointer",
			target:      nil,
			length:      10,
			expected:    false,
		},
		{
			description: "NotLenType",
			target:      pointerOf[any](2),
			length:      10,
			shouldPanic: true,
		},
		{
			description: "StringHasLength",
			target:      pointerOf[any]("length-of-12"),
			length:      12,
			expected:    true,
		},
		{
			description: "StringDoesNotHaveLength",
			target:      pointerOf[any]("length-of-12"),
			length:      15,
			expected:    false,
		},
		{
			description: "SliceHasLength",
			target:      pointerOf[any]([]int{1, 2, 3}),
			length:      3,
			expected:    true,
		},
		{
			description: "SliceDoesNotHaveLength",
			target:      pointerOf[any]([]int{1, 2, 3}),
			length:      2,
			expected:    false,
		},
		{
			description: "ArrayHasLength",
			target:      pointerOf[any]([3]int{1, 2, 3}),
			length:      3,
			expected:    true,
		},
		{
			description: "ArrayDoesNotHaveLength",
			target:      pointerOf[any]([3]int{1, 2, 3}),
			length:      2,
			expected:    false,
		},
		{
			description: "MapHasLength",
			target:      pointerOf[any](map[int]string{1: "1", 2: "2", 3: "3"}),
			length:      3,
			expected:    true,
		},
		{
			description: "MapDoesNotHaveLength",
			target:      pointerOf[any](map[int]string{1: "1", 2: "2", 3: "3"}),
			length:      2,
			expected:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Act
			if test.shouldPanic {
				require.Panics(t, func() { validate.Length(test.target, test.length) })
				return
			}

			result := validate.Length(test.target, test.length)

			// Assert
			expectedResult := validate.Field{
				Valid:   test.expected,
				Message: fmt.Sprintf(`{0} should have a length of %d`, test.length),
				Fields:  []any{test.target},
			}

			require.Equal(t, expectedResult, result)
		})
	}
}

func TestMinLength(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		target      *any
		length      int
		expected    bool
		shouldPanic bool
	}{
		{
			description: "NilPointer",
			target:      nil,
			length:      10,
			expected:    false,
		},
		{
			description: "NotLenType",
			target:      pointerOf[any](2),
			length:      10,
			shouldPanic: true,
		},
		{
			description: "StringHasLength",
			target:      pointerOf[any]("length-of-12"),
			length:      12,
			expected:    true,
		},
		{
			description: "StringHasMoreThanLength",
			target:      pointerOf[any]("length-of-12"),
			length:      10,
			expected:    true,
		},
		{
			description: "StringHasLessThanLength",
			target:      pointerOf[any]("length-of-12"),
			length:      15,
			expected:    false,
		},
		{
			description: "SliceHasLength",
			target:      pointerOf[any]([]int{1, 2, 3}),
			length:      3,
			expected:    true,
		},
		{
			description: "SliceHasMoreThanLength",
			target:      pointerOf[any]([]int{1, 2, 3}),
			length:      1,
			expected:    true,
		},
		{
			description: "SliceHasLessThanLength",
			target:      pointerOf[any]([]int{1, 2, 3}),
			length:      4,
			expected:    false,
		},
		{
			description: "ArrayHasLength",
			target:      pointerOf[any]([3]int{1, 2, 3}),
			length:      3,
			expected:    true,
		},
		{
			description: "ArrayHasMoreThanLength",
			target:      pointerOf[any]([3]int{1, 2, 3}),
			length:      1,
			expected:    true,
		},
		{
			description: "ArrayHasLessThanLength",
			target:      pointerOf[any]([3]int{1, 2, 3}),
			length:      4,
			expected:    false,
		},
		{
			description: "MapHasLength",
			target:      pointerOf[any](map[int]string{1: "1", 2: "2", 3: "3"}),
			length:      3,
			expected:    true,
		},
		{
			description: "MapHasMoreThanLength",
			target:      pointerOf[any](map[int]string{1: "1", 2: "2", 3: "3"}),
			length:      1,
			expected:    true,
		},
		{
			description: "MapHasLessThanLength",
			target:      pointerOf[any](map[int]string{1: "1", 2: "2", 3: "3"}),
			length:      4,
			expected:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Act
			if test.shouldPanic {
				require.Panics(t, func() { validate.MinLength(test.target, test.length) })
				return
			}

			result := validate.MinLength(test.target, test.length)

			// Assert
			expectedResult := validate.Field{
				Valid:   test.expected,
				Message: fmt.Sprintf(`{0} should have a length of at least %d`, test.length),
				Fields:  []any{test.target},
			}

			require.Equal(t, expectedResult, result)
		})
	}
}

func TestMaxLength(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		target      *any
		length      int
		expected    bool
		shouldPanic bool
	}{
		{
			description: "NilPointer",
			target:      nil,
			length:      10,
			expected:    false,
		},
		{
			description: "NotLenType",
			target:      pointerOf[any](2),
			length:      10,
			shouldPanic: true,
		},
		{
			description: "StringHasLength",
			target:      pointerOf[any]("length-of-12"),
			length:      12,
			expected:    true,
		},
		{
			description: "StringHasMoreThanLength",
			target:      pointerOf[any]("length-of-12"),
			length:      10,
			expected:    false,
		},
		{
			description: "StringHasLessThanLength",
			target:      pointerOf[any]("length-of-12"),
			length:      15,
			expected:    true,
		},
		{
			description: "SliceHasLength",
			target:      pointerOf[any]([]int{1, 2, 3}),
			length:      3,
			expected:    true,
		},
		{
			description: "SliceHasMoreThanLength",
			target:      pointerOf[any]([]int{1, 2, 3}),
			length:      1,
			expected:    false,
		},
		{
			description: "SliceHasLessThanLength",
			target:      pointerOf[any]([]int{1, 2, 3}),
			length:      4,
			expected:    true,
		},
		{
			description: "ArrayHasLength",
			target:      pointerOf[any]([3]int{1, 2, 3}),
			length:      3,
			expected:    true,
		},
		{
			description: "ArrayHasMoreThanLength",
			target:      pointerOf[any]([3]int{1, 2, 3}),
			length:      1,
			expected:    false,
		},
		{
			description: "ArrayHasLessThanLength",
			target:      pointerOf[any]([3]int{1, 2, 3}),
			length:      4,
			expected:    true,
		},
		{
			description: "MapHasLength",
			target:      pointerOf[any](map[int]string{1: "1", 2: "2", 3: "3"}),
			length:      3,
			expected:    true,
		},
		{
			description: "MapHasMoreThanLength",
			target:      pointerOf[any](map[int]string{1: "1", 2: "2", 3: "3"}),
			length:      1,
			expected:    false,
		},
		{
			description: "MapHasLessThanLength",
			target:      pointerOf[any](map[int]string{1: "1", 2: "2", 3: "3"}),
			length:      4,
			expected:    true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Act
			if test.shouldPanic {
				require.Panics(t, func() { validate.MaxLength(test.target, test.length) })
				return
			}

			result := validate.MaxLength(test.target, test.length)

			// Assert
			expectedResult := validate.Field{
				Valid:   test.expected,
				Message: fmt.Sprintf(`{0} should have a length of at most %d`, test.length),
				Fields:  []any{test.target},
			}

			require.Equal(t, expectedResult, result)
		})
	}
}

func TestBetweenLength(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		target      *any
		min         int
		max         int
		expected    bool
		shouldPanic bool
	}{
		{
			description: "NilPointer",
			target:      nil,
			min:         10,
			max:         10,
			expected:    false,
		},
		{
			description: "NotLenType",
			target:      pointerOf[any](2),
			min:         10,
			max:         10,
			shouldPanic: true,
		},
		{
			description: "StringHasLengthOfMin",
			target:      pointerOf[any]("length-of-12"),
			min:         12,
			max:         15,
			expected:    true,
		},
		{
			description: "StringHasLengthOfMax",
			target:      pointerOf[any]("length-of-12"),
			min:         10,
			max:         12,
			expected:    true,
		},
		{
			description: "StringHasLengthBetweenMinAndMax",
			target:      pointerOf[any]("length-of-12"),
			min:         10,
			max:         13,
			expected:    true,
		},
		{
			description: "StringHasMoreThanMax",
			target:      pointerOf[any]("length-of-12"),
			min:         10,
			max:         11,
			expected:    false,
		},
		{
			description: "StringHasLessThanMin",
			target:      pointerOf[any]("length-of-12"),
			min:         13,
			max:         15,
			expected:    false,
		},
		{
			description: "SliceHasLengthOfMin",
			target:      pointerOf[any]([]int{1, 2, 3}),
			min:         3,
			max:         15,
			expected:    true,
		},
		{
			description: "SliceHasLengthOfMax",
			target:      pointerOf[any]([]int{1, 2, 3}),
			min:         1,
			max:         3,
			expected:    true,
		},
		{
			description: "SliceHasLengthBetweenMinAndMax",
			target:      pointerOf[any]([]int{1, 2, 3}),
			min:         1,
			max:         4,
			expected:    true,
		},
		{
			description: "SliceHasMoreThanMax",
			target:      pointerOf[any]([]int{1, 2, 3}),
			min:         1,
			max:         2,
			expected:    false,
		},
		{
			description: "SliceHasLessThanMin",
			target:      pointerOf[any]([]int{1, 2, 3}),
			min:         4,
			max:         5,
			expected:    false,
		},
		{
			description: "ArrayHasLengthOfMin",
			target:      pointerOf[any]([3]int{1, 2, 3}),
			min:         3,
			max:         15,
			expected:    true,
		},
		{
			description: "ArrayHasLengthOfMax",
			target:      pointerOf[any]([3]int{1, 2, 3}),
			min:         1,
			max:         3,
			expected:    true,
		},
		{
			description: "ArrayHasLengthBetweenMinAndMax",
			target:      pointerOf[any]([3]int{1, 2, 3}),
			min:         1,
			max:         4,
			expected:    true,
		},
		{
			description: "ArrayHasMoreThanMax",
			target:      pointerOf[any]([3]int{1, 2, 3}),
			min:         1,
			max:         2,
			expected:    false,
		},
		{
			description: "ArrayHasLessThanMin",
			target:      pointerOf[any]([3]int{1, 2, 3}),
			min:         4,
			max:         5,
			expected:    false,
		},
		{
			description: "MapHasLengthOfMin",
			target:      pointerOf[any](map[int]string{1: "1", 2: "2", 3: "3"}),
			min:         3,
			max:         15,
			expected:    true,
		},
		{
			description: "MapHasLengthOfMax",
			target:      pointerOf[any](map[int]string{1: "1", 2: "2", 3: "3"}),
			min:         1,
			max:         3,
			expected:    true,
		},
		{
			description: "MapHasLengthBetweenMinAndMax",
			target:      pointerOf[any](map[int]string{1: "1", 2: "2", 3: "3"}),
			min:         1,
			max:         4,
			expected:    true,
		},
		{
			description: "MapHasMoreThanMax",
			target:      pointerOf[any](map[int]string{1: "1", 2: "2", 3: "3"}),
			min:         1,
			max:         2,
			expected:    false,
		},
		{
			description: "MapHasLessThanMin",
			target:      pointerOf[any](map[int]string{1: "1", 2: "2", 3: "3"}),
			min:         4,
			max:         5,
			expected:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Act
			if test.shouldPanic {
				require.Panics(t, func() { validate.BetweenLength(test.target, test.min, test.max) })
				return
			}

			result := validate.BetweenLength(test.target, test.min, test.max)

			// Assert
			expectedResult := validate.Field{
				Valid:   test.expected,
				Message: fmt.Sprintf(`{0} should have a length between %d and %d`, test.min, test.max),
				Fields:  []any{test.target},
			}

			require.Equal(t, expectedResult, result)
		})
	}
}
