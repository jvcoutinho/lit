package validate

import (
	"fmt"
	"time"
)

// After validates if target is after value.
func After(target *time.Time, value time.Time) Field {
	return Field{
		Valid:   target != nil && target.After(value),
		Message: fmt.Sprintf("{0} should be after %s", value.Format(time.RFC3339)),
		Fields:  []any{target},
	}
}

// Before validates if target is before value.
func Before(target *time.Time, value time.Time) Field {
	return Field{
		Valid:   target != nil && target.Before(value),
		Message: fmt.Sprintf("{0} should be before %s", value.Format(time.RFC3339)),
		Fields:  []any{target},
	}
}

// BetweenTime validates if target is after min and before max.
func BetweenTime(target *time.Time, min, max time.Time) Field {
	return Field{
		Valid: target != nil && target.After(min) && target.Before(max),
		Message: fmt.Sprintf("{0} should be after %s and before %s", min.Format(time.RFC3339),
			max.Format(time.RFC3339)),
		Fields: []any{target},
	}
}
