package routes

import (
	"errors"
	"strings"

	"github.com/jvcoutinho/lambda/maps"
	"github.com/jvcoutinho/lambda/slices"
)

var (
	ErrMethodIsEmpty                   = errors.New("method should not be empty")
	ErrPatternContainsDoubleSlash      = errors.New("pattern should not contain double slashes (//)")
	ErrPatternHasBeenDefinedAlready    = errors.New("route already exists")
	ErrPatternHasConflictingParameters = errors.New("parameters are conflicting with defined ones in another route")
)

// Trie stores pattern segments as nodes so they can be matched with incoming requests later.
type Trie struct {
	root       *Node
	parameters map[*Node][]string
}

// NewTrie creates a new Trie instance.
func NewTrie() *Trie {
	return &Trie{
		root:       NewNode(),
		parameters: make(map[*Node][]string),
	}
}

// Match matches pattern and method against the nodes in this trie, returning the node corresponding
// to the HTTP method and the matched arguments.
//
// If a match can not be found, Match returns (nil, nil).
func (t *Trie) Match(pattern, method string) (*Node, map[string]string) {
	var (
		segments     = getSegments(pattern, method)
		arguments    = make([]string, 0)
		terminalNode = t.findTerminal(t.root, segments, &arguments)
	)

	if terminalNode == nil {
		return nil, nil
	}

	return terminalNode, maps.FromSlices(t.parameters[terminalNode], arguments)
}

func (t *Trie) findTerminal(parent *Node, segments []string, arguments *[]string) *Node {
	if len(segments) == 0 {
		if parent.IsTerminal() {
			return parent
		}

		return nil
	}

	segment := segments[0]

	if child, ok := parent.StaticChildren[segment]; ok {
		if terminal := t.findTerminal(child, segments[1:], arguments); terminal != nil {
			return terminal
		}
	}

	if parent.DynamicChild != nil {
		*arguments = append(*arguments, segment)

		if terminal := t.findTerminal(parent.DynamicChild, segments[1:], arguments); terminal != nil {
			return terminal
		}

		*arguments = (*arguments)[:len(*arguments)-1]
	}

	return nil
}

// Insert adds a pattern and method pair to this trie, returning the terminal node.
//
// If a pattern already exists for method in this trie, Insert returns an error.
func (t *Trie) Insert(pattern, method string) (*Node, error) {
	if method == "" {
		return nil, ErrMethodIsEmpty
	}

	if strings.Contains(pattern, "//") {
		return nil, ErrPatternContainsDoubleSlash
	}

	segments := getSegments(pattern, method)

	if err := t.validate(segments); err != nil {
		return nil, err
	}

	var (
		parent     = t.root
		parameters = make([]string, 0)
	)

	for _, segment := range segments {
		if isParameter(segment) {
			parameters = append(parameters, segment)

			if parent.DynamicChild != nil {
				parent = parent.DynamicChild

				continue
			}

			parent.DynamicChild = NewNode()
			parent = parent.DynamicChild

			continue
		}

		if child, ok := parent.StaticChildren[segment]; ok {
			parent = child

			continue
		}

		parent.StaticChildren[segment] = NewNode()
		parent = parent.StaticChildren[segment]
	}

	t.parameters[parent] = parameters

	return parent, nil
}

func (t *Trie) validate(segments []string) error {
	var (
		parent          = t.root
		foundParameters = make([]string, 0)
	)

	for _, segment := range segments {
		if isParameter(segment) {
			if parent.DynamicChild == nil {
				return nil
			}

			foundParameters = append(foundParameters, segment)
			parent = parent.DynamicChild

			continue
		}

		child, ok := parent.StaticChildren[segment]
		if !ok {
			return nil
		}

		parent = child
	}

	routeParameters, ok := t.parameters[parent]
	if ok && !slices.Equal(routeParameters, foundParameters) {
		return ErrPatternHasConflictingParameters
	}

	return ErrPatternHasBeenDefinedAlready
}

func getSegments(pattern, method string) []string {
	pattern = strings.Trim(pattern, "/")
	method = strings.ToUpper(method)

	segments := strings.Split(pattern, "/")
	segments = append(segments, method)

	return segments
}

func isParameter(str string) bool {
	return strings.HasPrefix(str, ":")
}
