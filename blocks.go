package main

import (
	"fmt"
	"strings"

	"github.com/alecthomas/chroma/v2"
)

// Block keeps together multiple tokens. This is used to create easier to handle "open" and "close" tokens.
type Block struct {
	Keyword
	Value    string
	OpenTag  string // if a keyword, holds the opening tag, may contain a -.
	CloseTag string
	Step     int // How to indent: -1, indent less, 1 indent more.
}

// Blocks converts a stream of tokens to blocks. I.e. an if "statement" is: {{if .Flash }}
// [CommentPreproc {{] [Keyword if] [TextWhitespace  ][NameAttribute .Flash][CommentPreproc }}]. The basic algo is to
// grab everything between commentpreprocs and slam a keyword on it. We need to track the opening and closing brace
// because they may contain a -.
func Blocks(tokens []chroma.Token) []Block {
	blocks := []Block{}
	b := Block{}
	open := false
	for _, t := range tokens {
		if *flagToken {
			fmt.Printf("[%s] %s\n", t.Type.String(), t.Value)
		}
		if strings.Count(t.Value, "\n") > 1 && strings.TrimSpace(t.Value) == "" { // closing newline and empty line
			blocks = append(blocks, Block{Keyword: keyOTHER})
			continue
		}
		t.Value = strings.TrimSpace(t.Value)

		switch t.Type {
		case chroma.CommentPreproc:
			if !open {
				b.OpenTag = t.Value
			} else {
				b.CloseTag = t.Value
				blocks = append(blocks, b)
				b = Block{}
			}

			open = !open
			continue

		case chroma.Keyword:
			if open {
				switch t.Value {
				case "if":
					// this "eats" else if seen before, track that.
					if b.Keyword == keyELSE {
						b.Keyword = keyELSEIF
					} else {
						b.Keyword = keyIF
						b.Step = 1
					}
				case "else":
					b.Keyword = keyELSE
				case "with":
					b.Keyword = keyWITH
					b.Step = 1
				case "range":
					b.Keyword = keyRANGE
					b.Step = 1
				case "end":
					b.Keyword = keyEND
					b.Step = -1
				case "define":
					b.Keyword = keyDEFINE
					b.Step = 1
				case "template":
					b.Keyword = keyTEMPLATE
				case "block":
					b.Keyword = keyBLOCK
					b.Step = 1
				case "break":
					b.Keyword = keyBREAK
				case "continue":
					b.Keyword = keyCONTINUE
				}
			}
			continue
		}

		if t.Value == "" {
			continue
		}

		if open {
			if b.Value == "" {
				b.Value = t.Value
			} else {
				b.Value += " " + t.Value
			}
			continue
		}

		// not open, add the token at hand to the list as a block.
		blocks = append(blocks, Block{Keyword: keyOTHER, Value: t.Value})
	}
	return blocks
}

func (b Block) String() string {
	if b.Keyword == keyOTHER {
		return b.Value
	}

	keyword := string(b.Keyword)
	if keyword != "" { // it can be "" is there is just a pipeline without any keywords
		keyword = " " + keyword + " "
	} else {
		keyword = " "
	}

	if b.Value != "" {
		return b.OpenTag + keyword + b.Value + " " + b.CloseTag
	}
	return b.OpenTag + keyword + b.CloseTag
}
