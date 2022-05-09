package queue_test

import (
	"testing"

	"github.com/fgm/container"
)

func testDequeue(t *testing.T, s container.Queue[int], expectCountable bool) {
	input := 1
	s.Enqueue(input)
	c, ok := s.(container.Countable)
	if ok != expectCountable {
		t.Fatalf("queue countable is %t but expected %t", ok, expectCountable)
	}
	if ok && c.Len() != 1 {
		t.Fatalf("got len %d but expected 1", c.Len())
	}
	actual, ok := s.Dequeue()
	if !ok {
		t.Fatalf("failed dequeueing enqueued element")
	}
	if actual != input {
		t.Fatalf("popped elements differing from pushed one")
	}
	_, ok = s.Dequeue()
	if ok {
		t.Fatalf("successfully dequeued from empty queue")
	}

}
