package slices

// ElementAtOrDefault returns the element of arr pointed by index.
//
// If index is out of range, the default value of T is returned.
func ElementAtOrDefault[T any](arr []T, index int) T {
	var value T

	if index >= 0 && index < len(arr) {
		value = arr[index]
	}

	return value
}
