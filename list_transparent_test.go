package container

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestListQueue(t *testing.T) {
	checks := [...]struct {
		name     string
		input    []int
		dequeues int
		expectOK bool
		expected []int
	}{
		{"Q 0, DQ 1", nil, 1, false, []int{0}},
		{"Q 1, DQ 1", []int{1}, 1, true, []int{1}},
		{"Q 1, DQ 2", []int{1}, 2, false, []int{1, 0}},
		{"Q 2, DQ 2", []int{1, 2}, 2, true, []int{1, 2}},
	}

	for _, check := range checks {
		t.Run(check.name, func(t *testing.T) {
			q := NewListQueue[int](len(check.input))
			for _, x := range check.input {
				q.Enqueue(x)
			}
			actual := make([]int, 0, len(check.input))
			actualOK := true
			for i := 0; i < check.dequeues; i++ {
				e, ok := q.Dequeue()
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
			stack := NewListStack[int](len(check.input))
			for _, x := range check.input {
				stack.Push(x)
			}
			actual := make([]int, 0, len(check.input))
			actualOK := true
			for i := 0; i < check.pops; i++ {
				e, ok := stack.Pop()
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
