package validate

import (
	"fmt"
	"strings"
)

// HasPrefix validates if target starts with prefix.
func HasPrefix(target *string, prefix string) Field {
	return Field{
		Valid:   target != nil && strings.HasPrefix(*target, prefix),
		Message: fmt.Sprintf(`{0} should start with "%s"`, prefix),
		Fields:  []any{target},
	}
}

// HasSuffix validates if target ends with suffix.
func HasSuffix(target *string, suffix string) Field {
	return Field{
		Valid:   target != nil && strings.HasSuffix(*target, suffix),
		Message: fmt.Sprintf(`{0} should end with "%s"`, suffix),
		Fields:  []any{target},
	}
}

// Substring validates if target contains the given substring.
func Substring(target *string, substring string) Field {
	return Field{
		Valid:   target != nil && strings.Contains(*target, substring),
		Message: fmt.Sprintf(`{0} should contain "%s"`, substring),
		Fields:  []any{target},
	}
}

// Lowercase validates if target contains only lowercase characters.
func Lowercase(target *string) Field {
	return Field{
		Valid:   target != nil && strings.ToLower(*target) == *target,
		Message: "{0} should contain only lowercase characters",
		Fields:  []any{target},
	}
}

// Uppercase validates if target contains only uppercase characters.
func Uppercase(target *string) Field {
	return Field{
		Valid:   target != nil && strings.ToUpper(*target) == *target,
		Message: "{0} should contain only uppercase characters",
		Fields:  []any{target},
	}
}
