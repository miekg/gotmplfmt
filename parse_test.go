package main

import "testing"

func TestParseCloseTag(t *testing.T) {
	lexer := NewLexer("</main>\n</body>\n")
	tokens := lexer.Lex()
	Parse(tokens)
}
