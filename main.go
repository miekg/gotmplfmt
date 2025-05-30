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
	flagDebug = flag.Bool("d", false, "Show debug information")
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
			fmt.Printf("Type: %v, Subtype: %v, Value: %q\n", token.Type, token.Subtype, token.Value)
		}
	}
	if *flagDebug {
		indent += "+"
	}

	w := NewSuppressWriter(os.Stdout)
	PrettyDumb(w, tokens)
}
