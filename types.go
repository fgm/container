package container

// Queue is generic queue with no concurrency guarantees.
// Instantiate by queue.New<implementation>Queue(sizeHint).
// The size hint MAY be used by some implementations to optimize storage.
type Queue[E any] interface {
	Enqueue(E)
	// Dequeue removes the first element from the queue. If the queue is empty,
	// it returns the zero value of the element type, and ok is false.
	Dequeue() (e E, ok bool)
}

// Stack is a generic queue with no concurrency guarantees.
// Instantiate by stack.New<implementation>Stack(sizeHint).
// The size hint MAY be used by some implementations to optimize storage.
type Stack[E any] interface {
	Push(E)
	// Pop removes the top element from the stack. If the stack is empty,
	// it returns the zero value of the element type, and ok is false.
	Pop() (e E, ok bool)
}

// Countable MAY be provided by some implementations.
type Countable interface {
	Len() int
}
