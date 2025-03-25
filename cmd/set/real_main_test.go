package main

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRealMain(t *testing.T) {
	var buf bytes.Buffer
	exitCode := realMain(&buf)

	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	expectedOutput := `elements in set: 42
Element:   0 ok: true
Element:   1 ok: true
Element:   8 ok: false
Element:  27 ok: false
Element:  64 ok: true
Element: 125 ok: false
Element: 216 ok: false
Element: 343 ok: false
Element: 512 ok: false
Element: 729 ok: true
`
	if buf.String() != expectedOutput {
		t.Fatal(cmp.Diff(expectedOutput, buf.String()))
	}
}
