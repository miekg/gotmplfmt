package main

import (
	"bytes"
	"io"
	"strings"
)

const indent = "    "

type W struct {
	lineWrite bool
	w         io.Writer
}

func New(w io.Writer) *W { return &W{w: w} }

func (w *W) Write(data []byte) (int, error) {
	w.lineWrite = true
	if bytes.HasSuffix(data, []byte("\n")) {
		w.lineWrite = false
	}
	return w.w.Write(data)
}

// Indent writes an indent to the underlaying writer, but only if we haven't seen a something written yet.
func (w *W) Indent(level int) error {
	if w.lineWrite {
		return nil
	}
	_, err := io.WriteString(w.w, strings.Repeat(indent, level))
	return err
}
