package container

import "iter"

// OrderedMap has the same API as a sync.Map for the specific case of OrderedMap[any, any].
type OrderedMap[K comparable, V any] interface {
	Delete(key K)
	Load(key K) (value V, loaded bool)
	// Range is similar to the sync.Map Range method but can fail if the callback deletes map entries.
	Range(func(key K, value V) bool)
	Store(key K, value V)
}

// Queue is a generic queue with no concurrency guarantees.
// Instantiate by queue.New<implementation>Queue(sizeHint).
// The size hint MAY be used by some implementations to optimize storage.
type Queue[E any] interface {
	Enqueue(E)
	// Dequeue removes the first element from the queue. If the queue is empty,
	// it returns the zero value of the element type, and ok is false.
	Dequeue() (e E, ok bool)
}

//go:generate go run golang.org/x/tools/cmd/stringer -type=WaitableQueueState -output=types_waitablequeuestate_string.go
type WaitableQueueState int

const (
	_                         WaitableQueueState = iota
	QueueIsBelowLowWatermark                     // Below low watermark
	QueueIsNominal                               // Between low and high watermarks
	QueueIsAboveHighWatermark                    // Above high watermark
)

// Unit is a zero-sized struct used as a placeholder in some generic types.
type Unit = struct{}

// WaitableQueue is a concurrency-safe generic unbounded queue.
// It is meant to be used in a producer-consumer pattern,
// where the blocking behavior and capacity limits of channels are an issue.
type WaitableQueue[E any] interface {
	// Dequeue removes the first element from the queue if any is present.
	// If the queue is empty, it returns the zero value of the element type, ok is false, and the result is QueueIsBelowLowWatermark.
	Dequeue() (e E, ok bool, result WaitableQueueState)
	// Enqueue adds an element to the queue.
	Enqueue(E) WaitableQueueState
	// WaitChan returns a channel that signals when an item might be available or when the queue is closed.
	WaitChan() <-chan Unit
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
// For concurrency-safe types, it is not atomic vs other operations,
// meaning it MUST NOT be used to take decisions, but only as an observability/debugging tool.
type Countable interface {
	Len() int
}

type Set[E comparable] interface {
	Add(item E) (found bool)
	Remove(item E) (found bool)
	Contains(item E) bool
	Clear() (count int)
	Items() iter.Seq[E]

	Union(other Set[E]) Set[E]
	Intersection(other Set[E]) Set[E]
	Difference(other Set[E]) Set[E]
}
