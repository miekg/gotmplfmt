package main

import (
	"io"
	"os"
)

func main() {
	buf, _ := io.ReadAll(os.Stdin)
	lexer := NewLexer(string(buf))

	tree := Parse(lexer.Lex())

	println("**AST**")
	w := New(os.Stdout)
	Pretty(w, tree, 0)
}
