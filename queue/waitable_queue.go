package queue

import (
	"fmt"
	"sync"

	"github.com/fgm/container"
)

type unit = container.Unit

// waitable implements WaitableQueue
type waitable[E any] struct {
	closed bool
	items  []E
	hi, lo int // Low and high watermarks
	mu     sync.Mutex
	signal chan unit // Used to signal availability or closure
}

// NewWaitableQueue creates a new WaitableQueue with the given initial capacity and watermarks.
//
// The three arguments are in number of elements, not in bytes.
// Implementations MAY use the initial capacity to preallocate storage.
func NewWaitableQueue[E any](initialCapacity int, lowWatermark, highWatermark int) (container.WaitableQueue[E], error) {
	if initialCapacity < 0 {
		return nil, fmt.Errorf("container: initialCapacity (%d) cannot be negative", initialCapacity)
	}
	if lowWatermark < 0 {
		return nil, fmt.Errorf("container: lowWatermark (%d) cannot be negative", lowWatermark)
	}
	if highWatermark < 0 {
		return nil, fmt.Errorf("container: highWatermark (%d) cannot be negative", highWatermark)
	}
	if lowWatermark > highWatermark {
		return nil, fmt.Errorf("container: lowWatermark (%d) cannot be greater than highWatermark (%d)", lowWatermark, highWatermark)
	}

	// Use a buffered channel of size 1. This prevents Enqueue
	// from blocking if the Dequeue side isn't waiting *at the exact moment*.
	// It acts like a latch: if signal is sent and no one is waiting,
	// the next wait will immediately succeed.
	return &waitable[E]{
		closed: false,
		items:  make([]E, 0, initialCapacity),
		hi:     highWatermark,
		lo:     lowWatermark,
		signal: make(chan unit, 1),
	}, nil
}

// getState returns the current state of the queue.
//
// It MUST only be called while holding the mutex to avoid race conditions.
func (bq *waitable[E]) getState() container.WaitableQueueState {
	l := len(bq.items) // Do not use bq.Len() here, it would deadlock.
	switch {
	case l < bq.lo:
		return container.QueueIsBelowLowWatermark
	case l > bq.hi:
		return container.QueueIsAboveHighWatermark
	default:
		return container.QueueIsNominal
	}
}

// Enqueue adds an item and signals *if* necessary.
func (bq *waitable[E]) Enqueue(item E) container.WaitableQueueState {
	bq.mu.Lock()
	defer bq.mu.Unlock()

	if bq.closed {
		panic("enqueue on closed queue") // Or return an error
	}

	bq.items = append(bq.items, item)

	// Signal that an item is available.
	// Use a non-blocking send because the channel has buffer size 1.
	// If the buffer is full, it means a signal is already pending,
	// and we don't need to send another one.
	select {
	case bq.signal <- unit{}:
		// Signal sent
	default:
		// Signal already pending or channel closed
	}
	// Return the current state of the queue
	return bq.getState()
}

// Dequeue removes and returns an item if available.
func (bq *waitable[E]) Dequeue() (E, bool, container.WaitableQueueState) {
	bq.mu.Lock()
	defer bq.mu.Unlock()

	if len(bq.items) == 0 {
		var zero E // Create the zero value for type E
		// Cannot return item if empty
		return zero, false, container.QueueIsBelowLowWatermark
	}

	item := bq.items[0]
	// Efficiently remove the first element (avoids memory leak)
	bq.items[0] = *new(E) // Assign zero value to prevent memory leak if E is a pointer type
	bq.items = bq.items[1:]

	return item, true, bq.getState()
}

// Len returns the number of items in the queue.
//
// It MUST NOT be called while holding the mutex to avoid deadlocks.
func (bq *waitable[E]) Len() int {
	bq.mu.Lock()
	defer bq.mu.Unlock()
	return len(bq.items)
}

// WaitChan returns the signal channel.
func (bq *waitable[E]) WaitChan() <-chan container.Unit {
	return bq.signal
}

// Close marks the queue as closed and closes the signal channel.
func (bq *waitable[E]) Close() {
	bq.mu.Lock()
	defer bq.mu.Unlock()
	if !bq.closed {
		bq.closed = true
		// Close the channel to permanently unblock any waiting Dequeue operations
		// and signal that no more items will arrive.
		close(bq.signal)
	}
}
