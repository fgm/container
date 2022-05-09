package stack

import (
	"sync"

	"github.com/fgm/container"
	types "github.com/fgm/container/internal"
)

type listSyncPoolStack[E any] struct {
	pool *sync.Pool
	top  *types.ListElement[E]
}

func (ss *listSyncPoolStack[E]) Push(e E) {
	le := ss.pool.Get().(*types.ListElement[E])
	le.Value = e
	le.Next = ss.top
	ss.top = le
}

func (ss *listSyncPoolStack[E]) Pop() (E, bool) {
	// Empty stack: fail.
	if ss.top == nil {
		return *new(E), false
	}

	// Non-empty stack.
	le := ss.top
	e := le.Value
	ss.top = ss.top.Next
	ss.pool.Put(le)

	return e, true
}

func NewListSyncPoolStack[E any](sizeHint int) container.Stack[E] {
	pool := &sync.Pool{New: func() any { return new(types.ListElement[E]) }}
	sl := make([]*types.ListElement[E], 0, sizeHint)
	for i := 0; i < sizeHint; i++ {
		sl = append(sl, pool.Get().(*types.ListElement[E]))
	}
	var e *types.ListElement[E]
	for i := 0; i < sizeHint; i++ {
		e, sl = sl[0], sl[1:]
		pool.Put(e)
	}
	return &listSyncPoolStack[E]{pool: pool}
}
