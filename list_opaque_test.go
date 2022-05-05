package container_test

import (
	"testing"

	"github.com/fgm/container"
)

func BenchmarkListQueue_Dequeue(b *testing.B) {
	var queue = container.NewListQueue[int](b.N)

	for i := 0; i < b.N; i++ {
		queue.Enqueue(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n, _ = queue.Dequeue()
	}
	b.StopTimer()
}

func BenchmarkListQueue_Enqueue(b *testing.B) {
	var queue = container.NewListQueue[int](b.N)

	for i := 0; i < b.N; i++ {
		queue.Enqueue(i)
	}
	b.StopTimer()
}

func BenchmarkListStack_Pop(b *testing.B) {
	var stack = container.NewListStack[int](b.N)

	for i := 0; i < b.N; i++ {
		stack.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n, _ = stack.Pop()
	}
	b.StopTimer()
}

func BenchmarkListStack_Push(b *testing.B) {
	var stack = container.NewListStack[int](b.N)

	for i := 0; i < b.N; i++ {
		stack.Push(i)
	}
	b.StopTimer()
}
