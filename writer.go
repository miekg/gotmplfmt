package main

import (
	"bytes"
	"io"
	"strings"
)

const indent = "----"

type W struct {
	lineWritten bool // if true data has been written to the current line.
	w           io.Writer
}

func New(w io.Writer) *W { return &W{w: w} }

func (w *W) Write(data []byte) (int, error) {
	w.lineWritten = true
	if bytes.HasSuffix(data, []byte("\n")) {
		w.lineWritten = false
	}
	//	bytes.Replace(data, []byte("\n"), []byte("*\n"), -1)
	return w.w.Write(data)
}

func (w *W) Text(t []byte) error {
	_, err := w.w.Write(t)
	return err
}

func (w *W) Newline() error {
	w.Write([]byte("\n"))
	return nil
}

// Indent writes an indent to the underlaying writer, but only if we haven't seen a something written yet.
func (w *W) Indent(depth int) error {
	if w.lineWritten {
		return nil
	}
	_, err := io.WriteString(w.w, strings.Repeat(indent, depth))
	return err
}
