package stack_test

import (
	"testing"

	"github.com/fgm/container/stack"
)

func BenchmarkListInternalPoolStack_Pop_raw(b *testing.B) {
	var s = stack.NewListInternalPoolStack[int](0)

	for i := 0; i < b.N; i++ {
		s.Push(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stack.N, _ = s.Pop()
	}
	b.StopTimer()
}

func BenchmarkListInternalPoolStack_Pop_prealloc(b *testing.B) {
	var s = stack.NewListInternalPoolStack[int](b.N)

	for i := 0; i < b.N; i++ {
		s.Push(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stack.N, _ = s.Pop()
	}
	b.StopTimer()
}

func BenchmarkListInternalPoolStack_Push_raw(b *testing.B) {
	var s = stack.NewListInternalPoolStack[int](0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Push(i)
	}
	b.StopTimer()
}

func BenchmarkListInternalPoolStack_Push_prealloc(b *testing.B) {
	var s = stack.NewListInternalPoolStack[int](b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Push(i)
	}
	b.StopTimer()
}

func TestListInternalPoolStackPop(t *testing.T) {
	testPop(t, stack.NewListInternalPoolStack[int](0), false)
}
