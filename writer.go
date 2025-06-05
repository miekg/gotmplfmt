package main

import (
	"bufio"
	"bytes"
	"io"
)

// SuppressWriter suppresses multiple newlines and only writes indents if a line is not active.
type SuppressWriter struct {
	active bool // Line is active; sutff has been written after a newline.
	indent bool // Only the indent has been written so far.
	len    int  // Character written on the current line.
	keep   bool // Dont insert line breaks, this is (only?) important for textarea tags, as ws is significant tere.
	b      *bytes.Buffer
	w      io.Writer
}

func NewSuppressWriter(w io.Writer) *SuppressWriter { return &SuppressWriter{w: w, b: &bytes.Buffer{}} }

func (s *SuppressWriter) Write(data []byte) (int, error) {
	if s.active && isTabs(data) {
		return len(data), nil
	}

	if s.keep && bytes.Equal(data, []byte("\n")) {
		return len(data), nil
	}

	// this will break with wrongly nested tags as we don't see if its open or close...
	if tag := htmlTag(string(data)); tag == "textarea" {
		s.keep = !s.keep
	}

	s.active = true
	if bytes.HasSuffix(data, []byte("\n")) {
		s.active = false
		s.len = 0
	}
	if s.active {
		s.len += len(data)
	}
	return s.b.Write(data)
}

// isTabs returns true if all bytes are tabs (and thus is an indent)
func isTabs(b []byte) bool {
	for i := range b {
		if b[i] != '\t' {
			return false
		}
	}
	return true
}

func Len(w io.Writer) int {
	s, ok := w.(*SuppressWriter)
	if !ok {
		return 0
	}
	return s.len
}

// Flushes flushes the reformatted template to w. If w is a SuppressWriter any blank lines that are only indentation are removed.
// For a define and block a newline is inserted.
func Flush(w io.Writer) {
	s, ok := w.(*SuppressWriter)
	if !ok {
		return
	}
	sc := bufio.NewScanner(s.b)
	i := 0
	for sc.Scan() {
		// if the line is indent + newline, suppress
		line := sc.Bytes()
		if isTabs(line) {
			continue
		}
		if i > 0 && extraNewline(line) {
			s.w.Write([]byte("\n"))
		}

		s.w.Write(line)
		s.w.Write([]byte("\n"))
		i++
	}
}

// Should be done in the lexer, not here...
func extraNewline(line []byte) bool {
	l1 := bytes.TrimSpace(line)

	if !bytes.HasPrefix(l1, []byte("{{")) {
		return false
	}

	if ok := bytes.HasPrefix(l1, []byte("{{define ")); ok {
		return true
	}
	if ok := bytes.HasPrefix(l1, []byte("{{- define ")); ok {
		return true
	}
	if ok := bytes.HasPrefix(l1, []byte("{{block ")); ok {
		return true
	}
	if ok := bytes.HasPrefix(l1, []byte("{{- block ")); ok {
		return true
	}
	return false
}
