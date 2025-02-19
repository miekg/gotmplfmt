package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

var (
	flagToken = flag.Bool("t", false, "Show the tokens")
	flagDebug = flag.Bool("d", false, "Show debug information")
)

func main() {
	flag.Parse()
	buf, _ := io.ReadAll(os.Stdin)
	lexer := NewLexer(string(buf))

	tokens := lexer.Lex()

	tree := Parse(tokens)
	if *flagToken {
		for _, token := range tokens {
			fmt.Printf("Type: %v, Subtype: %v, Value: %q\n", token.Type, token.Subtype, token.Value)
		}
	}
	if *flagDebug {
		indent = "   +"
	}

	w := New(os.Stdout)
	l := &Layout{}
	l.Pretty(w, tree, 0)
}

type Layout struct {
	Single bool // if true inhibit newlines
	Output bool // if true something has been written
}

func (l *Layout) Pretty(w *W, n *Node, depth int) {
	// The root token, depth = 0, does not contain anything, just the beginning of tree, skip it.
	if n.Parent == nil {
		for i := range n.List {
			l.Pretty(w, n.List[i], depth+1)
		}
		return
	}

	l.Render(w, n, depth, true)

	for _, n := range n.List {
		if n.Token.Type == TokenTemplate {
			if n.Token.Subtype == Else || n.Token.Subtype == ElseIf {
				l.Pretty(w, n, depth)
				continue
			}
		}
		l.Pretty(w, n, depth+1)
	}

	if Container(n.Token.Subtype) {
		l.Render(w, n, depth, false)
	}
}

func (l *Layout) Render(w *W, n *Node, depth int, entering bool) {
	w.Indent(depth - 1)

	if !entering { // a container type is the only one that gets false here.
		l.Single = false // we use Println anyway here
		if n.Token.Type == TokenHTML {
			fmt.Fprintln(w, EndTag(n.Token.Value))
			return
		}
		switch n.MinusEnd {
		case MinusBoth:
			fmt.Fprintln(w, "{{- end -}}")
		case MinusLeft:
			fmt.Fprintln(w, "{{- end}}")
		case MinusRight:
			fmt.Fprintln(w, "{{end -}}")
		default:
			fmt.Fprintln(w, "{{end}}")
		}
		return
	}
	defer func() {
		l.Output = true
	}()

	// entering
	if n.Token.Type == TokenHTML {
		endtag := EndTag(n.Token.Value)
		if _, ok := SingleLineTag[endtag]; ok {
			l.Single = true
		}
		// Exception alert... a <script src ==... is also a one-liner
		if strings.HasPrefix(n.Token.Value, "<script src") {
			l.Single = true
		}
	}
	if l.Single {
		fmt.Fprint(w, n.Token.Value)
		return
	}

	if n.Token.Type == TokenTemplate && l.Output == true {
		switch n.Token.Subtype {
		case Template, Define, Block:
			fmt.Fprintln(w)
			w.Indent(depth - 1)
		}
	}

	fmt.Fprintln(w, n.Token.Value)
}

// EndTag takes an HTML tag and returns it as an end tag.
func EndTag(s string) string {
	s1 := ""
	for _, r := range s {
		if r == '<' {
			continue
		}
		if r == '>' {
			continue
		}
		if r == '/' {
			continue
		}
		if unicode.IsSpace(r) {
			break
		}
		s1 += string(r)
	}

	return "</" + s1 + ">"
}

// SingleLineTag holds the tag that should be rendered on a single line. We use endtags here so we can (re)use EndTag.
var SingleLineTag = map[string]struct{}{
	"</h1>":    {},
	"</h2>":    {},
	"</h3>":    {},
	"</h4>":    {},
	"</h5>":    {},
	"</h6>":    {},
	"</title>": {},

	"</a>":      {},
	"</i>":      {},
	"</u>":      {},
	"</b>":      {},
	"</tt>":     {},
	"</em>":     {},
	"</strike>": {},
	"</strong>": {},
	"</mark>":   {},
	"</ins>":    {},
	"</del>":    {},
	"</small>":  {},
	"</big>":    {},
	"</sub>":    {},
	"</sup>":    {},
}
