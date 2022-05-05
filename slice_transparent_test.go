package container

import (
	"testing"
)

func BenchmarkSliceQueue_Dequeue_raw(b *testing.B) {
	var queue Queue[int] = &sliceQueue[int]{}

	for i := 0; i < b.N; i++ {
		queue.Enqueue(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n, _ = queue.Dequeue()
	}
	b.StopTimer()
}

func BenchmarkSliceQueue_Enqueue_raw(b *testing.B) {
	var queue Queue[int] = &sliceQueue[int]{}

	for i := 0; i < b.N; i++ {
		queue.Enqueue(i)
	}
	b.StopTimer()
}

func BenchmarkSliceStack_Pop_raw(b *testing.B) {
	var stack Stack[int] = &sliceStack[int]{}

	for i := 0; i < b.N; i++ {
		stack.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n, _ = stack.Pop()
	}
	b.StopTimer()
}

func BenchmarkSliceStack_Push_raw(b *testing.B) {
	var stack Stack[int] = &sliceStack[int]{}

	for i := 0; i < b.N; i++ {
		stack.Push(i)
	}
	b.StopTimer()
}
