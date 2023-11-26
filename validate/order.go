package validate

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

// Greater validates if target is greater than value.
func Greater[T constraints.Ordered](target *T, value T) Field {
	return Field{
		Valid:   target == nil || *target > value,
		Message: fmt.Sprintf("{0} should be greater than %v", value),
		Fields:  []any{target},
	}
}

// GreaterOrEqual validates if target is greater or equal than value.
func GreaterOrEqual[T constraints.Ordered](target *T, value T) Field {
	return Field{
		Valid:   target == nil || *target >= value,
		Message: fmt.Sprintf("{0} should be greater or equal than %v", value),
		Fields:  []any{target},
	}
}

// Less validates if target is less than value.
func Less[T constraints.Ordered](target *T, value T) Field {
	return Field{
		Valid:   target == nil || *target < value,
		Message: fmt.Sprintf("{0} should be less than %v", value),
		Fields:  []any{target},
	}
}

// LessOrEqual validates if target is less or equal than value.
func LessOrEqual[T constraints.Ordered](target *T, value T) Field {
	return Field{
		Valid:   target == nil || *target <= value,
		Message: fmt.Sprintf("{0} should be less or equal than %v", value),
		Fields:  []any{target},
	}
}

// Between validates if target is greater than min and less than max.
func Between[T constraints.Ordered](target *T, min, max T) Field {
	return Field{
		Valid:   target == nil || (*target > min && *target < max),
		Message: fmt.Sprintf("{0} should be between %v and %v", min, max),
		Fields:  []any{target},
	}
}
