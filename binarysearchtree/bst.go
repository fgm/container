package binarysearchtree

import (
	"golang.org/x/exp/constraints"
)

type BST[E constraints.Ordered] struct {
	data        *E
	left, right *BST[E]
}

// Upsert adds a value to the tree, replacing and returning the previous one if any.
// If none existed, it returns nil.
func (t *BST[E]) Upsert(e *E) *E {
	switch {
	case t == nil, e == nil:
		return nil
	case t.data == nil:
		t.data = e
		return nil
	case *e == *t.data:
		past := t.data
		t.data = e
		return past
	case *e > *t.data:
		if t.right == nil {
			t.right = &BST[E]{data: e}
			return nil
		}
		return t.right.Upsert(e)
	// The default case only covers the last "*e < *t.data" case, but if we only
	// use that clause, the compiler thinks some cases are not covered.
	default:
		if t.left == nil {
			t.left = &BST[E]{data: e}
			return nil
		}
		return t.left.Upsert(e)
	}
}

func (t *BST[E]) Delete(e *E) {
	// leaf: do it
	// one child: promote it
	// two children: promote the leftmost child of the right tree as the root.
	// If it had a right child (can't have a left child since it is rightmost), promote it.
}

// IndexOf returns the position of the value among those in the tree.
// If the value cannot be found, it will return 0, false, otherwise the position
// starting at 0, and true.
func (t *BST[E]) IndexOf(e *E) (int, bool) {
	index, found := 0, false
	t.WalkInOrder(func(x *E) {
		if *x == *e {
			found = true
		}
		if !found {
			index++
		}
	})
	if !found {
		index = 0
	}
	return index, found
}

// WalkInOrder in useful for search and listing nodes in order.
func (t *BST[E]) WalkInOrder(fn func(e *E)) {
	if t == nil {
		return
	}
	if t.left != nil {
		t.left.WalkInOrder(fn)
	}
	fn(t.data)
	if t.right != nil {
		t.right.WalkInOrder(fn)
	}
}

// WalkPostOrder in useful for deleting subtrees.
func (t *BST[E]) WalkPostOrder(fn func(e *E)) {
	if t == nil {
		return
	}
	if t.left != nil {
		t.left.WalkPostOrder(fn)
	}
	if t.right != nil {
		t.right.WalkPostOrder(fn)
	}
	fn(t.data)
}

// WalkPreOrder is useful to clone the tree.
func (t *BST[E]) WalkPreOrder(fn func(e *E)) {
	if t == nil {
		return
	}
	fn(t.data)
	if t.left != nil {
		t.left.WalkPreOrder(fn)
	}
	if t.right != nil {
		t.right.WalkPreOrder(fn)
	}
}

func (t *BST[E]) Clone() *BST[E] {
	clone := &BST[E]{}
	t.WalkPreOrder(func(e *E) {
		clone.Upsert(e)
	})
	return clone
}
