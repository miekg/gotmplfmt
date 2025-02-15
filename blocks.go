package main

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/alecthomas/chroma/v2"
	"github.com/yosssi/gohtml"
)

// Block keeps together multiple tokens. This is used to create easier to handle "open" and "close" tokens.
type Block struct {
	Keyword
	Value    string
	OpenTag  string // if a keyword, holds the opening tag, may contain a -.
	CloseTag string
}

// Blocks converts a stream of tokens to blocks. I.e. an if "statement" is: {{if .Flash }}
// [CommentPreproc {{] [Keyword if] [TextWhitespace  ][NameAttribute .Flash][CommentPreproc }}]. The basic algo is to
// grab everything between commentpreprocs and slam a keyword on it. We need to track the opening and closing brace
// because they may contain a -.
func Blocks(tokens []chroma.Token) []Block {
	blocks := []Block{}
	b := Block{}
	open := false

	defer func() {
		if *flagToken {
			fmt.Println("***")
		}
	}()

	for _, t := range tokens {
		t.Value = strings.TrimSpace(t.Value)

		if *flagToken {
			fmt.Printf("[%s] %s<END>\n", t.Type.String(), t.Value)
		}

		switch t.Type {
		case chroma.CommentMultiline:
			b.Keyword = keyCOMMENT
			b.Value = t.Value
			blocks = append(blocks, b)
			b = Block{}
			continue

		case chroma.TextWhitespace:
			t.Value = " "

		case chroma.CommentPreproc:
			if !open {
				b.OpenTag = t.Value
				if strings.Contains(b.OpenTag, "-") {
					b.OpenTag += " "
				}
			} else {
				b.CloseTag = t.Value
				if strings.Contains(b.CloseTag, "-") {
					b.CloseTag = " " + b.CloseTag
				}
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

		// the lexer comes with this
		// [Other] <input type="radio"<END>
		// [CommentPreproc] {{<END>
		// Where it misses the whitespace between the two, this is annoying. We fix this by inserting a block if
		// the current value ends with a " (may need more later)
		if strings.HasSuffix(t.Value, "\"") {
			blocks = append(blocks, Block{Keyword: keyOTHER, Value: " "})
		}
		// If we see a closing tag or just a '>' we assume a keyOTHER was a oneliner and we insert a newline token
		if strings.HasSuffix(t.Value, ">") {
			// this fcks up the indenting and newline administration
			blocks = append(blocks, Block{Keyword: keyLINE, Value: "\n"})
		}

	}
	return blocks
}

func (b Block) String() string {
	if b.Keyword == keyOTHER {
		// if b.Value is a html snippet the formatting will fail (empty string is returned)
		if formatted := gohtml.Format(b.Value); formatted != "" {
			return formatted + "\n"
		}
		return b.Value
	}
	b.Value = strings.TrimLeftFunc(b.Value, unicode.IsSpace)
	return b.OpenTag + b.Value + b.CloseTag
}
