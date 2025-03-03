package main

import (
	"fmt"
	"unicode"
)

type layout struct {
	Single bool // if true inhibit newlines
	Output bool // if true something has been written
}

// Pretty walks the tree and formats the template.
func Pretty(w *W, n *Node, depth int) {
	l := &layout{}
	l.pretty(w, n, depth)
}

// Pretty walks the tree and formats the template.
func (l *layout) pretty(w *W, n *Node, depth int) {
	// The root token, depth = 0, does not contain anything, just the beginning of tree, skip it.
	if n.Parent == nil {
		for i := range n.List {
			l.pretty(w, n.List[i], depth+1)
		}
		return
	}

	l.Render(w, n, depth, true)

	for _, n := range n.List {
		if n.Token.Type == TokenTemplate {
			if n.Token.Subtype == Else {
				l.pretty(w, n, depth)
				continue
			}
		}
		l.pretty(w, n, depth+1)
	}

	if Container(n.Token.Subtype) {
		l.Render(w, n, depth, false)
	}
}

// Render output a formatted token from the node  n.
func (l *layout) Render(w *W, n *Node, depth int, entering bool) {
	// !entering
	if !entering {
		l.Single = false // we use Println anyway here.

		if n.Token.Type == TokenHTML {
			htmltag := tag(n.Token.Value)
			if _, ok := SingleLineTag[htmltag]; ok {
				fmt.Fprintln(w)
			}
			return
		}
		w.Indent(depth - 1)
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

	// entering
	defer func() { l.Output = true }()

	// Even though we're entering, we get the end tag of an html element here as well, because we add them to the
	// AST before we close the block. So we need to indent one-less to put these on the right level.
	if n.Token.Type == TokenHTML && n.Token.Subtype == TagClose {
		w.Indent(depth - 2)
	} else {
		w.Indent(depth - 1)
	}

	if n.Token.Type == TokenHTML {
		htmltag := tag(n.Token.Value)
		if _, ok := SingleLineTag[htmltag]; ok {
			l.Single = true
		}
	}
	if l.Single {
		// range and with always start on a new line. If l.Single is true we are quaranteed that we printed some
		// html, so in that case we can just insert a newline
		if n.Token.Subtype == Range || n.Token.Subtype == With {
			fmt.Fprintln(w)
			w.Indent(depth - 1)
			l.Single = false
		}
		fmt.Fprint(w, n.Token.Value)
		return
	}

	// Empty line before these, but only when we've seen something be outputted.
	if n.Token.Type == TokenTemplate && l.Output == true {
		switch n.Token.Subtype {
		case Template, Define, Block:
			fmt.Fprintln(w)
			w.Indent(depth - 1)
		}
	}

	fmt.Fprintln(w, n.Token.Value)
}

// tag takes an HTML elements and returns the tag.
func tag(s string) (s1 string) {
	for _, r := range s {
		if r == '<' || r == '>' || r == '/' {
			continue
		}
		if unicode.IsSpace(r) {
			break
		}
		s1 += string(r)
	}
	return s1
}

// SingleLineTag holds the tags that should be rendered on a single line.
var SingleLineTag = map[string]struct{}{
	"h1":    {},
	"h2":    {},
	"h3":    {},
	"h4":    {},
	"h5":    {},
	"h6":    {},
	"title": {},

	"a":      {},
	"i":      {},
	"u":      {},
	"b":      {},
	"tt":     {},
	"em":     {},
	"strike": {},
	"strong": {},
	"mark":   {},
	"label":  {},
	"span":   {},
	"ins":    {},
	"del":    {},
	"small":  {},
	"big":    {},
	"sub":    {},
	"sup":    {},
}
