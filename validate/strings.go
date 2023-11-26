package validate

import (
	"fmt"
	"strings"
)

// HasPrefix validates if target starts with prefix.
func HasPrefix(target *string, prefix string) Field {
	return Field{
		Valid:   target == nil || strings.HasPrefix(*target, prefix),
		Message: fmt.Sprintf(`{0} should start with "%s"`, prefix),
		Fields:  []any{target},
	}
}
