package binarysearchtree

import (
	"fmt"
	"testing"

	"github.com/fgm/container"
)

var (
	One, Two, Three, Four, Five, Six = 1, 2, 3, 4, 5, 6
)

// Simple builds this tree:
//       3
//      / \
//     2   4
//    /     \
//   1       5
func Simple() container.BinarySearchTree[int] {
	simple := Intrinsic[int]{}
	simple.Upsert(&Three, &Two, &Four, &One, &Five)
	return &simple
}

// HalfFull builds this tree, which contains all deletion cases
//       3
//      / \
//     2   5
//    /   / \
//   1   4   6
func HalfFull() container.BinarySearchTree[int] {
	hf := Intrinsic[int]{}
	hf.Upsert(&Three, &Two, &Five, &One, &Four, &Six)
	return &hf
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
	bst := Simple().(*Intrinsic[int])
	clone := bst.Clone().(*Intrinsic[int])
	input := bst.root
	output := clone.root
	checks := [...]struct {
		name             string
		expected, actual any
	}{
		{"root", *input.data, *output.data},
		{"left.data", *input.left.data, *output.left.data},
		{"left.left.data", *input.left.left.data, *output.left.left.data},
		{"left.left.left", input.left.left.left, output.left.left.left},
		{"left.left.right", input.left.left.right, output.left.left.right},
		{"left.right", input.left.right, output.left.right},
		{"right.data", *input.right.data, *output.right.data},
		{"right.left", input.right.left, output.right.left},
		{"right.right.data", *input.right.right.data, *output.right.right.data},
		{"right.right.left", input.right.right.left, output.right.right.left},
		{"right.right.right", input.right.right.right, output.right.right.right},
	}
	for _, check := range checks {
		t.Run(check.name, func(t *testing.T) {
			if check.actual != check.expected {
				t.Fatalf("got %v but expected %d", check.actual, check.expected)
			}
		})
	}
}

func TestTree_Delete(t *testing.T) {
	checks := [...]struct {
		name     string
		delendum int
	}{
		{"one: leaf", One},
		{"two: one child", Two},
		{"five: two children", Five},
		{"three: root with two children", Three},
	}

	for _, check := range checks {
		t.Run(check.name, func(t *testing.T) {
			bst := HalfFull()
			bst.Delete(&check.delendum)
		})
	}
}
