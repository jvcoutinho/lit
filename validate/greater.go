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
		Targets: []any{target},
	}
}
