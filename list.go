package container

type listElement[E any] struct {
	value E
	next  *listElement[E]
}

type listQueue[E any] struct {
	head, tail *listElement[E]
}

func (sq *listQueue[E]) Enqueue(e E) {
	le := &listElement[E]{value: e}
	if sq.head == nil {
		sq.head = le
	}
	if sq.tail != nil {
		sq.tail.next = le
	}
	sq.tail = le
}

func (sq *listQueue[E]) Dequeue() (E, bool) {
	// Empty queue: fail.
	if sq.head == nil {
		return *new(E), false
	}
	// Single element queue: convert to empty.
	if sq.head.next == nil {
		e := sq.head.value
		sq.head = nil
		sq.tail = nil
		return e, true
	}
	// Normal queue.
	e := sq.head.value
	sq.head = sq.head.next
	return e, true
}

func NewListQueue[E any](_ int) Queue[E] {
	return &listQueue[E]{}
}

type listStack[E any] struct {
	top *listElement[E]
}

func (ss *listStack[E]) Push(e E) {
	ss.top = &listElement[E]{value: e, next: ss.top}
}

func (ss *listStack[E]) Pop() (E, bool) {
	// Empty stack: fail.
	if ss.top == nil {
		return *new(E), false
	}

	// Non-empty stack.
	e := ss.top.value
	ss.top = ss.top.next

	return e, true
}

func NewListStack[E any](_ int) Stack[E] {
	return &listStack[E]{}
}
