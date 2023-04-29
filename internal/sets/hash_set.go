package sets

// HashSet is a collection for unordered unique elements.
type HashSet[T comparable] map[T]bool

// NewHashSet creates a new instance of HashSet.
func NewHashSet[T comparable]() HashSet[T] {
	return make(map[T]bool)
}

// Add adds an element to this set.
func (s HashSet[T]) Add(elem T) {
	s[elem] = true
}

// Contains returns true if and only if elem is contained in this set.
func (s HashSet[T]) Contains(elem T) bool {
	_, ok := s[elem]
	return ok
}
