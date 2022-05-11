package queue

import (
	"github.com/fgm/container"
	types "github.com/fgm/container/internal"
)

type internalPoolQueue[E any] struct {
	maxPool    int
	pool       []types.ListElement[E]
	head, tail *types.ListElement[E]
}

func (sq *internalPoolQueue[E]) Enqueue(e E) {
	var le *types.ListElement[E]
	if len(sq.pool) != 0 {
		l := len(sq.pool) - 1
		le, sq.pool = &sq.pool[l], sq.pool[0:l]
	} else {
		le = &types.ListElement[E]{}
	}
	le.Value = e

	if sq.head == nil {
		sq.head = le
	}
	if sq.tail != nil {
		sq.tail.Next = le
	}
	sq.tail = le
}

func (sq *internalPoolQueue[E]) Dequeue() (E, bool) {
	// Empty queue: fail.
	if sq.head == nil {
		return *new(E), false
	}
	le := sq.head
	e := le.Value
	if len(sq.pool) <= sq.maxPool {
		sq.pool = append(sq.pool, *le)
	}
	// Single element queue : convert to empty.
	if sq.head.Next == nil {
		sq.head = nil
		sq.tail = nil
		return e, true
	}
	// Normal queue.
	sq.head = sq.head.Next
	return e, true
}

// NewListInternalPoolQueue preheats the element pool.
func NewListInternalPoolQueue[E any](sizeHint int) container.Queue[E] {
	q := &internalPoolQueue[E]{}
	q.maxPool = sizeHint
	q.pool = make([]types.ListElement[E], sizeHint)
	return q
}
