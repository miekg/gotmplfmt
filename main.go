package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
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
	if n.Parent != nil {
		d := depth - 1
		w.Indent(d)
		// debug flag?
		//fmt.Fprintf(w, "[%d] %q\n", d, n.Token.Value)
		if n.Token.Type == TokenText && strings.Count(n.Token.Value, "\n") > 0 { // formatted multiline html
			n.Token.Value = IndentString(n.Token.Value, d)
		}

		fmt.Fprintln(w, n.Token.Value)
	}
	for i := range n.List {
		Pretty(w, n.List[i], depth+1)
	}
}
