package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	flagToken = flag.Bool("t", false, "Show the tokens")
	flagWidth = flag.Int("w", 120, "Maximum line width")
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

	if *flagToken {
		for _, token := range tokens {
			fmt.Printf("Type: %s, Subtype: %2d, Value: %q\n", token.Type, token.Subtype, token.Value)
		}
	}

	w := NewSuppressWriter(os.Stdout)
	PrettyDumb(w, tokens)
}
