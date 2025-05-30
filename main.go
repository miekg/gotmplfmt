package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var (
	flagToken = flag.Bool("t", false, "Show the tokens")
	flagDebug = flag.Bool("d", false, "Show debug information")
	flagDumb  = flag.Bool("b", false, "Try dumb pretty printing")
)

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("gotmplfmt: %s", err)
		}
		Reformat(data)
		return
	}

	for _, a := range flag.Args() {
		data, err := os.ReadFile(a)
		if err != nil {
			log.Fatalf("gotmplfmt: %s", err)
		}
		Reformat(data)
	}
}

func Reformat(data []byte) {
	lexer := NewLexer(string(data))
	tokens := lexer.Lex()

	tree := Parse(tokens)
	if *flagToken {
		for _, token := range tokens {
			fmt.Printf("Type: %v, Subtype: %v, Value: %q\n", token.Type, token.Subtype, token.Value)
		}
	}
	if *flagDumb {
		w := NewSuppressWriter(os.Stdout)
		PrettyDumb(w, tokens)
		return
	}

	if *flagDebug {
		indent += "+"
	}

	w := New(os.Stdout)
	Pretty(w, tree, 0)

	if *flagDebug {
		structure(tree, 0)
	}
}

func structure(n *Node, depth int) {
	if n.Parent == nil {
		for i := range n.List {
			structure(n.List[i], depth+1)
		}
		return
	}

	fmt.Printf("[%s] List %d token: %q \t\t(parent: %q)\n", strings.Repeat(" ", depth), len(n.List), n.Token.Value, n.Parent.Token.Value)

	for _, n := range n.List {
		structure(n, depth+1)
	}
}
