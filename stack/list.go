package stack

import (
	"github.com/fgm/container"
	types "github.com/fgm/container/internal"
)

type listStack[E any] struct {
	top *types.ListElement[E]
}

func (ss *listStack[E]) Push(e E) {
	ss.top = &types.ListElement[E]{Value: e, Next: ss.top}
}

func (ss *listStack[E]) Pop() (E, bool) {
	// Empty stack: fail.
	if ss.top == nil {
		return *new(E), false
	}

	// Non-empty stack.
	e := ss.top.Value
	ss.top = ss.top.Next

	return e, true
}

func NewListStack[E any](_ int) container.Stack[E] {
	return &listStack[E]{}
}
