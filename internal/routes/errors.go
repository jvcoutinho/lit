package routes

import "fmt"

// RouteAlreadyDefinedError is the error for already defined routes.
type RouteAlreadyDefinedError struct {
	Route Route
}

func (e RouteAlreadyDefinedError) Error() string {
	return fmt.Sprintf("%s has been already defined", e.Route)
}

// DuplicateArgumentsError is the error for duplicate arguments in the same pattern.
//
// Example: /users/:id/books/:id.
type DuplicateArgumentsError struct {
	Duplicate string
}

func (e DuplicateArgumentsError) Error() string {
	return fmt.Sprintf("a pattern can not contain two arguments with the same name (%s)", e.Duplicate)
}
