package queue

import "github.com/fgm/container"

type sliceQueue[E any] struct {
	store []E
}

func (sq *sliceQueue[E]) Enqueue(e E) {
	sq.store = append(sq.store, e)
}

func (sq *sliceQueue[E]) Dequeue() (E, bool) {
	if len(sq.store) == 0 {
		return *new(E), false
	}
	e := sq.store[0]
	sq.store = sq.store[1:]
	return e, true
}

func (sq *sliceQueue[E]) Len() int {
	return len(sq.store)
}

func NewSliceQueue[E any](sizeHint int) container.Queue[E] {
	q := &sliceQueue[E]{}
	q.store = make([]E, 0, sizeHint)
	return q
}
