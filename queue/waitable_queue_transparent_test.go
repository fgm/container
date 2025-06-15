package queue

import (
	"testing"

	"github.com/fgm/container"
)

const (
	WQCap   = 10
	WQLow   = 2
	WQHigh  = 8
	WQInput = 42
)

func TestWaitable_Enqueue(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		initialItems  []int
		expectPanic   bool
		expectedState container.WaitableQueueState
		setup         func(*waitable[int])
	}{
		{
			name:          "enqueue to empty queue",
			initialItems:  []int{},
			expectedState: container.QueueIsBelowLowWatermark,
		},
		{
			name:          "enqueue to reach low watermark",
			initialItems:  []int{1},
			expectedState: container.QueueIsBelowLowWatermark,
		},
		{
			name:          "enqueue to nominal queue",
			initialItems:  []int{1, 2, 3},
			expectedState: container.QueueIsNominal,
		},
		{
			name:          "enqueue to reach high watermark",
			initialItems:  []int{1, 2, 3, 4, 5, 6, 7},
			expectedState: container.QueueIsAboveHighWatermark,
		},
		{
			name:          "enqueue to reach saturation",
			initialItems:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			expectedState: container.QueueIsNearSaturation,
		},
		{
			name:         "enqueue to closed queue",
			initialItems: []int{1, 2, 3},
			expectPanic:  true,
			setup:        func(q *waitable[int]) { q.Close() },
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			q, err := NewWaitableQueue[int](WQCap, WQLow, WQHigh)
			if err != nil {
				t.Fatalf("Failed to create queue: %v", err)
			}
			wq := q.(*waitable[int])

			// Add initial items
			wq.items = append(wq.items, test.initialItems...)

			// Apply the setup function if provided.
			if test.setup != nil {
				test.setup(wq)
			}

			// Handle panic cases
			if test.expectPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("Enqueue() expected panic but didn't get one")
					}
				}()
			}

			// Call Enqueue and verify the state.
			state := q.Enqueue(WQInput)

			// Only check the state if we don't expect a panic.
			if test.expectPanic {
				return
			}
			if state != test.expectedState {
				t.Errorf("Enqueue() returned state = %v, want %v", state, test.expectedState)
			}

			// Verify an item was actually added
			if wq.items[len(wq.items)-1] != WQInput {
				t.Errorf("Enqueue() failed to add item to queue")
			}

			// Verify that the signal channel has a value if the queue was empty before Enqueue.
			if len(test.initialItems) == 0 {
				select {
				case <-wq.signal:
					// Signal received, as expected
				default:
					t.Errorf("Enqueue() failed to send signal when queue was empty")
				}
			}
		})
	}
}
