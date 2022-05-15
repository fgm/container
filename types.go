package container

import (
	"cmp"
	"iter"
)

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
	QueueIsEmpty              WaitableQueueState = iota // Only set at creation, not in regular use
	QueueIsBelowLowWatermark                            // Below low watermark
	QueueIsNominal                                      // Between low and high watermarks
	QueueIsAboveHighWatermark                           // Above high watermark
	QueueIsNearSaturation                               // Queue almost at full capacity: this is an emergency signal which may be used to trigger load shedding.
)

// Unit is a zero-sized struct used as a placeholder in some generic types.
type Unit = struct{}

// WaitableQueue is a concurrency-safe generic unbounded queue.
// It is meant to be used in a producer-consumer pattern,
// where the blocking behavior and capacity limits of channels are an issue.
type WaitableQueue[E any] interface {
	// Close the queue, preventing any further enqueueing, and unblocking all consumers waiting on WaitChan.
	// In most cases, it only makes sense to have it closed by the producer.
	Close()
	// Dequeue removes the first element from the queue if any is present.
	// If the queue is empty, it returns the zero value of the element type, ok is false, and the result is QueueIsBelowLowWatermark.
	// QueueIsAboveHighWatermark should be used to scale the consumer up or trigger a producer throttle.
	// QueueIsNearSaturation is the same, just more urgent, and is more useful on the Enqueue method.
	Dequeue() (e E, ok bool, result WaitableQueueState)
	// Enqueue adds an element to the queue.
	// Most applications will ignore the result of this call:
	// the most common reason to use it is checking for QueueIsNearSaturation as a trigger for producer throttling.
	// Beware of using QueueIsBelowLowWatermark as a sign to resume production,
	// as you will never get it if you do not call the method,
	// meaning your producer could never unthrottle if it completely stops emitting.
	// In most cases, use the Dequeue / WaitChan side for flow control instead.
	Enqueue(E) WaitableQueueState
	// WaitChan returns a channel that signals when an item might be available to dequeue or when the queue is closed.
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
	// Len returns the number of elements in a structure.
	// Its complexity may be higher than O(1), e.g. O(n) when it relies on Enumerable.
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

// FIXME replace by an iterator-based version like the one in Set.Items.
type Enumerable[E any] interface {
	Elements() []E
}

// BinarySearchTree is a generic binary search tree implementation with no concurrency guarantees.
// Instantiate by a zero value of the implementation.
type BinarySearchTree[E cmp.Ordered] interface {
	Clone() BinarySearchTree[E]
	Delete(*E)
	IndexOf(*E) (int, bool)
	Upsert(...*E) []*E
	WalkInOrder(cb WalkCB[E]) error
	WalkPostOrder(cb WalkCB[E]) error
	WalkPreOrder(cb WalkCB[E]) error
}

type WalkCB[E any] func(*E) error
