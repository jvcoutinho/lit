package slices

// Contains returns true if and only if arr contains elem.
func Contains[T comparable](arr []T, elem T) bool {
	for _, e := range arr {
		if e == elem {
			return true
		}
	}

	return false
}
