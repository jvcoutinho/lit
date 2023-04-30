package structures

// HashSet is a collection for unordered unique elements.
type HashSet[T comparable] map[T]bool

// NewHashSet creates a new instance of HashSet.
func NewHashSet[T comparable](initialElements ...T) HashSet[T] {
	set := make(map[T]bool)
	for _, element := range initialElements {
		set[element] = true
	}

	return set
}

// Add adds an element to this set.
func (s HashSet[T]) Add(elem T) {
	s[elem] = true
}

// Len is the number of elements this set contains.
func (s HashSet[T]) Len() int {
	return len(s)
}

// Contains returns true if and only if elem is contained in this set.
func (s HashSet[T]) Contains(elem T) bool {
	_, ok := s[elem]
	return ok
}
