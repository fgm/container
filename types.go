package container

import (
	"cmp"
	"iter"
)

type OrderedMap[K comparable, V any] interface {
	Delete(key K)
	Load(key K) (value V, loaded bool)
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
	WalkInOrder(cb WalkCB[E])
	WalkPostOrder(cb WalkCB[E])
	WalkPreOrder(cb WalkCB[E])
}

type WalkCB[E any] func(*E)
