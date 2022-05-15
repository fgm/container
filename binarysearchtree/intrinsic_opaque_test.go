package binarysearchtree_test

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/fgm/container"
	bst "github.com/fgm/container/binarysearchtree"
)

func TestIntrinsic_nil(t *testing.T) {
	var tree *bst.Intrinsic[int]
	tree.WalkInOrder(bst.P)
	tree.WalkPostOrder(bst.P)
	tree.WalkPreOrder(bst.P)
	tree.Upsert(nil)
	tree.Delete(nil)

	tree = &bst.Intrinsic[int]{}
	tree.WalkInOrder(bst.P)
	tree.WalkPostOrder(bst.P)
	tree.WalkPreOrder(bst.P)
	// Output:
}

func TestIntrinsic_Upsert(t *testing.T) {
	tree := bst.Simple()
	actual := tree.Upsert(&bst.One)
	if len(actual) != 1 {
		t.Fatalf("expected overwriting upsert to return one value, got %v", actual)
	}
	if *actual[0] != bst.One {
		t.Fatalf("expected overwriting upsert to return %d, got %d", bst.One, *actual[0])
	}

	actual = tree.Upsert(&bst.Six)
	if len(actual) != 1 {
		t.Fatalf("expected non-overwriting upsert to return one value, got %v", actual)
	}
	if actual[0] != nil {
		t.Fatalf("expected non-overwriting upsert to return one nil, got %v", actual[0])
	}
}

func TestIntrinsic_IndexOf(t *testing.T) {
	tree := bst.Simple()
	checks := [...]struct {
		input         int
		expectedOK    bool
		expectedIndex int
	}{
		{bst.One, true, 0},
		{bst.Two, true, 1},
		{bst.Three, true, 2},
		{bst.Four, true, 3},
		{bst.Five, true, 4},
		{bst.Six, false, 0},
	}
	for _, check := range checks {
		t.Run(strconv.Itoa(check.input), func(t *testing.T) {
			actualIndex, actualOK := tree.IndexOf(&check.input)
			if actualOK != check.expectedOK {
				t.Fatalf("%d found: %t but expected %t", check.input, actualOK, check.expectedOK)
			}
			if actualIndex != check.expectedIndex {
				t.Fatalf("%d at index %d but expected %d", check.input, actualIndex, check.expectedIndex)
			}
		})
	}
}

func TestIntrinsic_Len(t *testing.T) {
	si := bst.Simple().(container.Enumerable[int]).Elements()
	hf := bst.HalfFull().(container.Enumerable[int]).Elements()

	checks := [...]struct {
		name      string
		input     []int
		deletions []int
		expected  int
	}{
		{"empty", nil, nil, 0},
		{"simple", si, nil, 5},
		{"half-full", hf, nil, 6},
		{"overwrite element", append(si, bst.Three), nil, 5},
		{"delete nonexistent", si, []int{bst.Six}, 5},
		{"delete existing childless", si, []int{bst.One}, 4},
		{"delete existing with 1 left child", si, []int{bst.Two}, 4},
		{"delete existing with 1 right child", si, []int{bst.Four}, 4},
		{"delete existing with 2 children", hf, []int{bst.Three}, 5},
	}
	for _, check := range checks {
		t.Run(check.name, func(t *testing.T) {
			tree := bst.Intrinsic[int]{}
			// In these loops, e is always the same variable: without cloning,
			// each iteration reuses the same pointer, overwriting the tree.
			for _, e := range check.input {
				clone := e
				tree.Upsert(&clone)
			}
			for _, e := range check.deletions {
				clone := e
				tree.Delete(&clone)
			}
			if tree.Len() != check.expected {
				t.Fatalf("Found len %d, but expected %d", tree.Len(), check.expected)
			}
		})
	}
}

func TestIntrinsic_Walk_canceling(t *testing.T) {
	tree := bst.Simple()
	checks := [...]struct {
		name                   string
		walker                 func(cb container.WalkCB[int]) error
		calls1, calls3, calls5 int
	}{
		{"in order", tree.WalkInOrder, 1, 3, 5},
		{"post order", tree.WalkPostOrder, 1, 5, 3},
		{"pre order", tree.WalkPreOrder, 3, 1, 5},
	}
	tree.WalkInOrder(func(e *int) error { log.Println(*e); return nil })
	stopper := func(at int) container.WalkCB[int] {
		called := 0
		return container.WalkCB[int](func(e *int) error {
			called++
			if *e == at {
				return fmt.Errorf("%d", called)
			}
			return nil
		})
	}
	for _, check := range checks {
		t.Run(check.name, func(t *testing.T) {
			for _, val := range []struct {
				input, expected int
			}{
				{1, check.calls1},
				{3, check.calls3},
				{5, check.calls5},
			} {
				err := check.walker(stopper(val.input))
				actual, err := strconv.Atoi(err.Error())
				if err != nil {
					t.Fatalf("unexpected non-int error: %v", err)
				}
				if actual != val.expected {
					t.Fatalf("got %d but expected %d", actual, val.expected)
				}
			}
		})
	}
}
