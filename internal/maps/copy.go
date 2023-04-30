package maps

// Copy creates an identical copy of the given map.
//
// If the given map is nil, Copy also returns nil.
func Copy[T comparable, V any](_map map[T]V) map[T]V {
	if _map == nil {
		return nil
	}

	mapCopy := make(map[T]V)

	for k, v := range _map {
		mapCopy[k] = v
	}

	return mapCopy
}
