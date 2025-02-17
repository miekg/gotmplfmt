package main

import (
	"fmt"
)

// Node is the parse tree of a template.
type Node struct {
	Tokens []Token
	List   []*Node
	Parent *Node
}

// Parse parses the tokens and adds them the list in n. It returns the token that are not yet consumed.
func (n *Node) parse(tokens []Token) []Token {
	if len(tokens) == 0 {
		return nil
	}
	if tokens[0].Type == TokenText {
		// add to the current node and continue with the rest
		n.Tokens = append(n.Tokens, tokens[0])
		return n.parse(tokens[1:])
	}

	// TokenTemplate, add them all as a list
	switch tokens[0].Subtype {
	case End:
		// end ends the current list addition and we fall back to the parent
		n.Parent.Tokens = append(n.Parent.Tokens, tokens[0])
		return n.Parent.parse(tokens[1:])

	case Define:
		n.Tokens = append(n.Tokens, tokens[0])
		n1 := &Node{Parent: n}
		n.List = append(n.List, n1)
		return n1.parse(tokens[1:])

	default:
		// add to the current open node
		n.Tokens = append(n.Tokens, tokens[0])
		return n.parse(tokens[1:])

	}
}

func Parse(tokens []Token) error {
	for _, token := range tokens {
		fmt.Printf("Type: %v, Subtype: %v, Value: %q\n", token.Type, token.Subtype, token.Value)
	}

	tree := &Node{}
	tree.Parent = tree

	tree.parse(tokens)

	println("**AST**")
	Tree(tree, 0)
	println("**AST**")
	Pretty(tree, 0)

	return nil
}

func Pretty(n *Node, depth int) {

	for _, t := range n.Tokens {
		fmt.Printf("[%d] %s\n", depth, t.Value)
	}
	for i := range n.List {
		Pretty(n.List[i], depth+1)
	}
}

func Tree(n *Node, depth int) {
	fmt.Printf("[%d] %d tokens %d nodes\n", depth, len(n.Tokens), len(n.List))
	for i := range n.List {
		Tree(n.List[i], depth+1)
	}

}
