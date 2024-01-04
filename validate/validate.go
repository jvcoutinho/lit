// Package validate contains field validations for Go structs, appropriated to use with
// [*lit.Request].
//
// There are two ways a struct can be validated:
//
//   - Explicitly, by calling [Fields] and passing the validations required;
//   - Implicitly, by making the struct implement [Validatable] with a pointer receiver
//     and using the [binding functions].
//
// When a validation fails, [Fields] use the Message and Fields attributes of the [Field] struct
// to build the validation error.
//
// # Custom validations
//
// A validation is simply an instance of the struct [Field]. In order to create a new validation, is enough to just
// create your own instance, passing the arguments. You could also change only the Message field, if that
// meets your use case. Check the [package-level examples].
//
// [binding functions]: https://pkg.go.dev/github.com/jvcoutinho/lit/bind
// [package-level examples]: https://pkg.go.dev/github.com/jvcoutinho/lit/validate#pkg-examples
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
			msg.WriteString("; ")
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
