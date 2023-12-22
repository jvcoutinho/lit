package validate_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/jvcoutinho/lit/validate"
	"github.com/stretchr/testify/require"
)

func TestAfter(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		description string
		target      *time.Time
		value       time.Time
		expected    bool
	}{
		{
			description: "NilPointer",
			target:      nil,
			value:       now,
			expected:    false,
		},
		{
			description: "After",
			target:      &now,
			value:       now.Add(-20 * time.Second),
			expected:    true,
		},
		{
			description: "NotAfter",
			target:      &now,
			value:       now.Add(20 * time.Second),
			expected:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Act
			result := validate.After(test.target, test.value)

			// Assert
			expectedResult := validate.Field{
				Valid:   test.expected,
				Message: fmt.Sprintf("{0} should be after %s", test.value.Format(time.RFC3339)),
				Fields:  []any{test.target},
			}

			require.Equal(t, expectedResult, result)
		})
	}
}

func TestBefore(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		description string
		target      *time.Time
		value       time.Time
		expected    bool
	}{
		{
			description: "NilPointer",
			target:      nil,
			value:       now,
			expected:    false,
		},
		{
			description: "Before",
			target:      &now,
			value:       now.Add(20 * time.Second),
			expected:    true,
		},
		{
			description: "NotBefore",
			target:      &now,
			value:       now.Add(-20 * time.Second),
			expected:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Act
			result := validate.Before(test.target, test.value)

			// Assert
			expectedResult := validate.Field{
				Valid:   test.expected,
				Message: fmt.Sprintf("{0} should be before %s", test.value.Format(time.RFC3339)),
				Fields:  []any{test.target},
			}

			require.Equal(t, expectedResult, result)
		})
	}
}

func TestBetweenTime(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		description string
		target      *time.Time
		min         time.Time
		max         time.Time
		expected    bool
	}{
		{
			description: "NilPointer",
			target:      nil,
			min:         now,
			max:         now,
			expected:    false,
		},
		{
			description: "BeforeMin",
			target:      &now,
			min:         time.Now().Add(20 * time.Second),
			max:         time.Now().Add(30 * time.Second),
			expected:    false,
		},
		{
			description: "AfterMax",
			target:      &now,
			min:         time.Now().Add(-20 * time.Second),
			max:         time.Now().Add(-10 * time.Second),
			expected:    false,
		},
		{
			description: "Between",
			target:      &now,
			min:         time.Now().Add(-20 * time.Second),
			max:         time.Now().Add(20 * time.Second),
			expected:    true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Act
			result := validate.BetweenTime(test.target, test.min, test.max)

			// Assert
			expectedResult := validate.Field{
				Valid: test.expected,
				Message: fmt.Sprintf("{0} should be after %s and before %s",
					test.min.Format(time.RFC3339), test.max.Format(time.RFC3339)),
				Fields: []any{test.target},
			}

			require.Equal(t, expectedResult, result)
		})
	}
}
