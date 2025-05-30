package main

import (
	"fmt"
	"io"
	"strings"
)

const indent = "\t"

func printIndent(w io.Writer, level int) {
	if level < 0 {
		level = 0
	}
	io.WriteString(w, strings.Repeat(indent, level))
}

// PrettyDumb does not use a tree, just the list of tokens as parsed and indents them as appropriate.
func PrettyDumb(w io.Writer, tokens []Token) {
	// We sometimes write too many newlines, we fix this in "post" with the Flush function.
	level := 0
	for _, token := range tokens {
		if level < 0 {
			level = 0
		}
		printIndent(w, level)
		ti := TokenIndent(token.Subtype)
		switch ti {
		case IndentDecKeep:
			fmt.Fprintln(w)
			printIndent(w, level-1)
			fmt.Fprintf(w, "%s\n", token.Value)
		case IndentDec:
			fmt.Fprintln(w)
			printIndent(w, level-1)
			fmt.Fprintf(w, "%s\n", token.Value)
			level -= 1
		case IndentKeep:
			fmt.Fprintf(w, "%s", token.Value)
		case IndentInc:
			fmt.Fprintf(w, "%s\n", token.Value)
			level += 1
		case IndentNewlineKeep:
			fmt.Fprintln(w)
			printIndent(w, level)
			fmt.Fprintf(w, "%s\n", token.Value)
		}
		if Len(w) > *flagWidth {
			fmt.Fprintln(w)
			printIndent(w, level)
		}
	}
	Flush(w)
}
