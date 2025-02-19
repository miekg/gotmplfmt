package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPretty(t *testing.T) {
	dir := "testdata"
	testFiles, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("Could not read %q: %s", dir, err)
	}
	for _, f := range testFiles {
		if f.IsDir() {
			continue
		}
		if filepath.Ext(f.Name()) != ".tmpl" {
			continue
		}
		buf, err := os.ReadFile(dir + "/" + f.Name())
		if err != nil {
			t.Errorf("Could not read %q: %s", dir+"/"+f.Name(), err)
		}

		prettyfile := dir + "/" + f.Name()[:len(f.Name())-4] + "pretty"
		prettybuf, err := os.ReadFile(prettyfile)
		if err != nil {
			continue
		}

		lexer := NewLexer(string(buf))
		tokens := lexer.Lex()
		tree := Parse(tokens)

		b := &bytes.Buffer{}
		w := New(b)
		Pretty(w, tree, 0)

		if diff := cmp.Diff(string(prettybuf), b.String()); diff != "" {
			t.Errorf("TestPretty (%s) mismatch (-want +got):\n%s", f.Name(), diff)
		}
	}

}
