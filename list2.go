package container

import "sync"

type list2Queue[E any] struct {
	pool       *sync.Pool
	head, tail *listElement[E]
}

func (sq *list2Queue[E]) Enqueue(e E) {
	le := sq.pool.Get().(*listElement[E])
	le.value = e
	if sq.head == nil {
		sq.head = le
	}
	if sq.tail != nil {
		sq.tail.next = le
	}
	sq.tail = le
}

func (sq *list2Queue[E]) Dequeue() (E, bool) {
	// Empty queue: fail.
	if sq.head == nil {
		return *new(E), false
	}
	le := sq.head
	e := le.value
	sq.pool.Put(le)
	// Single element queue: convert to empty.
	if sq.head.next == nil {
		sq.head = nil
		sq.tail = nil
		return e, true
	}
	// Normal queue.
	sq.head = sq.head.next
	return e, true
}

// NewList2Queue preheats the element pool.
func NewList2Queue[E any](sizeHint int) Queue[E] {
	pool := &sync.Pool{New: func() any { return new(listElement[E]) }}
	sl := make([]*listElement[E], 0, sizeHint)
	for i := 0; i < sizeHint; i++ {
		sl = append(sl, pool.Get().(*listElement[E]))
	}
	var e *listElement[E]
	for i := 0; i < sizeHint; i++ {
		e, sl = sl[0], sl[1:]
		pool.Put(e)
	}
	return &list2Queue[E]{pool: pool}
}

type list2Stack[E any] struct {
	pool *sync.Pool
	top  *listElement[E]
}

func (ss *list2Stack[E]) Push(e E) {
	le := ss.pool.Get().(*listElement[E])
	le.value = e
	le.next = ss.top
	ss.top = le
}

func (ss *list2Stack[E]) Pop() (E, bool) {
	// Empty stack: fail.
	if ss.top == nil {
		return *new(E), false
	}

	// Non-empty stack.
	le := ss.top
	e := le.value
	ss.top = ss.top.next
	ss.pool.Put(le)

	return e, true
}

func NewList2Stack[E any](sizeHint int) Stack[E] {
	pool := &sync.Pool{New: func() any { return new(listElement[E]) }}
	sl := make([]*listElement[E], 0, sizeHint)
	for i := 0; i < sizeHint; i++ {
		sl = append(sl, pool.Get().(*listElement[E]))
	}
	var e *listElement[E]
	for i := 0; i < sizeHint; i++ {
		e, sl = sl[0], sl[1:]
		pool.Put(e)
	}
	return &list2Stack[E]{pool: pool}
}
