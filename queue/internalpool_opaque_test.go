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
