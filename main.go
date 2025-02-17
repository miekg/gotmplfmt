package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	buf, _ := io.ReadAll(os.Stdin)
	lexer := NewLexer(string(buf))
	tokens := lexer.Lex()

	for _, token := range tokens {
		fmt.Printf("Type: %v, Value: %q\n", token.Type, token.Value)
	}
}
