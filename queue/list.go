package queue

import (
	"github.com/fgm/container"
	types "github.com/fgm/container/internal"
)

type listQueue[E any] struct {
	head, tail *types.ListElement[E]
}

func (sq *listQueue[E]) Enqueue(e E) {
	le := &types.ListElement[E]{Value: e}
	if sq.head == nil {
		sq.head = le
	}
	if sq.tail != nil {
		sq.tail.Next = le
	}
	sq.tail = le
}

func (sq *listQueue[E]) Dequeue() (E, bool) {
	// Empty queue: fail.
	if sq.head == nil {
		return *new(E), false
	}
	// Single element queue : convert to empty.
	if sq.head.Next == nil {
		e := sq.head.Value
		sq.head = nil
		sq.tail = nil
		return e, true
	}
	// Normal queue.
	e := sq.head.Value
	sq.head = sq.head.Next
	return e, true
}

func NewListQueue[E any](_ int) container.Queue[E] {
	return &listQueue[E]{}
}
