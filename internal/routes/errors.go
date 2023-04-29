package routes

import "fmt"

// ErrDuplicateArguments is the error for duplicate arguments in the same pattern.
//
// Example: /users/:id/books/:id
type ErrDuplicateArguments struct {
	Duplicate string
}

func (e ErrDuplicateArguments) Error() string {
	return fmt.Sprintf("a pattern can not contain two arguments with the same name (%s)", e.Duplicate)
}
