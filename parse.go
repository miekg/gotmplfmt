package main

import (
	"fmt"
	"strings"
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
	tree.parse(tokens)

	return tree
}

func Pretty(w *W, n *Node, depth int) {
	// The root token, depth = 0, does not contain anything, just the beginning of tree, skip it.
	if n.Parent != nil {
		d := depth - 1
		if n.Token.Subtype == End { // pull in {{end}}s
			d--
		}

		w.Indent(d)
		// debug flag?
		//fmt.Fprintf(w, "[%d] %q\n", d, n.Token.Value)
		if n.Token.Type == TokenText && strings.Count(n.Token.Value, "\n") > 0 { // formatted multieline html
			n.Token.Value = IndentString(n.Token.Value, d)
		}

		fmt.Fprintln(w, n.Token.Value)
	}
	for i := range n.List {
		Pretty(w, n.List[i], depth+1)
	}
}
