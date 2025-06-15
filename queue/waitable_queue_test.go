package queue_test

import (
	"context"
	"errors"
	"log"
	"sort"
	"sync"
	"testing"
	"time"

	// Import the container interface package
	"github.com/fgm/container"
	"github.com/fgm/container/queue"
)

// checkQueueAndCollect sends processed items to an output channel for testing verification.
func checkQueueAndCollect[E any](ctx context.Context, q container.WaitableQueue[E], output chan<- E, wg *sync.WaitGroup) {
	defer wg.Done() // Signal completion when this goroutine exits

	for {
		// First, try to dequeue in case items are already available.
		item, ok, _ := q.Dequeue()
		if ok {
			select {
			case output <- item:
				// Item sent successfully.
			case <-ctx.Done():
				// Context was cancelled while trying to send: exit.
				log.Printf("Consumer exiting: context cancelled while sending item.\n")
				return
			}
			continue // Successfully dequeued, try again immediately
		}

		// Queue is currently empty, wait for a signal or context cancellation.
		select {
		case <-ctx.Done():
			// Context was cancelled, stop processing.
			// Attempt one last drain, respecting context cancellation during send
			log.Printf("Consumer exiting: context cancelled while waiting.\n")
			for {
				item, ok, _ := q.Dequeue()
				if !ok {
					log.Printf("Consumer exiting: queue drained after context cancellation.\n")
					return // Queue is now drained.
				}
				select {
				case output <- item:
					// Final item sent.
				case <-ctx.Done():
					log.Printf("Consumer exiting: context cancelled during final drain.\n")
					return // Context was cancelled during drain.
				}
			}

		case _, sigOK := <-q.WaitChan():
			// Received a signal OR the signal channel was closed unexpectedly.
			if !sigOK {
				// Channel closing indicates an issue or unexpected state.
				// We'll treat it like a cancellation signal for shutdown.
				log.Printf("Consumer exiting: WaitChan closed unexpectedly.\n")
				// Perform a final drain similar to context cancellation
				for {
					item, ok, _ := q.Dequeue()
					if !ok {
						log.Printf("Consumer exiting: queue drained after WaitChan closed.\n")
						return // Queue is now drained.
					}
					select {
					case output <- item:
						// Final item sent.
					case <-ctx.Done():
						log.Printf("Consumer exiting: context cancelled during final drain after WaitChan closed.\n")
						return // Context was cancelled during drain.
					}
				}
			}
			// At this point, we know that a signal was received, because sigOK is true.
			// Loop back to the top to attempt Dequeue().
		}
	}
}

// TestConcurrentWaitableQueue tests the concurrent behavior of the WaitableQueue,
// in a producer / consumer scenario.
func TestConcurrentWaitableQueue(t *testing.T) {
	const (
		numProducers     = 5
		numConsumers     = 3
		itemsPerProducer = 20
		totalItems       = numProducers * itemsPerProducer
	)

	q, err := queue.NewWaitableQueue[int](totalItems*2, 0, totalItems*2)
	if err != nil {
		t.Fatalf("Failed to create WaitableQueue: %v", err)
	}

	// Context for managing consumer lifecycle
	ctx, cancel := context.WithCancel(t.Context())
	// No defer cancel() here, we cancel explicitly later

	var producerWg, consumerWg sync.WaitGroup

	// Channel to collect results, buffered to hold all items
	processedItemsChan := make(chan int, totalItems)
	// Channel to signal when all items have been collected
	collectionDone := make(chan container.Unit)

	// Start Consumers (checkQueue goroutines).
	consumerWg.Add(numConsumers)
	for i := 0; i < numConsumers; i++ {
		go checkQueueAndCollect(ctx, q, processedItemsChan, &consumerWg)
	}

	// Start Producers.
	producerWg.Add(numProducers)
	for i := 0; i < numProducers; i++ {
		go func(producerID int) {
			defer producerWg.Done()
			for j := 0; j < itemsPerProducer; j++ {
				item := producerID*1000 + j // Guarantee a unique item value
				_ = q.Enqueue(item)
			}
		}(i)
	}

	// Wait for producers, then signal collection goroutine.
	go func() {
		producerWg.Wait()
		t.Logf("All producers finished.")
		// All producers are done, but consumers might still be processing.
	}()

	// Collect results and signal when done.
	processedItems := make([]int, 0, totalItems)
	go func() {
		for item := range processedItemsChan {
			processedItems = append(processedItems, item)
			if len(processedItems) == totalItems {
				close(collectionDone) // Signal that all expected items are collected
				return                // Stop collecting from this goroutine
			}
		}
	}()

	// Wait until all expected items have been collected by the collection goroutine.
	t.Logf("Waiting for %d items to be collected...", totalItems)
	select {
	case <-collectionDone:
		t.Logf("All %d items collected.", totalItems)
	case <-time.After(10 * time.Second): // Add a timeout
		t.Fatalf("Timeout waiting for items to be collected. Collected %d items.", len(processedItems))
	}

	// Now that all items produced have been collected, cancel the context to stop consumers
	t.Logf("Cancelling context to stop consumers...")
	cancel()

	// Wait for all consumers to finish processing and exit
	consumerWaitDone := make(chan container.Unit)
	go func() {
		consumerWg.Wait()
		close(consumerWaitDone)
	}()

	select {
	case <-consumerWaitDone:
		t.Log("All consumers finished.")
	case <-time.After(5 * time.Second): // Timeout for consumers to stop
		t.Fatal("Timeout waiting for consumers to stop after cancellation.")
	}

	// Close the results channel *after* ensuring consumers are done
	// (though the collection goroutine already stopped reading)
	close(processedItemsChan)

	if len(processedItems) != totalItems {
		t.Errorf("Expected %d items processed, but got %d", totalItems, len(processedItems))
	}

	// Verify uniqueness and content
	sort.Ints(processedItems)
	expected := make([]int, 0, totalItems)
	seen := make(map[int]bool)
	for i := 0; i < numProducers; i++ {
		for j := 0; j < itemsPerProducer; j++ {
			item := i*1000 + j
			expected = append(expected, item)
			if seen[item] {
				t.Errorf("Duplicate item processed: %d", item)
			}
			seen[item] = true
		}
	}
	sort.Ints(expected) // Sort expected items

	// Check count just in case map logic missed something (unlikely).
	if len(processedItems) != len(seen) {
		t.Errorf("Processed item count (%d) does not match unique item count (%d)", len(processedItems), len(seen))
	}

	// Direct comparison of sorted slices.
	match := true
	if len(processedItems) != len(expected) {
		match = false // Should have been caught earlier, but double-check
	} else {
		for i := range processedItems {
			if processedItems[i] != expected[i] {
				match = false
				break
			}
		}
	}
	if !match {
		// Provide more details if mismatch occurs
		t.Errorf("Processed items do not match expected items.\nExpected (len %d): %v\nGot (len %d):      %v", len(expected), expected, len(processedItems), processedItems)
	} else {
		t.Logf("Successfully processed and verified %d items.", len(processedItems))
	}
}

