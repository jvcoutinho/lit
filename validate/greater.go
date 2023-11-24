package validate

import "golang.org/x/exp/constraints"

// Greater validates if target is greater than value.
func Greater[T constraints.Ordered](target *T, value T) Validation {
	return Validation{
		Valid:     target != nil && *target > value,
		Format:    "%v should be greater than %v",
		Arguments: []any{target, value},
	}
}
