package main

import "strings"

// Node is the parse tree of a template.
type Node struct {
	Token    Token
	List     []*Node
	Parent   *Node
	MinusEnd int // did the {{end}} tag contain {{- or -}}
}

// Parse parses the tokens and adds them the list in n. It returns the token that are not yet consumed.
func (n *Node) parse(tokens []Token) []Token {
	if len(tokens) == 0 {
		return nil
	}

	switch tokens[0].Subtype {
	case End:
		// Add to current list and then proceed to parse again a level higher. But don't add the token itself.
		n.MinusEnd = minus(tokens[0].Value)
		// a single end can close a chain of else ifs, we need to find the ultimate if/when/range parent here
		parent := n.Parent
		for parent != nil {
			if parent.Token.Subtype != Else &&
				parent.Token.Subtype != If &&
				parent.Token.Subtype != Range &&
				parent.Token.Subtype != With {
				break
			}
			parent = parent.Parent
		}
		if parent == nil {
			return n.parse(tokens[1:])
		} else {
			return parent.parse(tokens[1:])
		}

	case TagClose:
		n1 := &Node{Token: tokens[0], Parent: n}
		n.List = append(n.List, n1)

		if n.Parent == nil { // can happen when we see just the close tag in the root of the node.
			return n.parse(tokens[1:])
		} else {
			return n.Parent.parse(tokens[1:])
		}

	case Define, If, Else, Block, Range, With:
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
	tree := &Node{}
	tree.parse(tokens)

	return tree
}

const (
	MinusNone = iota
	MinusLeft
	MinusRight
	MinusBoth
)

func minus(s string) int {
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
