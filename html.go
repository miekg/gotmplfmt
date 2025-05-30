package main

import "unicode"

var inlineTag = map[string]struct{}{
	"a":        {},
	"i":        {},
	"u":        {},
	"b":        {},
	"br":       {},
	"tt":       {},
	"em":       {},
	"img":      {},
	"font":     {},
	"textarea": {},
	"input":    {},
	"strike":   {},
	"strong":   {},
	"mark":     {},
	"label":    {},
	"span":     {},
	"ins":      {},
	"del":      {},
	"small":    {},
	"big":      {},
	"sub":      {},
	"sup":      {},
}

var openOnceTag = map[string]struct{}{
	"html": {},
	"body": {},
	"head": {},
	"meta": {},
	"main": {},
	"nav":  {},
}

// htmlTag takes an HTML elements and returns the tag.
func htmlTag(s string) (s1 string) {
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

func isInLineTag(s string) bool {
	tag := htmlTag(s)
	_, ok := inlineTag[tag]
	return ok
}

func isOpenOnceTag(s string) bool {
	tag := htmlTag(s)
	_, ok := openOnceTag[tag]
	return ok
}
