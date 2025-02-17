package main

import (
	"io"
	"os"
)

func main() {
	buf, _ := io.ReadAll(os.Stdin)
	lexer := NewLexer(string(buf))
	Parse(lexer.Lex())
}
