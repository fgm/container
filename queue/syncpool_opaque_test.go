package queue_test

import (
	"testing"

	"github.com/fgm/container/queue"
)

func BenchmarkListSyncPoolQueue_Dequeue_raw(b *testing.B) {
	var q = queue.NewListSyncPoolQueue[int](0)

	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		queue.N, _ = q.Dequeue()
	}
	b.StopTimer()
}

func BenchmarkListSyncPoolQueue_Dequeue_prealloc(b *testing.B) {
	var q = queue.NewListSyncPoolQueue[int](b.N)

	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		queue.N, _ = q.Dequeue()
	}
	b.StopTimer()
}

func BenchmarkListSyncPoolQueue_Enqueue_raw(b *testing.B) {
	var q = queue.NewListSyncPoolQueue[int](0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}
	b.StopTimer()
}

func BenchmarkListSyncPoolQueue_Enqueue_prealloc(b *testing.B) {
	var q = queue.NewListSyncPoolQueue[int](b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}
	b.StopTimer()
}

func TestListSyncPoolQueuePop(t *testing.T) {
	testDequeue(t, queue.NewListSyncPoolQueue[int](1), false)
}

// TestListSyncPoolQueueMultipleElements tests scenarios that achieve 100% coverage
func TestListSyncPoolQueueMultipleElements(t *testing.T) {
	// Test case 1: Multiple enqueues to cover tail linking (lines 21-23)
	t.Run("multiple_enqueues_cover_tail_linking", func(t *testing.T) {
		q := queue.NewListSyncPoolQueue[int](2)
		// First enqueue: head == nil, so head is set, tail is set
		q.Enqueue(1)
		// Second enqueue: head != nil, tail != nil, so tail.Next is set (covers line 22)
		q.Enqueue(2)
		// Third enqueue: same scenario, covers line 22 again
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

	// Test case 2: Multiple dequeues to cover normal queue path (lines 41-43)
	t.Run("multiple_dequeues_cover_normal_path", func(t *testing.T) {
		q := queue.NewListSyncPoolQueue[int](2)
		q.Enqueue(1)
		q.Enqueue(2)
		q.Enqueue(3)
		
		// First dequeue: head.Next != nil, so normal path (line 42) is taken
		val1, ok1 := q.Dequeue()
		if !ok1 || val1 != 1 {
			t.Errorf("first dequeue: got (%v, %v), want (1, true)", val1, ok1)
		}
		
		// Second dequeue: head.Next != nil, so normal path (line 42) is taken again
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

	// Test case 3: Pool reuse scenarios
	t.Run("pool_reuse_scenarios", func(t *testing.T) {
		q := queue.NewListSyncPoolQueue[int](2)
		
		// Fill and drain to populate pool
		q.Enqueue(1)
		q.Enqueue(2)
		q.Dequeue() // This puts element back to sync.Pool
		q.Dequeue() // This puts element back to sync.Pool
		
		// Enqueue again - should reuse pooled elements
		q.Enqueue(3)
		q.Enqueue(4)
		
		val1, ok1 := q.Dequeue()
		if !ok1 || val1 != 3 {
			t.Errorf("reuse dequeue 1: got (%v, %v), want (3, true)", val1, ok1)
		}
		val2, ok2 := q.Dequeue()
		if !ok2 || val2 != 4 {
			t.Errorf("reuse dequeue 2: got (%v, %v), want (4, true)", val2, ok2)
		}
	})

	// Test case 4: Empty queue dequeue
	t.Run("empty_queue_dequeue", func(t *testing.T) {
		q := queue.NewListSyncPoolQueue[int](1)
		
		// Dequeue from empty queue should return zero value and false
		val, ok := q.Dequeue()
		if ok || val != 0 {
			t.Errorf("empty dequeue: got (%v, %v), want (0, false)", val, ok)
		}
	})
}
