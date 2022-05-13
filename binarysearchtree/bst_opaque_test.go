package binarysearchtree_test

import (
	"strconv"
	"testing"

	"github.com/fgm/container/binarysearchtree"
)

func TestBST_nil(t *testing.T) {
	var bst *binarysearchtree.BST[int]
	bst.WalkInOrder(binarysearchtree.P)
	bst.WalkPostOrder(binarysearchtree.P)
	bst.WalkPreOrder(binarysearchtree.P)
	bst.Upsert(nil)
	// Output:
}

func TestBST_Upsert(t *testing.T) {
	one, six := 1, 6
	bst := binarysearchtree.Simple()
	actual := bst.Upsert(&one)
	if actual == nil {
		t.Fatalf("expected overwriting upsert to return value, got nil")
	}
	if *actual != one {
		t.Fatalf("expected overwriting upsert to return %d, got %d", one, *actual)
	}

	actual = bst.Upsert(&six)
	if actual != nil {
		t.Fatalf("expected non-overwriting upsert to return nil, got %d", *actual)
	}
}

func TestBST_IndexOf(t *testing.T) {
	bst := binarysearchtree.Simple()
	checks := []struct {
		input         int
		expectedOK    bool
		expectedIndex int
	}{
		{1, true, 0},
		{2, true, 1},
		{3, true, 2},
		{4, true, 3},
		{5, true, 4},
		{6, false, 0},
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
