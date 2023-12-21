package validate

import (
	"fmt"
	"regexp"
	"time"
)

var (
	uuidRegex        = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	emailRegex       = regexp.MustCompile(`^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`)
	ipv4AddressRegex = regexp.MustCompile(`(\b25[0-5]|\b2[0-4][0-9]|\b[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
	ipv6AddressRegex = regexp.MustCompile(`(([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]+|::(ffff(:0{1,4})?:)?((25[0-5]|(2[0-4]|1?[0-9])?[0-9])\.){3}(25[0-5]|(2[0-4]|1?[0-9])?[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1?[0-9])?[0-9])\.){3}(25[0-5]|(2[0-4]|1?[0-9])?[0-9]))`)
)

// UUID validates if target is a valid [UUID] of any version.
//
// [UUID]: https://en.wikipedia.org/wiki/Universally_unique_identifier
func UUID(target *string) Field {
	return Field{
		Valid:   target != nil && uuidRegex.MatchString(*target),
		Message: "{0} is not a valid UUID",
		Fields:  []any{target},
	}
}

// Email validates if target is a valid e-mail.
func Email(target *string) Field {
	return Field{
		Valid:   target != nil && emailRegex.MatchString(*target),
		Message: "{0} is not a valid e-mail",
		Fields:  []any{target},
	}
}

// IPAddress validates if target is a valid IPv4 or IPv6 address.
func IPAddress(target *string) Field {
	return Field{
		Valid:   target != nil && (ipv4AddressRegex.MatchString(*target) || ipv6AddressRegex.MatchString(*target)),
		Message: "{0} is not a valid IP address",
		Fields:  []any{target},
	}
}

// IPv4Address validates if target is a valid IPv4 address.
func IPv4Address(target *string) Field {
	return Field{
		Valid:   target != nil && ipv4AddressRegex.MatchString(*target),
		Message: "{0} is not a valid IPv4 address",
		Fields:  []any{target},
	}
}

// IPv6Address validates if target is a valid IPv6 address.
func IPv6Address(target *string) Field {
	return Field{
		Valid:   target != nil && ipv6AddressRegex.MatchString(*target),
		Message: "{0} is not a valid IPv6 address",
		Fields:  []any{target},
	}
}

// DateTime validates if target is a valid date-time string.
//
// Note that binding functions bind strings to [time.Time] fields if the layout is time.RFC3339, validating
// them in the process. In this case, this function is not needed.
func DateTime(target *string, layout string) Field {
	var err error
	if target != nil {
		_, err = time.Parse(layout, *target)
	}

	return Field{
		Valid:   target != nil && err == nil,
		Message: fmt.Sprintf(`{0} is not a valid date time in required format (ex: "%s")`, layout),
		Fields:  []any{target},
	}
}
