package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

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
	blocks := Blocks(iterator.Tokens())
	//	level := 0
	for _, b := range blocks {
		// dedent only the keyword
		if b.Keyword == keyELSE || b.Keyword == keyELSEIF {
			fmt.Printf("%s", b)
			continue
		}

		if b.Keyword == keyOTHER {
			fmt.Print(b)
			continue
		}
		fmt.Printf("%s", b)
	}
}