func TestConcurrentWaitableQueue_ContextCancel(t *testing.T) {
	const (
		numConsumers   = 3
		itemsToEnqueue = 500 // High enough to ensure some items are likely still in the queue.
	)

	q, err := queue.NewWaitableQueue[int](itemsToEnqueue, 0, itemsToEnqueue)
	if err != nil {
		t.Fatalf("Failed to create WaitableQueue: %v", err)
	}

	ctx, cancel := context.WithCancel(t.Context())
	// No defer cancel() here: we use explicit cancellation.

	var consumerWg sync.WaitGroup
	processedItemsChan := make(chan int, itemsToEnqueue) // Buffer size >= items

	// Start Consumers
	consumerWg.Add(numConsumers)
	for i := 0; i < numConsumers; i++ {
		go checkQueueAndCollect(ctx, q, processedItemsChan, &consumerWg)
	}

	// Enqueue some items
	for i := 0; i < itemsToEnqueue; i++ {
		_ = q.Enqueue(i)
	}
	t.Logf("Enqueued %d items.", itemsToEnqueue)

	// Give consumers a little time to process some items
	// time.Sleep(100 * time.Millisecond) // Increased sleep slightly

	// Cancel the context
	t.Log("Cancelling context...")
	cancel() // Explicitly cancel

	// Wait for consumers to exit due to cancellation
	// Use a timeout channel to prevent test hanging indefinitely
	waitChan := make(chan container.Unit)
	go func() {
		consumerWg.Wait()
		close(waitChan)
	}()

	select {
	case <-waitChan:
		t.Logf("All consumers finished after context cancellation.")
	case <-time.After(2 * time.Second): // Timeout
		t.Fatal("Consumers did not finish within timeout after context cancellation")
	}
	// Close the results channel *after* consumers are done.
	close(processedItemsChan)

	// Collect results - we expect *fewer* than totalItems were enqueued,
	// because cancellation happened while items were likely still in the queue or being processed.
	processedCount := 0
	finalItems := []int{}
	for item := range processedItemsChan {
		processedCount++
		finalItems = append(finalItems, item) // Collect for logging if needed
	}

	switch {
	case processedCount == itemsToEnqueue:
		// This could happen if consumers were extremely fast and processed everything
		// before the cancel signal was effectively received and acted upon by all of them.
		// It's less likely with the sleep but possible. Consider it a pass, but log it.
		t.Logf("WARN: Processed all %d items despite cancellation (potentially very fast consumers or timing).", itemsToEnqueue)
	case processedCount == 0 && itemsToEnqueue > 0:
		// This might happen if cancellation was extremely fast relative to consumer startup/dequeue.
		t.Logf("Processed 0 items before cancellation (potentially very fast cancellation or slow consumers).")
	default:
		t.Logf("Successfully processed %d items before/during context cancellation (expected < %d). Items: %v", processedCount, itemsToEnqueue, finalItems)
	}
	// The main point is that the consumers stopped gracefully after cancellation.
}

