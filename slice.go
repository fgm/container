package container

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

func NewSliceQueue[E any](sizeHint int) Queue[E] {
	q := &sliceQueue[E]{}
	q.store = make([]E, 0, sizeHint)
	return q
}

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

func NewSliceStack[E any](sizeHint int) Stack[E] {
	q := &sliceStack[E]{}
	q.store = make([]E, 0, sizeHint)
	return q
}
