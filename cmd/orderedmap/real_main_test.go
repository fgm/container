package main

import (
	"bytes"
	"sort"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRealMain(t *testing.T) {
	var w1, w2 bytes.Buffer
	exitCode := realMain(&w1, &w2)

	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	const (
		expectedStandardOutput = `1 11
2 2
3 3
4 4
6 6
7 7
8 8
`

		expectedOrderedOutput = `OrderedMap:
1 11
2 2
3 3
4 4
6 6
7 7
8 8
`
	)
	// Standard map is unordered, we only check presence of keys and values.
	actualStandardMap := strings.Split(strings.Trim(w1.String(), "\n"), "\n")[1:]
	sort.Strings(actualStandardMap)
	actualStandardOutput := strings.Join(actualStandardMap, "\n") + "\n"
	if actualStandardOutput != expectedStandardOutput {
		t.Fatal(cmp.Diff(expectedStandardOutput, actualStandardOutput))
	}

	// Ordered map is ordered, we check the whole output.
	actualOrderedOutput := w2.String()
	if actualOrderedOutput != expectedOrderedOutput {
		t.Fatal(cmp.Diff(w2.String(), expectedOrderedOutput))
	}
}