func TestNewWaitableQueue(t *testing.T) {
	const (
		Cap  = 10
		Low  = 2
		High = 8
	)
	tests := [...]struct {
		name             string
		capacity, lo, hi int
		expectErr        error
	}{
		{"capacity below 0", -1, -1, -1, queue.ErrCapacityIsNegative},
		{"low watermark below 0", Cap, -1, -1, queue.ErrLowWatermarkIsNegative},
		{"high watermark below 0", Cap, Low, -1, queue.ErrHighWatermarkIsNegative},
		{"high watermark below low watermark", Cap, High, Low, queue.ErrHighWatermarkIsLessThanLowWatermark},
		{"happy path", Cap, Low, High, nil},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, actualErr := queue.NewWaitableQueue[int](test.capacity, test.lo, test.hi)
			switch {
			case test.expectErr == nil && actualErr != nil:
				t.Fatalf("Got %v but expected success", actualErr)
			case test.expectErr != nil:
				if !errors.Is(actualErr, test.expectErr) {
					t.Fatalf("Got %v but expected error %v", actualErr, test.expectErr)
				}
				if actual != nil {
					t.Fatalf("Got %#v but expected nil queue", actual)
				}
				return // Queue is nil, nothing more to check.
			}
			var wqs container.WaitableQueueState

			if wqs = actual.Enqueue(1); wqs != container.QueueIsBelowLowWatermark {
				t.Errorf("Got %s but expected new queue to be below low watermark", wqs)
			}
			// Enqueue just enough to be in nominal state.
			mid := (High + Low - 1) / 2 // We already enqueued one item.
			for range mid {
				wqs = actual.Enqueue(1)
			}
			if wqs != container.QueueIsNominal {
				t.Errorf("Got %s but expected half-allocated queue to be nominal", wqs)
			}
			// Fill the queue to reach saturation.
			for range Cap - mid {
				wqs = actual.Enqueue(1)
			}
			if wqs != container.QueueIsNearSaturation {
				t.Errorf("Got %s but expected full queue to be near saturation", wqs)
			}
		})
	}
}

func TestWaitable_Len(t *testing.T) {
	tests := []struct {
		name           string
		initialItems   int
		operations     func(q container.WaitableQueue[int])
		expectedLength int
	}{
		{
			name:           "empty queue",
			initialItems:   0,
			operations:     nil,
			expectedLength: 0,
		},
		{
			name:           "queue with items",
			initialItems:   5,
			operations:     nil,
			expectedLength: 5,
		},
		{
			name:         "enqueue operations",
			initialItems: 2,
			operations: func(q container.WaitableQueue[int]) {
				q.Enqueue(42)
				q.Enqueue(43)
				q.Enqueue(44)
			},
			expectedLength: 5, // 2 initial + 3 added
		},
		{
			name:         "dequeue operations",
			initialItems: 5,
			operations: func(q container.WaitableQueue[int]) {
				_, _, _ = q.Dequeue()
				_, _, _ = q.Dequeue()
			},
			expectedLength: 3, // 5 initial - 2 removed
		},
		{
			name:         "mixed operations",
			initialItems: 3,
			operations: func(q container.WaitableQueue[int]) {
				_, _, _ = q.Dequeue()
				q.Enqueue(42)
				q.Enqueue(43)
				_, _, _ = q.Dequeue()
			},
			expectedLength: 3, // 3 initial - 2 removed + 2 added
		},
		{
			name:         "concurrent operations",
			initialItems: 0,
			operations: func(q container.WaitableQueue[int]) {
				var wg sync.WaitGroup
				// Ajouter 10 éléments en concurrence
				wg.Add(10)
				for i := 0; i < 10; i++ {
					go func(value int) {
						defer wg.Done()
						q.Enqueue(value)
					}(i)
				}
				wg.Wait()
			},
			expectedLength: 10, // 0 initial + 10 added concurrently
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			q, err := queue.NewWaitableQueue[int](queue.WQCap, queue.WQLow, queue.WQHigh)
			if err != nil {
				t.Fatalf("Failed to create queue: %v", err)
			}
			actual, ok := q.(interface {
				container.WaitableQueue[int]
				container.Countable
			})
			if !ok {
				t.Fatalf("expected both WaitableQueue and Countable interface")
			}

			// Ajouter les éléments initiaux
			for i := 0; i < test.initialItems; i++ {
				actual.Enqueue(i)
			}

			// Exécuter les opérations du test si définies
			if test.operations != nil {
				test.operations(actual)
			}

			// Vérifier que Len() retourne la longueur attendue
			if length := actual.Len(); length != test.expectedLength {
				t.Errorf("Len() = %d, want %d", length, test.expectedLength)
			}
		})
	}
}
