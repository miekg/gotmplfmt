package main

import "strings"

// Node is the parse tree of a template.
type Node struct {
	Token    Token
	List     []*Node
	Parent   *Node
	MinusEnd int // did the {{end}} tag contain {{- or -}}
}

func NewNode(parent *Node) *Node { return &Node{Parent: parent} }

// Parse parses the tokens and adds them the list in n. It returns the token that are not yet consumed.
func (n *Node) parse(tokens []Token) []Token {
	if len(tokens) == 0 {
		return nil
	}

	switch tokens[0].Subtype {
	case End:
		// Add to current list and then proceed to parse again a level higher. But don't add the token itself.
		n.MinusEnd = Minus(tokens[0].Value)
		return n.Parent.parse(tokens[1:])
	case TagClose:
		if n.Token.Type == TokenText /* 0 */ && n.Token.Value == "" {
			// Never seen an open for this tag, add to the tree, so it's outputted automatically.
			n1 := &Node{Token: tokens[0], Parent: n}
			n.List = append(n.List, n1)
		}
		if n.Parent == nil { // can happen when we see just the close tag in the root of the node.
			return n.parse(tokens[1:])
		} else {
			return n.Parent.parse(tokens[1:])
		}

	case Define, If, ElseIf, Else, Block, Range, With:
		n1 := &Node{Token: tokens[0], Parent: n}
		n.List = append(n.List, n1)
		return n1.parse(tokens[1:])

	case TagOpen:
		n1 := &Node{Token: tokens[0], Parent: n}
		n.List = append(n.List, n1)
		return n1.parse(tokens[1:])
	}

	n1 := &Node{Token: tokens[0], Parent: n}
	n.List = append(n.List, n1)
	return n.parse(tokens[1:])
}

// Parse parses tokens and returns to root node of the tree.
func Parse(tokens []Token) *Node {
	tree := NewNode(nil)
	tree.parse(tokens)

	return tree
}

const (
	MinusNone = iota
	MinusLeft
	MinusRight
	MinusBoth
)

func Minus(s string) int {
	if strings.HasPrefix(s, "{{-") {
		if strings.HasSuffix(s, "-}}") {
			return MinusBoth
		}
		return MinusLeft
	}
	if strings.HasSuffix(s, "-}}") {
		return MinusRight
	}
	return MinusNone
}
