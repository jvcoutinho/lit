package slices

// Filter returns all the elements of arr that match predicate.
func Filter[T comparable](arr []T, predicate func(T) bool) []T {
	result := make([]T, 0, len(arr))

	for _, e := range arr {
		if predicate(e) {
			result = append(result, e)
		}
	}

	return result
}
