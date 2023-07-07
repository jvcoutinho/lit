package trie

// Node represents an HTTP method or a segment in a URL pattern.
//
// A node can have multiple static children, but only one dynamic child.
type Node struct {
	// Either an HTTP method or a static segment in a URL pattern.
	StaticChildren map[string]*Node
	// A parameter segment in an URL pattern.
	DynamicChild *Node
}

// NewNode creates a new Node instance.
func NewNode() *Node {
	return &Node{
		StaticChildren: make(map[string]*Node),
		DynamicChild:   nil,
	}
}

// IsTerminal returns true if this node is a leaf (should represent an HTTP method node).
func (n *Node) IsTerminal() bool {
	return len(n.StaticChildren) == 0 && n.DynamicChild == nil
}
