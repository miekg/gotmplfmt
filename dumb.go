package main

import (
	"fmt"
	"io"
	"strings"
)

var indent = "\t"

func printIndent(w io.Writer, level int) {
	if level < 0 {
		level = 0
	}
	io.WriteString(w, strings.Repeat(indent, level))
}

// PrettyDumb does not use a tree, just the list of tokens as parsed and indents them as appropiate.
func PrettyDumb(w io.Writer, tokens []Token) {
	// We sometimes write too many newline, we fix this in "post" by collapsing \n\n into \n.
	level := 0
	for _, token := range tokens {
		printIndent(w, level)
		ti := TokenIndent(token.Subtype)
		switch ti {
		case -2:
			fmt.Fprintln(w)
			printIndent(w, level-1)
			fmt.Fprintf(w, "%s\n", token.Value)
		case -1:
			fmt.Fprintln(w)
			printIndent(w, level-1)
			fmt.Fprintf(w, "%s\n", token.Value)
			level -= 1
		case 0:
			fmt.Fprintf(w, "%s", token.Value)
		case 1:
			fmt.Fprintf(w, "%s\n", token.Value)
			level += 1
		case 2:
			fmt.Fprintln(w)
			printIndent(w, level)
			fmt.Fprintf(w, "%s\n", token.Value)
		}
	}
	Flush(w)
}
