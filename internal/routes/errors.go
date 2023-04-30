package routes

import "fmt"

// ErrRouteAlreadyDefined is the error for already defined routes.
type ErrRouteAlreadyDefined struct {
	Route Route
}

func (e ErrRouteAlreadyDefined) Error() string {
	return fmt.Sprintf("%s has been already defined", e.Route)
}

// ErrDuplicateArguments is the error for duplicate arguments in the same pattern.
//
// Example: /users/:id/books/:id
type ErrDuplicateArguments struct {
	Duplicate string
}

func (e ErrDuplicateArguments) Error() string {
	return fmt.Sprintf("a pattern can not contain two arguments with the same name (%s)", e.Duplicate)
}
