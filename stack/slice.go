package stack

import "github.com/fgm/container"

type sliceStack[E any] struct {
	store []E
}

func (ss *sliceStack[E]) Push(e E) {
	ss.store = append(ss.store, e)
}

func (ss *sliceStack[E]) Pop() (E, bool) {
	l := len(ss.store) - 1
	if l == -1 {
		return *new(E), false
	}
	e := ss.store[l]
	ss.store = ss.store[:l]
	return e, true
}

func (ss *sliceStack[E]) Len() int {
	return len(ss.store)
}

func NewSliceStack[E any](sizeHint int) container.Stack[E] {
	q := &sliceStack[E]{}
	q.store = make([]E, 0, sizeHint)
	return q
}
