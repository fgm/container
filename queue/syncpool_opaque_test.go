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
