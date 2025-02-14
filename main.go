package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/template/parse"
)

func main() {
	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	// Example template
	tmpl := string(buf)

	treeSet := make(map[string]*parse.Tree)
	t := parse.New("example")
	t.Mode = parse.ParseComments | parse.SkipFuncCheck
	tree, err := t.Parse(tmpl, "{{", "}}", treeSet)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	// Access the root of the AST
	root := tree.Root

	// Print the AST
	printAST(root, 0)
}

func printAST(node parse.Node, depth int) {
	indent := strings.Repeat("  ", depth)
	switch n := node.(type) {
	case *parse.ActionNode:
		fmt.Printf("%sActionNode: %s\n", indent, n.Pipe)
	case *parse.TextNode:
		fmt.Printf("%sTextNode: %q\n", indent, n.Text)
	case *parse.IfNode:
		fmt.Printf("%sIfNode:\n", indent)
		printAST(n.Pipe, depth+1)
		printAST(n.List, depth+1)
		if n.ElseList != nil {
			fmt.Printf("%sElse:\n", indent)
			printAST(n.ElseList, depth+1)
		}
	case *parse.RangeNode:
		fmt.Printf("%sRangeNode:\n", indent)
		printAST(n.Pipe, depth+1)
		printAST(n.List, depth+1)
		if n.ElseList != nil {
			fmt.Printf("%sElse:\n", indent)
			printAST(n.ElseList, depth+1)
		}
	case *parse.ListNode:
		fmt.Printf("%sListNode:\n", indent)
		for _, child := range n.Nodes {
			printAST(child, depth+1)
		}
	case *parse.PipeNode:
		fmt.Printf("%sPipeNode:\n", indent)
		for _, cmd := range n.Cmds {
			printAST(cmd, depth+1)
		}
	case *parse.CommandNode:
		fmt.Printf("%sCommandNode:\n", indent)
		for _, arg := range n.Args {
			printAST(arg, depth+1)
		}
	case *parse.FieldNode:
		fmt.Printf("%sFieldNode: %s\n", indent, n.Ident)
	case *parse.VariableNode:
		fmt.Printf("%sVariableNode: %s\n", indent, n.Ident)
	case *parse.TemplateNode:
		fmt.Printf("%sTemplateNode: %s\n", indent, n.String())
	default:
		fmt.Printf("%sUnknown Node: %T\n", indent, n)
	}
}
