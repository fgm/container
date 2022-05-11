package queue

import (
	"sync"

	"github.com/fgm/container"
	types "github.com/fgm/container/internal"
)

type listSyncPoolQueue[E any] struct {
	pool       *sync.Pool
	head, tail *types.ListElement[E]
}

func (sq *listSyncPoolQueue[E]) Enqueue(e E) {
	le := sq.pool.Get().(*types.ListElement[E])
	le.Value = e
	if sq.head == nil {
		sq.head = le
	}
	if sq.tail != nil {
		sq.tail.Next = le
	}
	sq.tail = le
}

func (sq *listSyncPoolQueue[E]) Dequeue() (E, bool) {
	// Empty queue: fail.
	if sq.head == nil {
		return *new(E), false
	}
	le := sq.head
	e := le.Value
	sq.pool.Put(le)
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

// NewListSyncPoolQueue preheats the element pool.
func NewListSyncPoolQueue[E any](sizeHint int) container.Queue[E] {
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
	return &listSyncPoolQueue[E]{pool: pool}
}
