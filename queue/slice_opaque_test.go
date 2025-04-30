package queue_test

import (
	"testing"

	"github.com/fgm/container/queue"
)

func BenchmarkSliceQueue_Dequeue_prealloc(b *testing.B) {
	var q = queue.NewSliceQueue[int](b.N)

	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		queue.N, _ = q.Dequeue()
	}
	b.StopTimer()
}

func BenchmarkSliceQueue_Enqueue_prealloc(b *testing.B) {
	var q = queue.NewSliceQueue[int](b.N)

	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}
	b.StopTimer()
}

func BenchmarkSliceQueue_Enqueue_preallocSad(b *testing.B) {
	var q = queue.NewSliceQueue[int](b.N / 10)

	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}
	b.StopTimer()
}

func TestSliceQueuePop(t *testing.T) {
	testDequeue(t, queue.NewSliceQueue[int](1), true)
}
