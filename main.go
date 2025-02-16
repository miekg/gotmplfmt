package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/template/parse"
)

const (
	Failsafe = "ZgotmplfmtZ"
	End      = "{{end}}"
)

func main() {
	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	w := New(os.Stdout)

	treeSet, err := parseTree(buf)
	if err != nil {
		log.Fatalf("Failed to parse template: %s", err)
	}
	for n, ts := range treeSet {
		fmt.Fprintf(w, "{{define %q}}\n", n)
		w.Pretty(ts.Root, 0)
		fmt.Fprintf(w, "%s\n", End)
	}
}

func parseTree(buf []byte) (map[string]*parse.Tree, error) {
	treeSet := make(map[string]*parse.Tree)
	t := parse.New(Failsafe)
	t.Mode = parse.ParseComments | parse.SkipFuncCheck

	_, err := t.Parse(string(buf), "{{", "}}", treeSet)
	return treeSet, err
}

func (w *W) Pretty(node parse.Node, depth int, elseif ...bool) {
	//fmt.Fprintf(w, "**KEY %T***\n", node)
	w.Indent(depth)
	switch n := node.(type) {
	case *parse.ActionNode:
		fmt.Fprintf(w, "%s", n.String())
	case *parse.TextNode:
		// depth needed
		w.Text(n.Text)
		w.Newline()
	case *parse.StringNode:
		fmt.Fprintf(w, "%s", n.Quoted)
	case *parse.IdentifierNode:
		fmt.Fprintf(w, "%s", n.Ident)
	case *parse.IfNode:
		if len(elseif) > 0 { // we are in a {{else if ...
			fmt.Fprintf(w, "if ")
		} else {
			fmt.Fprint(w, "{{if ")
		}
		w.Pretty(n.Pipe, depth+1)

		fmt.Fprintln(w, "}}")

		w.Pretty(n.List, depth+1)
		if n.ElseList != nil {
			if _, ok := n.ElseList.Nodes[0].(*parse.IfNode); ok { // else if construct
				fmt.Fprintf(w, "%s{{else ", indent)
				w.Pretty(n.ElseList.Nodes[0], depth+1, true)
				for _, child := range n.ElseList.Nodes[1:] {
					w.Pretty(child, depth)
				}
				return
			} else {
				fmt.Fprintf(w, "%s{{else}}", indent)
				w.Pretty(n.ElseList, depth)
			}
		}
		fmt.Fprintf(w, "%s\n", End)
	case *parse.RangeNode:
		fmt.Fprintf(w, "%sRangeNode:\n", indent)
		w.Pretty(n.Pipe, depth+1)
		w.Pretty(n.List, depth+1)
		if n.ElseList != nil {
			fmt.Fprintf(w, "%sElse:\n", indent)
			w.Pretty(n.ElseList, depth+1)
		}
	case *parse.ListNode:
		for _, child := range n.Nodes {
			w.Pretty(child, depth+1)
		}
	case *parse.PipeNode:
		for _, cmd := range n.Cmds {
			w.Pretty(cmd, depth+1)
		}
	case *parse.CommandNode:
		for _, arg := range n.Args {
			w.Pretty(arg, depth+1)
		}
	case *parse.FieldNode:
		fmt.Fprintf(w, ".%s", strings.Join(n.Ident, "."))
	case *parse.VariableNode:
		fmt.Fprintf(w, "%s var %s", indent, n.Ident)
	case *parse.TemplateNode:
		fmt.Fprintf(w, "%s %s", indent, n.String())
	default:
		fmt.Fprintf(w, "%sUnknown Node: %T\n", indent, n)
	}
}
