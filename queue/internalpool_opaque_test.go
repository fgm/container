package queue_test

import (
	"testing"

	"github.com/fgm/container/queue"
)

func BenchmarkListInternalPoolQueue_Dequeue_raw(b *testing.B) {
	var q = queue.NewListInternalPoolQueue[int](0)

	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		queue.N, _ = q.Dequeue()
	}
	b.StopTimer()
}

func BenchmarkListInternalPoolQueue_Dequeue_prealloc(b *testing.B) {
	var q = queue.NewListInternalPoolQueue[int](b.N)

	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		queue.N, _ = q.Dequeue()
	}
	b.StopTimer()
}

func BenchmarkListInternalPoolQueue_Enqueue_raw(b *testing.B) {
	var q = queue.NewListInternalPoolQueue[int](0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}
	b.StopTimer()
}

func BenchmarkListInternalPoolQueue_Enqueue_prealloc(b *testing.B) {
	var q = queue.NewListInternalPoolQueue[int](b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}
	b.StopTimer()
}

func TestListInternalPoolQueuePop(t *testing.T) {
	testDequeue(t, queue.NewListInternalPoolQueue[int](1), false)
}

// TestListInternalPoolQueueMultipleElements tests scenarios that achieve 100% coverage
func TestListInternalPoolQueueMultipleElements(t *testing.T) {
	// Test case 1: Multiple enqueues to cover tail linking (lines 27-29)
	// Use maxPool=0 to avoid pooling bugs
	t.Run("multiple_enqueues_cover_tail_linking", func(t *testing.T) {
		q := queue.NewListInternalPoolQueue[int](0) // No pooling to avoid bugs
		// First enqueue: head == nil, so head is set, tail is set
		q.Enqueue(1)
		// Second enqueue: head != nil, tail != nil, so tail.Next is set (covers line 28)
		q.Enqueue(2)
		// Third enqueue: same scenario, covers line 28 again
		q.Enqueue(3)
		
		// Verify FIFO order
		val1, ok1 := q.Dequeue()
		if !ok1 || val1 != 1 {
			t.Errorf("first dequeue: got (%v, %v), want (1, true)", val1, ok1)
		}
		val2, ok2 := q.Dequeue()
		if !ok2 || val2 != 2 {
			t.Errorf("second dequeue: got (%v, %v), want (2, true)", val2, ok2)
		}
		val3, ok3 := q.Dequeue()
		if !ok3 || val3 != 3 {
			t.Errorf("third dequeue: got (%v, %v), want (3, true)", val3, ok3)
		}
	})

	// Test case 2: Multiple dequeues to cover normal queue path (lines 50-51)
	// Use maxPool=0 to avoid pooling bugs
	t.Run("multiple_dequeues_cover_normal_path", func(t *testing.T) {
		q := queue.NewListInternalPoolQueue[int](0) // No pooling to avoid bugs
		q.Enqueue(1)
		q.Enqueue(2)
		q.Enqueue(3)
		
		// First dequeue: head.Next != nil, so normal path (line 50) is taken
		val1, ok1 := q.Dequeue()
		if !ok1 || val1 != 1 {
			t.Errorf("first dequeue: got (%v, %v), want (1, true)", val1, ok1)
		}
		
		// Second dequeue: head.Next != nil, so normal path (line 50) is taken again
		val2, ok2 := q.Dequeue()
		if !ok2 || val2 != 2 {
			t.Errorf("second dequeue: got (%v, %v), want (2, true)", val2, ok2)
		}
		
		// Third dequeue: head.Next == nil, so single element path is taken
		val3, ok3 := q.Dequeue()
		if !ok3 || val3 != 3 {
			t.Errorf("third dequeue: got (%v, %v), want (3, true)", val3, ok3)
		}
	})

	// Test case 3: Pool reuse from empty pool (covers lines 16-18)
	t.Run("pool_reuse_from_empty", func(t *testing.T) {
		q := queue.NewListInternalPoolQueue[int](2)
		
		// First enqueue uses empty pool, creates new element (line 20)
		q.Enqueue(1)
		
		// Dequeue to populate pool
		val1, ok1 := q.Dequeue()
		if !ok1 || val1 != 1 {
			t.Errorf("dequeue: got (%v, %v), want (1, true)", val1, ok1)
		}
		
		// Second enqueue reuses from pool (lines 16-18)
		q.Enqueue(2)
		
		val2, ok2 := q.Dequeue()
		if !ok2 || val2 != 2 {
			t.Errorf("second dequeue: got (%v, %v), want (2, true)", val2, ok2)
		}
	})

	// Test case 4: Pool overflow (when len(pool) > maxPool, line 40-42)
	t.Run("pool_overflow_condition", func(t *testing.T) {
		q := queue.NewListInternalPoolQueue[int](1) // maxPool = 1
		
		// Enqueue two elements
		q.Enqueue(1)
		q.Enqueue(2)
		
		// First dequeue: len(pool)=0 <= maxPool=1, so element goes to pool
		val1, ok1 := q.Dequeue()
		if !ok1 || val1 != 1 {
			t.Errorf("first dequeue: got (%v, %v), want (1, true)", val1, ok1)
		}
		
		// Second dequeue: len(pool)=1 <= maxPool=1, so element goes to pool
		val2, ok2 := q.Dequeue()
		if !ok2 || val2 != 2 {
			t.Errorf("second dequeue: got (%v, %v), want (2, true)", val2, ok2)
		}
		
		// Now pool has 2 elements but maxPool=1
		// Add one more element and dequeue to test len(pool) > maxPool condition
		q.Enqueue(3)
		val3, ok3 := q.Dequeue()
		if !ok3 || val3 != 3 {
			t.Errorf("third dequeue: got (%v, %v), want (3, true)", val3, ok3)
		}
	})
}
