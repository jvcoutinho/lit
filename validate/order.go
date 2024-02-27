package validate

import (
	"fmt"
	"slices"

	"golang.org/x/exp/constraints"
)

// Equal validates if target is equal to value.
func Equal[T comparable](target *T, value T) Field {
	return Field{
		Valid:   target != nil && *target == value,
		Message: fmt.Sprintf("{0} should be equal to %v", value),
		Fields:  []any{target},
	}
}

// NotEqual validates if target is not equal to value.
func NotEqual[T comparable](target *T, value T) Field {
	return Field{
		Valid:   target != nil && *target != value,
		Message: fmt.Sprintf("{0} should not be equal to %v", value),
		Fields:  []any{target},
	}
}

// EqualField validates if target is equal to the value of field.
func EqualField[T comparable](target *T, field *T) Field {
	return Field{
		Valid:   (target == nil && field == nil) || (target != nil && field != nil && *target == *field),
		Message: "{0} should be equal to {1}",
		Fields:  []any{target, field},
	}
}

// NotEqualField validates if target is not equal to the value of field.
func NotEqualField[T comparable](target *T, field *T) Field {
	return Field{
		Valid: (target != nil && field == nil) || (field != nil && target == nil) || (field != nil &&
			*target != *field),
		Message: "{0} should not be equal to {1}",
		Fields:  []any{target, field},
	}
}

// OneOf validates if target is one of values.
func OneOf[T comparable](target *T, values ...T) Field {
	return Field{
		Valid:   target != nil && slices.Contains(values, *target),
		Message: fmt.Sprintf("{0} should be one of %v", values),
		Fields:  []any{target},
	}
}

// Greater validates if target is greater than value.
func Greater[T constraints.Ordered](target *T, value T) Field {
	return Field{
		Valid:   target != nil && *target > value,
		Message: fmt.Sprintf("{0} should be greater than %v", value),
		Fields:  []any{target},
	}
}

// GreaterField validates if target is greater than the value of field.
func GreaterField[T constraints.Ordered](target *T, field *T) Field {
	return Field{
		Valid:   target != nil && field != nil && *target > *field,
		Message: "{0} should be greater than {1}",
		Fields:  []any{target, field},
	}
}

// GreaterOrEqual validates if target is greater or equal than value.
func GreaterOrEqual[T constraints.Ordered](target *T, value T) Field {
	return Field{
		Valid:   target != nil && *target >= value,
		Message: fmt.Sprintf("{0} should be greater or equal than %v", value),
		Fields:  []any{target},
	}
}

// GreaterOrEqualField validates if target is greater or equal than the value of field.
func GreaterOrEqualField[T constraints.Ordered](target *T, field *T) Field {
	return Field{
		Valid:   target != nil && field != nil && *target >= *field,
		Message: "{0} should be greater or equal than {1}",
		Fields:  []any{target, field},
	}
}

// Less validates if target is less than value.
func Less[T constraints.Ordered](target *T, value T) Field {
	return Field{
		Valid:   target != nil && *target < value,
		Message: fmt.Sprintf("{0} should be less than %v", value),
		Fields:  []any{target},
	}
}

// LessField validates if target is less than the value of field.
func LessField[T constraints.Ordered](target *T, field *T) Field {
	return Field{
		Valid:   target != nil && field != nil && *target < *field,
		Message: "{0} should be less than {1}",
		Fields:  []any{target, field},
	}
}

// LessOrEqual validates if target is less or equal than value.
func LessOrEqual[T constraints.Ordered](target *T, value T) Field {
	return Field{
		Valid:   target != nil && *target <= value,
		Message: fmt.Sprintf("{0} should be less or equal than %v", value),
		Fields:  []any{target},
	}
}

// LessOrEqualField validates if target is less or equal than the value of field.
func LessOrEqualField[T constraints.Ordered](target *T, field *T) Field {
	return Field{
		Valid:   target != nil && field != nil && *target <= *field,
		Message: "{0} should be less or equal than {1}",
		Fields:  []any{target, field},
	}
}

// Between validates if target is greater than min and less than max.
func Between[T constraints.Ordered](target *T, min, max T) Field {
	return Field{
		Valid:   target != nil && (*target > min && *target < max),
		Message: fmt.Sprintf("{0} should be between %v and %v", min, max),
		Fields:  []any{target},
	}
}

// Required validates that target is not nil. Suited for when target is a pointer field.
func Required[T any](target *T) Field {
	return Field{
		Valid:   target != nil,
		Message: "{0} is required",
		Fields:  []any{target},
	}
}
