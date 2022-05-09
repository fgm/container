package queue_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/fgm/container/queue"
)

func BenchmarkListQueue_Dequeue(b *testing.B) {
	var q = queue.NewListQueue[int](b.N)

	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		queue.N, _ = q.Dequeue()
	}
	b.StopTimer()
}

func BenchmarkListQueue_Enqueue(b *testing.B) {
	var q = queue.NewListQueue[int](b.N)

	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}
	b.StopTimer()
}

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
			q := queue.NewListQueue[int](len(check.input))
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

func TestListQueuePop(t *testing.T) {
	testDequeue(t, queue.NewListQueue[int](1), false)
}
