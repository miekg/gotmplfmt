package main

import (
	"bytes"
	"io"
	"strings"
)

var indent = "    "

type W struct {
	active bool // if true data has been written to the current line, including an indentation
	w      io.Writer
}

func New(w io.Writer) *W { return &W{w: w} }

func (w *W) Write(data []byte) (int, error) {
	w.active = true
	if bytes.HasSuffix(data, []byte("\n")) {
		w.active = false
	}
	return w.w.Write(data)
}

// Indent writes an indent to the underlaying writer, but only if we haven't seen a something written yet.
func (w *W) Indent(depth int) error {
	if w.active {
		return nil
	}
	_, err := io.WriteString(w.w, strings.Repeat(indent, depth))
	w.active = true
	return err
}
