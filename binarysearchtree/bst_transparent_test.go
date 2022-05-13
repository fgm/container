package binarysearchtree

import (
	"fmt"
	"testing"
)

var (
	one, two, three, four, five = 1, 2, 3, 4, 5
)

func Simple() *BST[int] {
	return &BST[int]{
		data:  &three,
		left:  &BST[int]{data: &two, left: &BST[int]{data: &one}},
		right: &BST[int]{data: &four, right: &BST[int]{data: &five}},
	}
}

func P(e *int) {
	_, _ = fmt.Println(*e)
}

func ExampleBST_WalkInOrder() {
	bst := Simple()
	bst.WalkInOrder(P)
	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
}

func ExampleBST_WalkPostOrder() {
	bst := Simple()
	bst.WalkPostOrder(P)
	// Output:
	// 1
	// 2
	// 5
	// 4
	// 3
}

func ExampleBST_WalkPreOrder() {
	bst := Simple()
	bst.WalkPreOrder(P)
	// Output:
	// 3
	// 2
	// 1
	// 4
	// 5
}

func TestBST_Clone(t *testing.T) {
	bst := Simple()
	clone := bst.Clone()
	checks := []struct {
		name             string
		expected, actual any
	}{
		{"root", *bst.data, *clone.data},
		{"left.data", *bst.left.data, *clone.left.data},
		{"left.left.data", *bst.left.left.data, *clone.left.left.data},
		{"left.left.left", bst.left.left.left, clone.left.left.left},
		{"left.left.right", bst.left.left.right, clone.left.left.right},
		{"left.right", bst.left.right, clone.left.right},
		{"right.data", *bst.right.data, *clone.right.data},
		{"right.left", bst.right.left, clone.right.left},
		{"right.right.data", *bst.right.right.data, *clone.right.right.data},
		{"right.right.left", bst.right.right.left, clone.right.right.left},
		{"right.right.right", bst.right.right.right, clone.right.right.right},
	}
	for _, check := range checks {
		t.Run(check.name, func(t *testing.T) {
			if check.actual != check.expected {
				t.Fatalf("got %v but expected %d", check.actual, check.expected)
			}
		})
	}
}
