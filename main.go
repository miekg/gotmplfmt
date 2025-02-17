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
	Pretty(tree, 0)
}
