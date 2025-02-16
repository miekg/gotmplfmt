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
	for n, ts := range treeSet {
		println("TREE NAME", n)
		printAST(ts.Root, 0)
	}
	println("PARSED TREE")

	printAST(tree.Root, 0)
}

func printAST(node parse.Node, depth int) {
	indent := strings.Repeat("  ", depth)
	fmt.Printf("%s%s", indent, node.String())
}
