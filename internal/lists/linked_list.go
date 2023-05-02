package lists

type listNode[T any] struct {
	value T
	next  *listNode[T]
}

func newListNode[T any](value T) *listNode[T] {
	return &listNode[T]{value, nil}
}

// LinkedList is a structure suited to store elements without knowing in advance how much they will be.
type LinkedList[T any] struct {
	head *listNode[T]
	len  int
}

// NewLinkedList creates a new instance of LinkedList.
func NewLinkedList[T any](initialElements ...T) *LinkedList[T] {
	list := &LinkedList[T]{}

	for i := len(initialElements) - 1; i >= 0; i-- {
		list.InsertAtBeginning(initialElements[i])
	}

	return list
}

// InsertAtBeginning inserts value to this list at its beginning.
//
// This costs O(1).
func (l *LinkedList[T]) InsertAtBeginning(value T) {
	newNode := newListNode(value)

	currentHead := l.head
	l.head = newNode
	newNode.next = currentHead

	l.len++
}

// ElementAt gets the element at the given index.
//
// If index is less than zero or is greater than the length of this list, ok equals false.
func (l *LinkedList[T]) ElementAt(index int) (value T, ok bool) {
	if index < 0 || index >= l.Len() {
		return value, false
	}

	i := 0
	for currentNode := l.head; currentNode != nil; currentNode = currentNode.next {
		if i == index {
			value = currentNode.value
			break
		}

		i++
	}

	return value, true
}

// Len is the number of elements this list has.
//
// This costs O(1).
func (l *LinkedList[T]) Len() int {
	return l.len
}

// Traverse executes action in each element of this list.
func (l *LinkedList[T]) Traverse(action func(T)) {
	for currentNode := l.head; currentNode != nil; currentNode = currentNode.next {
		action(currentNode.value)
	}
}
