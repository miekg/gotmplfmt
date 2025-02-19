package main

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
		// Add to current list and then proceed to parse again a level higher. But don't add the token itself.
		// n1 := &Node{Token: tokens[0], Parent: n}
		// 	n.List = append(n.List, n1)
		return n.Parent.parse(tokens[1:])
	case TagClose:
		// add close token too? Or just like end?
		return n.Parent.parse(tokens[1:])

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
