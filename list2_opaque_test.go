package container_test

import (
	"testing"

	"github.com/fgm/container"
)

func BenchmarkList2Queue_Dequeue_raw(b *testing.B) {
	var queue = container.NewList2Queue[int](0)

	for i := 0; i < b.N; i++ {
		queue.Enqueue(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n, _ = queue.Dequeue()
	}
	b.StopTimer()
}

func BenchmarkList2Queue_Dequeue_prealloc(b *testing.B) {
	var queue = container.NewList2Queue[int](b.N)

	for i := 0; i < b.N; i++ {
		queue.Enqueue(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n, _ = queue.Dequeue()
	}
	b.StopTimer()
}

func BenchmarkList2Queue_Enqueue_raw(b *testing.B) {
	var queue = container.NewList2Queue[int](0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		queue.Enqueue(i)
	}
	b.StopTimer()
}

func BenchmarkList2Queue_Enqueue_prealloc(b *testing.B) {
	var queue = container.NewList2Queue[int](b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		queue.Enqueue(i)
	}
	b.StopTimer()
}

func BenchmarkList2Stack_Pop_raw(b *testing.B) {
	var stack = container.NewList2Stack[int](0)

	for i := 0; i < b.N; i++ {
		stack.Push(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n, _ = stack.Pop()
	}
	b.StopTimer()
}

func BenchmarkList2Stack_Pop_prealloc(b *testing.B) {
	var stack = container.NewList2Stack[int](b.N)

	for i := 0; i < b.N; i++ {
		stack.Push(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n, _ = stack.Pop()
	}
	b.StopTimer()
}

func BenchmarkList2Stack_Push_raw(b *testing.B) {
	var stack = container.NewList2Stack[int](0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stack.Push(i)
	}
	b.StopTimer()
}

func BenchmarkList2Stack_Push_prealloc(b *testing.B) {
	var stack = container.NewList2Stack[int](b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stack.Push(i)
	}
	b.StopTimer()
}
