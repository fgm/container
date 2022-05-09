package stack_test

import (
	"testing"

	"github.com/fgm/container/stack"
)

func BenchmarkSliceStack_Pop_prealloc(b *testing.B) {
	var s = stack.NewSliceStack[int](b.N)

	for i := 0; i < b.N; i++ {
		s.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stack.N, _ = s.Pop()
	}
	b.StopTimer()
}

func BenchmarkSliceStack_Push_prealloc(b *testing.B) {
	var s = stack.NewSliceStack[int](b.N)

	for i := 0; i < b.N; i++ {
		s.Push(i)
	}
	b.StopTimer()
}

func TestSliceStack_Pop(t *testing.T) {
	testPop(t, stack.NewSliceStack[int](0), true)
}
