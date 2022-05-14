package binarysearchtree_test

import (
	"strconv"
	"testing"

	"github.com/fgm/container/binarysearchtree"
)

func TestBST_nil(t *testing.T) {
	var bst *binarysearchtree.Tree[int]
	bst.WalkInOrder(binarysearchtree.P)
	bst.WalkPostOrder(binarysearchtree.P)
	bst.WalkPreOrder(binarysearchtree.P)
	bst.Upsert(nil)
	// Output:
}

func TestBST_Upsert(t *testing.T) {
	bst := binarysearchtree.Simple()
	actual := bst.Upsert(&binarysearchtree.One)
	if len(actual) != 1 {
		t.Fatalf("expected overwriting upsert to return one value, got %v", actual)
	}
	if *actual[0] != binarysearchtree.One {
		t.Fatalf("expected overwriting upsert to return %d, got %d", binarysearchtree.One, *actual[0])
	}

	actual = bst.Upsert(&binarysearchtree.Six)
	if len(actual) != 1 {
		t.Fatalf("expected non-overwriting upsert to return one value, got %v", actual)
	}
	if actual[0] != nil {
		t.Fatalf("expected non-overwriting upsert to return one nil, got %v", actual[0])
	}
}

func TestBST_IndexOf(t *testing.T) {
	bst := binarysearchtree.Simple()
	checks := [...]struct {
		input         int
		expectedOK    bool
		expectedIndex int
	}{
		{binarysearchtree.One, true, 0},
		{binarysearchtree.Two, true, 1},
		{binarysearchtree.Three, true, 2},
		{binarysearchtree.Four, true, 3},
		{binarysearchtree.Five, true, 4},
		{binarysearchtree.Six, false, 0},
	}
	for _, check := range checks {
		t.Run(strconv.Itoa(check.input), func(t *testing.T) {
			actualIndex, actualOK := bst.IndexOf(&check.input)
			if actualOK != check.expectedOK {
				t.Fatalf("%d found: %t but expected %t", check.input, actualOK, check.expectedOK)
			}
			if actualIndex != check.expectedIndex {
				t.Fatalf("%d at index %d but expected %d", check.input, actualIndex, check.expectedIndex)
			}
		})
	}
}
