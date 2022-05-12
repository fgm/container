package main

import (
	"fmt"

	"github.com/fgm/container"
	"github.com/fgm/container/queue"
	"github.com/fgm/container/stack"
)

type Element int

// SizeHint is an indication of the maximum number of elements expected in the
// queue or stack. It is not a hard limit. Implementations may use it or not.
const sizeHint = 100

func main() {
	var e Element = 42

	q := queue.NewSliceQueue[Element](sizeHint) // resp. NewListQueue
	q.Enqueue(e)
	if lq, ok := q.(container.Countable); ok {
		fmt.Printf("elements in queue: %d\n", lq.Len())
	}
	for i := 0; i < 2; i++ {
		e, ok := q.Dequeue()
		fmt.Printf("Element: %v, ok: %t\n", e, ok)
	}

	s := stack.NewSliceStack[Element](sizeHint) // resp. NewListStack
	s.Push(e)
	if ls, ok := s.(container.Countable); ok {
		fmt.Printf("elements in s: %d\n", ls.Len())
	}
	for i := 0; i < 2; i++ {
		e, ok := s.Pop()
		fmt.Printf("Element: %v, ok: %t\n", e, ok)
	}
}
