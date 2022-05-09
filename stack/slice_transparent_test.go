package stack

import (
	"testing"

	"github.com/fgm/container"
)

func BenchmarkSliceStack_Pop_raw(b *testing.B) {
	var stack container.Stack[int] = &sliceStack[int]{}

	for i := 0; i < b.N; i++ {
		stack.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		N, _ = stack.Pop()
	}
	b.StopTimer()
}

func BenchmarkSliceStack_Push_raw(b *testing.B) {
	var stack container.Stack[int] = &sliceStack[int]{}

	for i := 0; i < b.N; i++ {
		stack.Push(i)
	}
	b.StopTimer()
}
