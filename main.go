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

	tmpl := string(buf)
	treeSet := make(map[string]*parse.Tree)
	t := parse.New("Zgotmplfmt")
	t.Mode = parse.ParseComments | parse.SkipFuncCheck
	if _, err = t.Parse(tmpl, "{{", "}}", treeSet); err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}
	for n, ts := range treeSet {
		fmt.Printf("{{define %q }}\n", n)
		printAST(ts.Root, 0)
		fmt.Println("{{end}}")
	}
}

func printAST(node parse.Node, depth int, elseif ...bool) {
	//		fmt.Printf("**KEY %T***\n", node)
	indent := strings.Repeat("  ", depth)
	switch n := node.(type) {
	case *parse.ActionNode:
		fmt.Printf("%s", n.String())
	case *parse.TextNode:
		//format this html
		fmt.Printf("T %s", n.Text)
	case *parse.StringNode:
		fmt.Printf("%s", n.Quoted)
	case *parse.IdentifierNode:
		fmt.Printf("%s", n.Ident)
	case *parse.IfNode:
		if len(elseif) > 0 {
			fmt.Printf("if ")
		} else {
			fmt.Printf("%s{{if ", indent)
		}
		printAST(n.Pipe, depth+1)
		fmt.Println("}}")
		printAST(n.List, depth+1)
		if n.ElseList != nil {
			if _, ok := n.ElseList.Nodes[0].(*parse.IfNode); ok { // else if construct
				fmt.Printf("%s{{else ", indent)
				printAST(n.ElseList.Nodes[0], depth+1, true)
				for _, child := range n.ElseList.Nodes[1:] {
					printAST(child, depth+1)
				}
				return
			} else {
				fmt.Printf("%s{{else}}", indent)
				printAST(n.ElseList, depth+1)
			}
		}
		fmt.Printf("{{end}}\n")
	case *parse.RangeNode:
		fmt.Printf("%sRangeNode:\n", indent)
		printAST(n.Pipe, depth+1)
		printAST(n.List, depth+1)
		if n.ElseList != nil {
			fmt.Printf("%sElse:\n", indent)
			printAST(n.ElseList, depth+1)
		}
	case *parse.ListNode:
		for _, child := range n.Nodes {
			printAST(child, depth+1)
		}
	case *parse.PipeNode:
		for _, cmd := range n.Cmds {
			printAST(cmd, depth+1)
		}
	case *parse.CommandNode:
		for _, arg := range n.Args {
			printAST(arg, depth+1)
		}
	case *parse.FieldNode:
		fmt.Printf(".%s", strings.Join(n.Ident, "."))
	case *parse.VariableNode:
		fmt.Printf("%s var %s", indent, n.Ident)
	case *parse.TemplateNode:
		fmt.Printf("%s %s", indent, n.String())
	default:
		fmt.Printf("%sUnknown Node: %T\n", indent, n)
	}
}
