package main

import (
	"bytes"
	"io"
	"strings"

	"github.com/yosssi/gohtml"
)

const indent = "   +"

type W struct {
	lineWritten   bool // if true data has been written to the current line.
	indentWritten bool // when the lines indent has been written.
	w             io.Writer
}

func New(w io.Writer) *W { return &W{w: w} }

func (w *W) Write(data []byte) (int, error) {
	w.lineWritten = true
	if bytes.HasSuffix(data, []byte("\n")) {
		w.lineWritten = false
		w.indentWritten = false
	}
	return w.w.Write(data)
}

func (w *W) Text(t []byte, depth int) error {
	// Differ between oneliners (which may also result in that the rest of the line is one the same line) and
	// multiple lines.
	//
	// We also assume the first line is indented, and that we must indent the rest of the line.
	formatted := gohtml.FormatBytes(t)
	if len(formatted) != 0 {
		_, err := w.w.Write(indentBytes(formatted, depth))
		return err
	}

	_, err := w.w.Write(indentBytes(t, depth))
	return err
}

func (w *W) Newline() error {
	w.Write([]byte("\n"))
	w.indentWritten = false
	return nil
}

// Indent writes an indent to the underlaying writer, but only if we haven't seen a something written yet.
func (w *W) Indent(depth int) error {
	if w.lineWritten {
		return nil
	}
	if w.indentWritten {
		return nil
	}
	//	io.WriteString(w.w, fmt.Sprintf("[%d] ", depth))
	_, err := io.WriteString(w.w, strings.Repeat(indent, depth))
	return err
}

// indentBytes indents every line, except the first.
func indentBytes(b []byte, depth int) []byte {
	if bytes.Count(b, []byte("\n")) <= 1 { // single line, just return
		return b
	}

	in := bytes.Repeat([]byte(indent), depth)
	in = append([]byte("\n"), in...)
	b = bytes.Replace(b, []byte("\n"), in, -1)
	return b
}
