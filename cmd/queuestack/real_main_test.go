package main

import (
	"bytes"
	"testing"
)

func TestRealMain(t *testing.T) {
	var buf bytes.Buffer
	exitCode := realMain(&buf)

	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	expectedOutput := `elements in queue: 1
Element: 42, ok: true
Element: 0, ok: false
elements in stack: 1
Element: 42, ok: true
Element: 0, ok: false
`
	if buf.String() != expectedOutput {
		t.Fatalf("unexpected output:\n%s", buf.String())
	}
}
