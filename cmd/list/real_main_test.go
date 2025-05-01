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

	expectedOutput := `elements in list: 13
Element:   0 len: 12
Element: 144 len: 11
Element:   1 len: 10
Element: 121 len: 9
Element:   4 len: 8
Element: 100 len: 7
Element:   9 len: 6
Element:  81 len: 5
Element:  16 len: 4
Element:  64 len: 3
Element:  25 len: 2
Element:  49 len: 1
Element:  36 len: 0
Found 13 elements, expected 13
`
	if buf.String() != expectedOutput {
		t.Fatal(cmp.Diff(expectedOutput, buf.String()))
	}
}
