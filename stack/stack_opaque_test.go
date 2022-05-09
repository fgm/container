package stack_test

import (
	"testing"

	"github.com/fgm/container"
)

func testPop(t *testing.T, s container.Stack[int], expectCountable bool) {
	input := 1
	s.Push(input)
	c, ok := s.(container.Countable)
	if ok != expectCountable {
		t.Fatalf("stack countable is %t but expected %t", ok, expectCountable)
	}
	if ok && c.Len() != 1 {
		t.Fatalf("got len %d but expected 1", c.Len())
	}
	actual, ok := s.Pop()
	if !ok {
		t.Fatalf("failed popping pushed element")
	}
	if actual != input {
		t.Fatalf("popped elements differing from pushed one")
	}
	_, ok = s.Pop()
	if ok {
		t.Fatalf("successfully popped empty stack")
	}

}
