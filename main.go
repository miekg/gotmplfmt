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
	keyLINE     Keyword = "\n"
	keyCOMMENT  Keyword = "{{/*"
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
	keyPIPE     Keyword = "n/a" // bare pipeline
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

	// For pretty printing we want to indent only when we just have written a newline.
	blocks := Blocks(iterator.Tokens())
	newline := false
	eol := ""
	level := 0
	for _, b := range blocks {
		if b.Keyword == keyEND {
			level--
		}

		if newline {
			newline = !newline
			fmt.Print(indent(level))

			// if this block is a keyword block, we close it with a newline
			if b.Keyword != keyOTHER && b.Keyword != keyPIPE {
				eol = "\n"
			}
		}

		fmt.Printf("%s%s", b, eol)

		switch b.Keyword {
		case keyIF, keyRANGE, keyWITH:
			level++
		case keyDEFINE:
			level++
			fallthrough
		case keyBLOCK:
			level++
			fallthrough
		case keyCOMMENT:
			fallthrough
		case keyTEMPLATE:
			fmt.Println()
			newline = true
		}
		if !newline {
			newline = eol != ""
		}
		eol = ""
	}
	fmt.Println() // closing newline
}

func indent(level int) string { return strings.Repeat("    ", level) }
