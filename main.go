package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

var (
	flagToken = flag.Bool("t", false, "Show the tokens")
	flagDebug = flag.Bool("d", false, "Show debug information")
)

func main() {
	flag.Parse()
	buf, _ := io.ReadAll(os.Stdin)
	lexer := NewLexer(string(buf))

	tokens := lexer.Lex()

	tree := Parse(tokens)
	if *flagToken {
		for _, token := range tokens {
			fmt.Printf("Type: %v, Subtype: %v, Value: %q\n", token.Type, token.Subtype, token.Value)
		}
	}
	if *flagDebug {
		indent = "   +"
	}

	w := New(os.Stdout)
	Pretty(w, tree, 0)
}

func Pretty(w *W, n *Node, depth int) {
	// The root token, depth = 0, does not contain anything, just the beginning of tree, skip it.
	if n.Parent == nil {
		for i := range n.List {
			Pretty(w, n.List[i], depth+1)
		}
		return
	}

	Render(w, n, depth, true)

	for _, l := range n.List {
		if l.Token.Type == TokenTemplate {
			if l.Token.Subtype == Else || l.Token.Subtype == ElseIf {
				Pretty(w, l, depth)
				continue
			}
		}
		Pretty(w, l, depth+1)
	}

	if Container(n.Token.Subtype) {
		Render(w, n, depth, false)
	}
}

func Render(w *W, n *Node, depth int, entering bool) {
	d := depth - 1
	w.Indent(d)

	if !entering { // a container type is the only one that gets false here.
		// we don't know if it was
		fmt.Fprintln(w, "{{end}}")
		return
	}

	fmt.Fprintln(w, n.Token.Value)
}
