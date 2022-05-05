package container_test

import (
	"testing"

	"github.com/fgm/container"
)

func BenchmarkSliceQueue_Dequeue_prealloc(b *testing.B) {
	var queue = container.NewSliceQueue[int](b.N)

	for i := 0; i < b.N; i++ {
		queue.Enqueue(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n, _ = queue.Dequeue()
	}
	b.StopTimer()
}

func BenchmarkSliceQueue_Enqueue_prealloc(b *testing.B) {
	var queue = container.NewSliceQueue[int](b.N)

	for i := 0; i < b.N; i++ {
		queue.Enqueue(i)
	}
	b.StopTimer()
}

func BenchmarkSliceStack_Pop_prealloc(b *testing.B) {
	var stack = container.NewSliceStack[int](b.N)

	for i := 0; i < b.N; i++ {
		stack.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n, _ = stack.Pop()
	}
	b.StopTimer()
}

func BenchmarkSliceStack_Push_prealloc(b *testing.B) {
	var stack = container.NewSliceStack[int](b.N)

	for i := 0; i < b.N; i++ {
		stack.Push(i)
	}
	b.StopTimer()
}
