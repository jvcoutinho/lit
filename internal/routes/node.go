package routes

// Node represents a segment in an URL pattern or a HTTP method.
//
// A node can have multiple static children, but only one dynamic child.
type Node struct {
	// Children nodes that are static (are not parameters), including HTTP methods.
	StaticChildren map[string]*Node
	// A child node that is dynamic (is a parameter).
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
