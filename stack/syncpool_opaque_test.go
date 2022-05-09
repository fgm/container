package stack_test

import (
	"testing"

	"github.com/fgm/container/stack"
)

func BenchmarkListSyncPoolStack_Pop_raw(b *testing.B) {
	var s = stack.NewListSyncPoolStack[int](0)

	for i := 0; i < b.N; i++ {
		s.Push(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stack.N, _ = s.Pop()
	}
	b.StopTimer()
}

func BenchmarkListSyncPoolStack_Pop_prealloc(b *testing.B) {
	var s = stack.NewListSyncPoolStack[int](b.N)

	for i := 0; i < b.N; i++ {
		s.Push(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stack.N, _ = s.Pop()
	}
	b.StopTimer()
}

func BenchmarkListSyncPoolStack_Push_raw(b *testing.B) {
	var s = stack.NewListSyncPoolStack[int](0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Push(i)
	}
	b.StopTimer()
}

func BenchmarkListSyncPoolStack_Push_prealloc(b *testing.B) {
	var s = stack.NewListSyncPoolStack[int](b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Push(i)
	}
	b.StopTimer()
}

func TestListSyncPoolStack_Pop(t *testing.T) {
	testPop(t, stack.NewListSyncPoolStack[int](0), false)
}
