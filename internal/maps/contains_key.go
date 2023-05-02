package maps

// ContainsKey returns true if and only if elem is a key defined in _map.
func ContainsKey[T comparable, V any](_map map[T]V, elem T) bool {
	_, ok := _map[elem]

	return ok
}
