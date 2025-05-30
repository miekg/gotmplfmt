package main

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// TokenType represents the type of token.
type TokenType int

const (
	TokenText     TokenType = iota // Contains plain text of html.
	TokenTemplate                  // Contains template actions. The Subtype contains the actual verb.
	TokenHTML                      // Contains HTML tags. The Subtype contains OpenTag or CloseTag if this is a block level element.
)

func (t TokenType) String() string {
	switch t {
	case TokenText:
		return "text"
	case TokenHTML:
		return "html"
	case TokenTemplate:
		return "tmpl"
	}
	return ""
}

// Token represents a token in the input text.
type Token struct {
	Type    TokenType
	Subtype TokenSubtype // If the token is a TokenTemplate, this holds the keyword.
	Value   string
}

// TokenSubtype describe the deeper type of a token, like what kind of template action or if that html is an open tag or not.
type TokenSubtype int

const (
	IndentKeep        = 0
	IndentInc         = 1
	IndentDec         = -1
	IndentNewlineKeep = 2
	IndentDecKeep     = -2
)

// TokenIndent returns:
//
//   - 0: no indententation is required, keep indent level the same (0)
//   - 1: print our token, newline, then increase indent with 1
//   - 2: print a newline, then our token, then newline, then keep indent.
//   - -1: print a newline, decrease the indent with 1, print our token.
//   - -2: print a newline, decrease the indent with 1, print out token, then newline, then keep indent.
func TokenIndent(s TokenSubtype) int {
	switch s {
	case Pipe:
		return IndentKeep
	case Block:
		return IndentInc
	case Define:
		return IndentInc
	case Template:
		return IndentNewlineKeep
	case Break:
		return IndentKeep
	case Continue:
		return IndentKeep
	case Else:
		return IndentDecKeep
	case If:
		return IndentInc
	case Range:
		return IndentInc
	case With:
		return IndentInc
	case End:
		return IndentDec
	case TagOpen:
		return IndentInc
	case TagClose:
		return IndentDec
	case TagNoop:
		return IndentNewlineKeep
	case Comment:
		return IndentNewlineKeep
	}
	return 0
}

const (
	None TokenSubtype = iota
	Pipe              // Pipe is anything in {{ that is not a keyword
	Block
	Define
	Template
	Break
	Continue
	Else
	// else if caught in Else
	If
	Range
	With
	// else with caught in With
	End

	TagOpen
	TagClose
	TagNoop
	Comment
)

var Subtypes = map[string]TokenSubtype{
	"block":    Block,
	"define":   Define,
	"break":    Break,
	"continue": Continue,
	"template": Template,
	"else":     Else, // 'else with' and 'else if' fall in this category
	"end":      End,
	"if":       If,
	"range":    Range,
	"with":     With,
	"/*":       Comment,
}

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

// backup steps back one rune.
func (l *Lexer) backup() { l.pos -= l.width }

// Lex runs the lexer and returns the tokens.
func (l *Lexer) Lex() []Token {
	l.lexText()
	tokens := l.tokens
	return tokens
}

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
	// Same for the end, internal whitespace is reduced to a single space. After this we are left with {{-<space>
	// (anything else is reject by the go parser) or {{<space>thing, the later is reduced to {{thing. And again also
	// at the end.
	subtype := None
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
		for s := range Subtypes {
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

	if t == TokenText || t == TokenHTML {
		// If the token start with spaces and when trimmed is empty, we skip this token.
		if trimmed := strings.TrimLeftFunc(value, unicode.IsSpace); len(trimmed) == 0 {
			return
		}
		// If the token start with a newline we trimleft the value. We do add it to the list.
		if strings.HasPrefix(value, "\n") {
			value = strings.TrimLeftFunc(value, unicode.IsSpace)
		}
		// If the remainder contains 1 newline, we trim the whitespace at the end too.
		if strings.Count(value, "\n") == 1 {
			value = strings.TrimRightFunc(value, unicode.IsSpace)
		}
		if t == TokenText && strings.HasPrefix(value, "<!--") {
			subtype = Comment
		}
	}

	if t == TokenHTML {
		value = strings.Replace(value, "\n", "", -1) // remove all newlines inside a html tag
		// and also trim all multi spaces
		fields := strings.Fields(value)
		value = strings.Join(fields, " ")

		if !isInLineTag(value) {
			switch {
			case strings.HasPrefix(value, "</"):
				subtype = TagClose
			case strings.HasSuffix(value, "/>"):
				// none
			case strings.HasPrefix(value, "<"):
				subtype = TagOpen
			}
			if isOpenOnceTag(value) {
				subtype = TagNoop
			}
		}
	}

	l.tokens = append(l.tokens, Token{Type: t, Value: value, Subtype: subtype})
}

// lexText scans plain text until it encounters a template or html tag.
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
		if r == '<' {
			r2 := l.next()
			if r2 == '/' || unicode.IsLetter(r2) {
				l.backup()
				l.backup()
				if l.pos > l.start {
					l.emit(TokenText)
				}
				l.lexHTML()
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

// lexHTML scans an HTML tag.
func (l *Lexer) lexHTML() {
	l.next() // consume '<'
	for {
		r := l.next()
		if r == -1 {
			break
		}
		if r == '>' {
			l.emit(TokenHTML)
			break
		}
	}
}
