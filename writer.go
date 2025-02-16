package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

const indent = "-*-*"

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

func (w *W) Text(t []byte) error {
	_, err := w.w.Write(t)
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
	io.WriteString(w.w, fmt.Sprintf("[%d] ", depth))
	_, err := io.WriteString(w.w, strings.Repeat(indent, depth))
	return err
}
