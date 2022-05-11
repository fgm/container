package stack

import (
	"github.com/fgm/container"
	types "github.com/fgm/container/internal"
)

type internalPoolStack[E any] struct {
	maxPool int
	pool    []types.ListElement[E]
	top     *types.ListElement[E]
}

func (ss *internalPoolStack[E]) Push(e E) {
	var le *types.ListElement[E]
	if len(ss.pool) != 0 {
		l := len(ss.pool) - 1
		le, ss.pool = &ss.pool[l], ss.pool[0:l]
	} else {
		le = &types.ListElement[E]{}
	}
	le.Value = e
	le.Next = ss.top
	ss.top = le
}

func (ss *internalPoolStack[E]) Pop() (E, bool) {
	// Empty stack: fail.
	if ss.top == nil {
		return *new(E), false
	}

	// Non-empty stack.
	le := ss.top
	e := le.Value
	ss.top = ss.top.Next
	if len(ss.pool) <= ss.maxPool {
		ss.pool = append(ss.pool, *le)
	}
	return e, true
}

// NewListInternalPoolStack preheats the element pool.
func NewListInternalPoolStack[E any](sizeHint int) container.Stack[E] {
	q := &internalPoolStack[E]{}
	q.maxPool = sizeHint
	q.pool = make([]types.ListElement[E], sizeHint)
	return q
}
