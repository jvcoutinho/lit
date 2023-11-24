package validate

import "fmt"

// Validatable can be validated.
type Validatable interface {
	// Validate returns a validation of itself.
	Validate() Validation
}

// Validation represents an assertion of a condition.
type Validation struct {
	// Determines if this validation succeeded.
	Valid bool
	// A format that can be displayed to the user if the validation fails.
	Format string
	// Arguments of this validation. It can contain struct fields as well as argument values.
	Arguments []any
}

func (v Validation) Error() string {
	return fmt.Sprintf(v.Format, v.Arguments...)
}
