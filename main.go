package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/alecthomas/chroma/v2/lexers"
)

type Keyword string

const (
	keyOTHER    Keyword = "other"
	keyEND      Keyword = "end"
	keyIF       Keyword = "if"
	keyRANGE    Keyword = "range"
	keyWITH     Keyword = "with"
	keyELSE     Keyword = "else"
	keyELSEIF   Keyword = "else if"
	keyDEFINE   Keyword = "define"
	keyBLOCK    Keyword = "block"
	keyTEMPLATE Keyword = "template"
	keyBREAK    Keyword = "break"
	keyCONTINUE Keyword = "continue"
)

var (
	flagToken = flag.Bool("t", false, "print tokens")
)

func main() {
	flag.Parse()
	lexer := lexers.Get("go-template")
	if lexer == nil {
		log.Fatal("No lexer seen")
	}
	contents, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	iterator, err := lexer.Tokenise(nil, string(contents))
	if err != nil {
		log.Fatal(err)
	}
	blocks := Blocks(iterator.Tokens())
	level := 0
	for _, b := range blocks {
		// if b.Step is postive, it applies after the element being printed
		if b.Step > 0 {
			fmt.Printf("%s%s\n", indent(level), b)
			level += b.Step
			continue
		}

		// dedent only the keyword (for else, and ...)
		if b.Keyword == keyELSE || b.Keyword == keyELSEIF {
			fmt.Printf("%s%s\n", indent(level-1), b)
			continue
		}

		level += b.Step
		if level < 0 { // {{end}} are also used when we haven't raised the indentlevel
			level = 0
		}
		fmt.Printf("%s%s\n", indent(level), b)
	}
}

func indent(level int) string {
	if level < 0 {
		level = 0
	}
	return strings.Repeat("\t", level)
}
