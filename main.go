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
	if _, err = t.Parse(tmpl, "{{", "}}", treeSet); err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}
	for n, ts := range treeSet {
		println("TREE NAME", n)
		printAST(ts.Root, 0)
	}
}

func printAST(node parse.Node, depth int) {
	indent := strings.Repeat("  ", depth)
	switch n := node.(type) {
	case *parse.ActionNode:
		fmt.Printf("%s action %s\n", indent, n.String())
	case *parse.TextNode:
		fmt.Printf("%s text %s %s\n", indent, n.String(), n.Text)
	case *parse.StringNode:
		fmt.Printf("%s string %s %s\n", indent, n.String(), n.Text)
	case *parse.IdentifierNode:
		fmt.Printf("%s IdentifierNode %s %s\n", indent, n.String(), n.Ident)
	case *parse.IfNode:
		fmt.Printf("%sIfNode:%s\n", indent, n.String())
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
		fmt.Printf("%sCommandNode: %s\n", indent, n.String())
		for _, arg := range n.Args {
			printAST(arg, depth+1)
		}
	case *parse.FieldNode:
		fmt.Printf("%s field %s", indent, n.Ident)
	case *parse.VariableNode:
		fmt.Printf("%s var %s", indent, n.Ident)
	case *parse.TemplateNode:
		fmt.Printf("%s %s", indent, n.String())
	case *parse.CommentNode:
		fmt.Printf("%s %s", indent, n.String())
	default:
		fmt.Printf("%sUnknown Node: %T\n", indent, n)
	}
}
