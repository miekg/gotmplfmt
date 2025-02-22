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
