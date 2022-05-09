package queue

import (
	"testing"

	"github.com/fgm/container"
)

func BenchmarkSliceQueue_Dequeue_raw(b *testing.B) {
	var queue container.Queue[int] = &sliceQueue[int]{}

	for i := 0; i < b.N; i++ {
		queue.Enqueue(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		N, _ = queue.Dequeue()
	}
	b.StopTimer()
}

func BenchmarkSliceQueue_Enqueue_raw(b *testing.B) {
	var queue container.Queue[int] = &sliceQueue[int]{}

	for i := 0; i < b.N; i++ {
		queue.Enqueue(i)
	}
	b.StopTimer()
}
