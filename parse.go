package main

import (
	"fmt"
)

// Node is the parse tree of a template.
type Node struct {
	Token  Token
	List   []*Node
	Parent *Node
}

func NewNode(parent *Node) *Node { return &Node{Parent: parent} }

// Parse parses the tokens and adds them the list in n. It returns the token that are not yet consumed.
func (n *Node) parse(tokens []Token) []Token {
	if len(tokens) == 0 {
		return nil
	}

	switch tokens[0].Subtype {
	case End:
		// If n is Else of ElseIf we go up 2 because the else adds a branch, but we only get one
		// {{end}}
		n1 := &Node{Token: tokens[0], Parent: n}

		if n.Token.Subtype == Else || n.Token.Subtype == ElseIf {
			n.Parent.Parent.List = append(n.Parent.Parent.List, n1)
			return n.Parent.Parent.parse(tokens[1:])
		}

		// Add to current list and then proceed to parse again a level higher.
		n.List = append(n.List, n1)
		return n.Parent.parse(tokens[1:])

	case Define, If, ElseIf, Else, Block, Range, With:
		n1 := &Node{Token: tokens[0], Parent: n}
		n.List = append(n.List, n1)
		return n1.parse(tokens[1:])

	default:
		// also includes TokenText, by having a zero subtype
		n1 := &Node{Token: tokens[0], Parent: n}
		n.List = append(n.List, n1)
		return n.parse(tokens[1:])
	}
}

// Parse parses tokens and returns to root node of the tree.
func Parse(tokens []Token) *Node {
	for _, token := range tokens {
		fmt.Printf("Type: %v, Subtype: %v, Value: %q\n", token.Type, token.Subtype, token.Value)
	}

	tree := NewNode(nil)
	tree.Parent = tree
	tree.parse(tokens)

	return tree
}

func Pretty(n *Node, depth int) {
	fmt.Printf("[%d] %q\n", depth, n.Token.Value)
	for i := range n.List {
		Pretty(n.List[i], depth+1)
	}
}
