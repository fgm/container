package main

import (
	"fmt"
	"io"

	"github.com/fgm/container/list"
)

type Element int

func realMain(w io.Writer) int {
	var e Element = 13

	l := list.New[Element]()
	// Add squares.
	for i := range e {
		l.PushBack(i * i)
	}
	fmt.Fprintf(w, "elements in list: %d\n", l.Len())

	found := 0
	for {
		if cur := l.Front(); cur != nil {
			found++
			l.Remove(cur)
			fmt.Fprintf(w, "Element: %3v len: %d\n", cur.Value, l.Len())
		} else {
			break
		}
		if cur := l.Back(); cur != nil {
			found++
			l.Remove(cur)
			fmt.Fprintf(w, "Element: %3v len: %d\n", cur.Value, l.Len())
		} else {
			break
		}
	}
	fmt.Fprintf(w, "Found %d elements, expected %d\n", found, e)
	return 0
}
