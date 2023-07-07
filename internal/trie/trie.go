package trie

import (
	"errors"
	"strings"

	"github.com/jvcoutinho/lambda/maps"
	"github.com/jvcoutinho/lambda/slices"
)

var (
	ErrPatternHasBeenDefinedAlready    = errors.New("route already exists")
	ErrPatternHasConflictingParameters = errors.New("parameters are conflicting with defined ones in another route")
	ErrMatchNotFound                   = errors.New("match not found")
)

// Trie is a tree whose nodes are either a URL pattern segment or an HTTP method.
//
// The trie is constructed in such a way that every path from the root to a leaf corresponds to a route.
// Inner nodes are pattern segments and leaves are HTTP methods.
type Trie struct {
	root       *Node
	parameters map[*Node][]string
}

// New creates a new Trie.
func New() *Trie {
	return &Trie{
		root:       NewNode(),
		parameters: make(map[*Node][]string),
	}
}

// Match checks if there is a path from the root to a leaf that represents the route of pattern and method.
// It also returns the arguments matched against the parameter segments.
//
// If such path does not exist, Match returns ErrMatchNotFound.
func (t *Trie) Match(pattern, method string) (*Node, map[string]string, error) {
	var (
		segments     = getSegments(pattern, method)
		arguments    = make([]string, 0)
		terminalNode = t.findTerminal(t.root, segments, &arguments)
	)

	if terminalNode == nil {
		return nil, nil, ErrMatchNotFound
	}

	return terminalNode, maps.FromSlices(t.parameters[terminalNode], arguments), nil
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

// Insert adds a path corresponding to the route of pattern and method and returns the terminal (leaf) node.
func (t *Trie) Insert(pattern, method string) (*Node, error) {
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
	pattern = strings.TrimLeft(pattern, "/")
	method = strings.ToUpper(method)

	segments := strings.Split(pattern, "/")
	segments = append(segments, method)

	return segments
}

func isParameter(str string) bool {
	return strings.HasPrefix(str, ":")
}
