package main

import (
	"fmt"

	"github.com/fgm/container"
)

type Element int

// SizeHint is an indication of the maximum number of elements expected in the
// queue or stack. It is not a hard limit. Implementations may use it or not.
const sizeHint = 100

func main() {
	var e Element = 42

	queue := container.NewSliceQueue[Element](sizeHint) // resp. NewListQueue
	queue.Enqueue(e)
	if lq, ok := queue.(container.Countable); ok {
		fmt.Printf("elements in queue: %d\n", lq.Len())
	}
	for i := 0; i < 2; i++ {
		e, ok := queue.Dequeue()
		fmt.Printf("Element: %v, ok: %t\n", e, ok)
	}

	stack := container.NewSliceStack[Element](sizeHint) // resp. NewListStack
	stack.Push(e)
	if ls, ok := stack.(container.Countable); ok {
		fmt.Printf("elements in stack: %d\n", ls.Len())
	}
	for i := 0; i < 2; i++ {
		e, ok := stack.Pop()
		fmt.Printf("Element: %v, ok: %t\n", e, ok)
	}
}
