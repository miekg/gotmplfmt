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
		// end ends the current list addition and we fall back to the parent, add the {{end}} to the parent List
		// instead of on this level.
		n1 := &Node{Token: tokens[0], Parent: n}
		n.Parent.List = append(n.Parent.List, n1)
		return n.Parent.parse(tokens[1:])

	case Define, If, ElseIf, Else:
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

func Parse(tokens []Token) error {
	for _, token := range tokens {
		fmt.Printf("Type: %v, Subtype: %v, Value: %q\n", token.Type, token.Subtype, token.Value)
	}

	tree := NewNode(nil)
	tree.Parent = tree

	tree.parse(tokens)

	println("**AST**")
	Pretty(tree, 0)

	return nil
}

func Pretty(n *Node, depth int) {

	fmt.Printf("[%d] %s\n", depth, n.Token.Value)
	for i := range n.List {
		Pretty(n.List[i], depth+1)
	}
}

func Tree(n *Node, depth int) {
	fmt.Printf("[%d] token: %v %d nodes\n", depth, n.Token, len(n.List))
	for i := range n.List {
		Tree(n.List[i], depth+1)
	}
}
