package validate

import "strings"

// Validatable structs can be validated.
type Validatable interface {
	// Validate returns a list of field validations.
	Validate() []Field
}

// Error occurs when at least one Field validation fails in validation steps.
type Error struct {
	// Validations that have failed.
	Violations []Field
}

func (e Error) Error() string {
	msg := strings.Builder{}

	for i, violation := range e.Violations {
		msg.WriteString(violation.Message)

		if i < len(e.Violations)-1 {
			msg.WriteString(", ")
		}
	}

	return msg.String()
}

// Field represents a field validation.
type Field struct {
	// Determines if this validation has succeeded.
	Valid bool

	// A user-friendly message that can be displayed if the validation fails. In this case,
	// "{i}" placeholders, where i is the index of the replacement in Fields, are replaced by the field's tag values.
	Message string

	// Pointers to the fields involved in this validation. It is a slice because the validation
	// can have multiple fields as arguments.
	Fields []any
}
