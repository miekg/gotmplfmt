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

	lexer1 := NewLexer(string(buf))
	tokens1 := lexer1.Lex()
	pretty1 := &bytes.Buffer{}
	w := NewSuppressWriter(pretty1)
	PrettyDumb(w, tokens1)

	lexer2 := NewLexer(string(buf))
	tokens2 := lexer2.Lex()
	pretty2 := &bytes.Buffer{}
	w = NewSuppressWriter(pretty2)
	PrettyDumb(w, tokens2)

	if diff := cmp.Diff(pretty1.String(), pretty2.String()); diff != "" {
		t.Errorf("TestPretty (%s) mismatch (-want +got):\n%s", base, diff)
	}
}
