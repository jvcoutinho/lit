package slices

// Any returns true if and only if there is an element in arr that matches predicate.
func Any[T any](arr []T, predicate func(T) bool) bool {
	for _, e := range arr {
		if predicate(e) {
			return true
		}
	}

	return false
}
