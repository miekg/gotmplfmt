package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestStable(t *testing.T) {
	base := "testdata/stablelarge"

	buf, err := os.ReadFile(base + ".tmpl")
	if err != nil {
		t.Errorf("Could not read %q: %s", base+".tmpl", err)
	}

	lexer := NewLexer(string(buf))
	tokens := lexer.Lex()
	tree := Parse(tokens)

	pretty1 := &bytes.Buffer{}
	w := New(pretty1)
	Pretty(w, tree, 0)

	lexer2 := NewLexer(string(buf))
	tokens2 := lexer2.Lex()
	tree2 := Parse(tokens2)
	pretty2 := &bytes.Buffer{}
	w = New(pretty2)
	Pretty(w, tree2, 0)

	if diff := cmp.Diff(pretty1.String(), pretty2.String()); diff != "" {
		t.Errorf("TestPretty (%s) mismatch (-want +got):\n%s", base, diff)
	}

}
