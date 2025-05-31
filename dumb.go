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

func linesIndent(s string, level int) string {
	lines := strings.Split(s, "\n")
	ind := strings.Repeat(indent, level)
	for i := range lines {
		if i == 0 {
			continue // first indent already writen above
		}
		lines[i] = ind + strings.TrimSpace(lines[i])
	}
	return strings.Join(lines, "\n")
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
		// embeded text with newline, like long comments need special treatment, to get indenting of each line
		// correct. Can only be done here, because we have the indent level handy.
		if token.Type == TokenText && strings.Count(token.Value, "\n") > 2 {
			token.Value = linesIndent(token.Value, level)
		}

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
			fmt.Fprintln(w)
			printIndent(w, level)
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
