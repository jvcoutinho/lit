package slices

import "errors"

var ErrNoElementFound = errors.New("no element found")

// First returns the first element from arr that matches predicate.
//
// If no element is found, First returns ErrNoElementFound.
func First[T any](arr []T, predicate func(T) bool) (value T, err error) {
	for _, e := range arr {
		if predicate(e) {
			return e, nil
		}
	}

	return value, ErrNoElementFound
}
