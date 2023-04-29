package slices

// First returns the first element from arr that matches predicate.
//
// If no element is found, ok is false.
func First[T any](arr []T, predicate func(T) bool) (value T, ok bool) {
	for _, e := range arr {
		if predicate(e) {
			return e, true
		}
	}

	return value, false
}
