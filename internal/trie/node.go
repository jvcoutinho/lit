package trie

// Node represents a HTTP method or a segment in an URL pattern.
//
// A node can have multiple static children, but only one dynamic child.
type Node struct {
	// Either a HTTP method or a static segment in an URL pattern.
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

// IsTerminal returns true if this node is a leaf (should represent a HTTP method node).
func (n *Node) IsTerminal() bool {
	return len(n.StaticChildren) == 0 && n.DynamicChild == nil
}
