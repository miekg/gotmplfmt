package main

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

// RenderHTML recursively prints the HTML node structure with indentation.
func RenderHTML(w io.Writer, n *html.Node, depth int) {
	indent := strings.Repeat(" ", depth)

	// Print the node type and data
	switch n.Type {
	case html.ElementNode:
		fmt.Printf("%s<%s>\n", indent, n.Data)
	case html.TextNode:
		text := strings.TrimSpace(n.Data)
		if text != "" {
			fmt.Printf("%s%s\n", indent, text)
		}
	case html.CommentNode:
		fmt.Printf("%s<!-- %s -->\n", indent, n.Data)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		RenderHTML(w, c, depth+1)
	}

	if n.Type == html.ElementNode {
		fmt.Printf("%s</%s>\n", indent, n.Data)
	}
}
