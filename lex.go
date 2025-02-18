package main

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/yosssi/gohtml"
)

// TokenType represents the type of token.
type TokenType int

const (
	TokenText TokenType = iota
	TokenTemplate
)

// Token represents a token in the input text.
type Token struct {
	Type    TokenType
	Subtype TokenSubtype // If the token is a TokenTemplate, this holds the keyword.
	Value   string
}

type TokenSubtype int

const (
	Pipe TokenSubtype = iota // Pipe is anything in {{ that is not a keyword
	Block
	Define
	Template
	Break
	Continue
	ElseIf
	Else
	If
	Range
	With
	End
)

// Container returns true if the TokenSubType is a container type
func Container(s TokenSubtype) bool {
	switch s {
	case Block:
		fallthrough
	case Define:
		fallthrough
	case ElseIf:
		fallthrough // elseif really container type, we need only 1 end
	case Else:
		fallthrough
	case If:
		fallthrough
	case Range:
		fallthrough
	case With:
		return true

	}

	return false
}

var Subtypes = map[string]TokenSubtype{
	"block":    Block,
	"define":   Define,
	"break":    Break,
	"continue": Continue,
	"template": Template,
	"else if":  ElseIf, // detect as seperate substype
	"else":     Else,
	"end":      End,
	"if":       If,
	"range":    Range,
	"with":     With,
}

// odererd list of keyword, so we are sure 'else if' becomes before 'else'
var sublist = []string{
	"block",
	"template",
	"define",
	"break",
	"continue",
	"else if",
	"else",
	"end",
	"if",
	"range",
	"with",
}

// order list of subtypes

// Lexer holds the state of the lexer.
type Lexer struct {
	input  string
	start  int
	pos    int
	width  int
	tokens []Token
}

// NewLexer creates a new lexer for the given input.
func NewLexer(input string) *Lexer { return &Lexer{input: input} }

// Lex runs the lexer and returns the tokens.
func (l *Lexer) Lex() []Token { l.lexText(); return l.tokens }

// backup steps back one rune.
func (l *Lexer) backup() { l.pos -= l.width }

// next returns the next rune in the input.
func (l *Lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return -1
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	return r
}

// emit adds a token to the token list.
func (l *Lexer) emit(t TokenType) {
	value := l.input[l.start:l.pos]

	// some cleanup TokenTemplates, {{<space>thing is reduced to {{thing, execpt for {{-, then it is {{- thing.
	// some for the end, internal whitespace is reduced to a single space. After this we are left with {{-<space>
	// (anything else is reject by the go parser) or {{<space>thing, the later is reduced to {{thing. And again also
	// at the end.
	subtype := Pipe
	if t == TokenTemplate {
		value = strings.Join(strings.Fields(value), " ")

		if strings.HasPrefix(value, "{{ ") {
			value = "{{" + strings.TrimPrefix(value, "{{ ")
		}
		if strings.HasSuffix(value, " }}") {
			value = strings.TrimSuffix(value, " }}") + "}}"
		}
		// beginning is now {{ or {{-, check if we can extract the subtype
	Loop:
		for _, s := range sublist {
			switch {
			case strings.HasPrefix(value, "{{"+s):
				subtype = Subtypes[s]
				break Loop

			case strings.HasPrefix(value, "{{- "+s):
				subtype = Subtypes[s]
				break Loop
			}
		}
	}

	defer func() { l.start = l.pos }()

	if t == TokenTemplate && subtype == End { // Skip ends in the AST
		return
	}

	if t == TokenText {
		// If the token start with spaces and when trimmed is empty, we skip this token.
		if trimmed := strings.TrimLeftFunc(value, unicode.IsSpace); len(trimmed) == 0 {
			return
		}
		// If the token start with a newline we trimleft the value. We do add it to the list.
		if strings.HasPrefix(value, "\n") {
			value = strings.TrimLeftFunc(value, unicode.IsSpace)
		}

		// If the remainder contains 1 newline, we trim the whitespace at the end too
		if strings.Count(value, "\n") == 1 {
			value = strings.TrimRightFunc(value, unicode.IsSpace)
		}
		// Try to fmt the HTML, if fails use the original value
		formatted := gohtml.Format(value)
		if formatted != "" {
			value = formatted
		}
	}

	l.tokens = append(l.tokens, Token{Type: t, Value: value, Subtype: subtype})
}

// lexText scans plain text until it encounters a template tag.
func (l *Lexer) lexText() {
	for {
		r := l.next()
		if r == -1 {
			break
		}
		if r == '{' {
			r2 := l.next()
			if r2 == '{' {
				l.backup()
				l.backup()
				if l.pos > l.start {
					l.emit(TokenText)
				}
				l.lexTemplate()
				continue
			}
			l.backup()
		}
	}
	if l.pos > l.start {
		l.emit(TokenText)
	}
}

// lexTemplate scans a template tag.
func (l *Lexer) lexTemplate() {
	l.next() // consume '{'
	l.next() // consume '{'
	for {
		r := l.next()
		if r == -1 {
			break
		}
		if r == '}' {
			r2 := l.next()
			if r2 == '}' {
				l.emit(TokenTemplate)
				break
			}
			l.backup()
		}
	}
}
