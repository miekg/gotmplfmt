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
}

const CommentPreprocClose = 1

// Blocks converts a stream of tokens to blocks. I.e. an if "statement" is: {{if .Flash }}
// [CommentPreproc {{] [Keyword if] [TextWhitespace  ][NameAttribute .Flash][CommentPreproc }}]. The basic algo is to
// grab everything between commentpreprocs and slam a keyword on it. We need to track the opening and closing brace
// because they may contain a -.
func Blocks(tokens []chroma.Token) []Block {
	blocks := []Block{}
	b := Block{}
	open := false

	for _, t := range tokens {
		t.Value = strings.TrimSpace(t.Value)

		if *flagToken {
			fmt.Printf("[%s] %s\n", t.Type.String(), t.Value)
		}

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

		case chroma.TextWhitespace:
			if open { // normalize in {{
				t.Value = " "
			}

		case chroma.Keyword:
			if open {
				switch t.Value {
				case "if":
					// this "eats" else if seen before, track that.
					if b.Keyword == keyELSE {
						b.Keyword = keyELSEIF
					} else {
						b.Keyword = keyIF
					}
				case "else":
					b.Keyword = keyELSE
				case "with":
					b.Keyword = keyWITH
				case "range":
					b.Keyword = keyRANGE
				case "end":
					b.Keyword = keyEND
				case "define":
					b.Keyword = keyDEFINE
				case "template":
					b.Keyword = keyTEMPLATE
				case "block":
					b.Keyword = keyBLOCK
				case "break":
					b.Keyword = keyBREAK
				case "continue":
					b.Keyword = keyCONTINUE
				default: // operators 'n such
					b.Keyword = keyPIPE
				}
			}
		}

		if open {
			b.Value += t.Value
			continue
		}
		if t.Type == chroma.Other {
			if t.Value == "" { // ignore the whitespace (as chroma.Other) after a preproc
				continue
			}
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

	return b.OpenTag + string(b.Keyword) + b.Value + b.CloseTag
}
