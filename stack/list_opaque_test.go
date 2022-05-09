package stack_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/fgm/container/stack"
)

func BenchmarkListStack_Pop(b *testing.B) {
	var s = stack.NewListStack[int](b.N)

	for i := 0; i < b.N; i++ {
		s.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stack.N, _ = s.Pop()
	}
	b.StopTimer()
}

func BenchmarkListStack_Push(b *testing.B) {
	var s = stack.NewListStack[int](b.N)

	for i := 0; i < b.N; i++ {
		s.Push(i)
	}
	b.StopTimer()
}

func TestListStack(t *testing.T) {
	checks := [...]struct {
		name     string
		input    []int
		pops     int
		expectOK bool
		expected []int
	}{
		{"Push 0, Pop 1", nil, 1, false, []int{0}},
		{"Push 1, Pop 1", []int{1}, 1, true, []int{1}},
		{"Push 1, Pop 2", []int{1}, 2, false, []int{1, 0}},
		{"Push 2, Pop 2", []int{1, 2}, 2, true, []int{2, 1}},
	}

	for _, check := range checks {
		t.Run(check.name, func(t *testing.T) {
			s := stack.NewListStack[int](len(check.input))
			for _, x := range check.input {
				s.Push(x)
			}
			actual := make([]int, 0, len(check.input))
			actualOK := true
			for i := 0; i < check.pops; i++ {
				e, ok := s.Pop()
				actual = append(actual, e)
				if !ok {
					actualOK = false
				}
			}
			if actualOK != check.expectOK {
				t.Fatalf("got OK %t but expected %t", actualOK, check.expectOK)
			}
			if !cmp.Equal(actual, check.expected) {
				t.Fatalf("unexpected result: %s", cmp.Diff(actual, check.expected))
			}
		})
	}
}

func TestListStack_Pop(t *testing.T) {
	testPop(t, stack.NewListStack[int](0), false)
}
