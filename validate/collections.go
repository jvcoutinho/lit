package validate

import (
	"fmt"
	"reflect"
)

// Empty validates if target is empty.
//
// If target is not a pointer to a slice, array, string or map, Empty panics.
func Empty[T any](target *T) Field {
	return Field{
		Valid:   target != nil && reflect.ValueOf(*target).Len() == 0,
		Message: "{0} should be empty",
		Fields:  []any{target},
	}
}

// NotEmpty validates if target is not empty.
//
// If target is not a pointer to a slice, array, string or map, NotEmpty panics.
func NotEmpty[T any](target *T) Field {
	return Field{
		Valid:   target != nil && reflect.ValueOf(*target).Len() > 0,
		Message: "{0} should not be empty",
		Fields:  []any{target},
	}
}

// Length validates if target has a length of value.
//
// If target is not a pointer to a slice, array, string or map, Length panics.
func Length[T any](target *T, value int) Field {
	return Field{
		Valid:   target != nil && reflect.ValueOf(*target).Len() == value,
		Message: fmt.Sprintf("{0} should have a length of %d", value),
		Fields:  []any{target},
	}
}

// MinLength validates if target has a length greater or equal than value.
//
// If target is not a pointer to a slice, array, string or map, MinLength panics.
func MinLength[T any](target *T, value int) Field {
	return Field{
		Valid:   target != nil && reflect.ValueOf(*target).Len() >= value,
		Message: fmt.Sprintf("{0} should have a length of at least %d", value),
		Fields:  []any{target},
	}
}

// MaxLength validates if target has a length less or equal than value.
//
// If target is not a pointer to a slice, array, string or map, MaxLength panics.
func MaxLength[T any](target *T, value int) Field {
	return Field{
		Valid:   target != nil && reflect.ValueOf(*target).Len() <= value,
		Message: fmt.Sprintf("{0} should have a length of at most %d", value),
		Fields:  []any{target},
	}
}

// BetweenLength validates if target has a length greater or equal than min and less or equal than max.
//
// If target is not a pointer to a slice, array, string or map, BetweenLength panics.
func BetweenLength[T any](target *T, min, max int) Field {
	var length int

	if target != nil {
		length = reflect.ValueOf(*target).Len()
	}

	return Field{
		Valid:   target != nil && length >= min && length <= max,
		Message: fmt.Sprintf("{0} should have a length between %d and %d", min, max),
		Fields:  []any{target},
	}
}
