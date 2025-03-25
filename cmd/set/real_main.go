package main

import (
	"fmt"
	"io"

	"github.com/fgm/container"
	"github.com/fgm/container/set"
)

type Element int

// SizeHint is an indication of the maximum number of elements expected in the
// set. It is not a hard limit. Implementations may use it or not.
const sizeHint = 100

func realMain(w io.Writer) int {
	var e Element = 42

	s := set.NewTrivial[Element](sizeHint)
	// Add squares.
	for i := range e {
		s.Add(i * i)
	}
	if cs, ok := s.(container.Countable); ok {
		fmt.Fprintf(w, "elements in set: %d\n", cs.Len())
	}
	// Remove elements to show that we
	// can also remove elements which are absent in the map.
	for i := Element(0); i < 10; i++ {
		del := i * i * i
		ok := s.Remove(del)
		fmt.Fprintf(w, "Element: %3v ok: %t\n", del, ok)
	}

	return 0
}
