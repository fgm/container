package binarysearchtree

import (
	"golang.org/x/exp/constraints"
)

type walkCB[E constraints.Ordered] func(*E)

type node[E constraints.Ordered] struct {
	data        *E
	left, right *node[E]
}

func (n *node[E]) delete(e *E) *node[E] {
	switch {
	case n == nil, e == nil:
		n = nil
	case *e < *n.data:
		n.left = n.left.delete(e)
	case *e > *n.data:
		n.right = n.right.delete(e)
	// Matched childless node: just drop it.
	case n.left == nil && n.right == nil:
		n = nil
	// Matched node with only one right child: promote that child.
	case n.left == nil:
		n = n.right
	// Matched node with only one left child: promote that child.
	case n.right == nil:
		n = n.left
	// Matched node with two children: promote leftmost child of right child.
	//
	// We could also have promoted the rightmost child of the left child.
	default:
		promovendum := n.right // Cannot be nil: that case was handled ed above.
		for {
			if promovendum.left == nil {
				break
			}
			promovendum = promovendum.left // Not nil either, per previous statement.
		}
		n.data = promovendum.data
		n.right = n.right.delete(promovendum.data) // As the leftmost child, it won't have two children.
	}
	return n
}

func (n *node[E]) upsert(m *node[E]) *E {
	switch {
	case *m.data < *n.data:
		if n.left == nil {
			n.left = m
			return nil
		} else {
			return n.left.upsert(m)
		}
	case *m.data > *n.data:
		if n.right == nil {
			n.right = m
			return nil
		} else {
			return n.right.upsert(m)
		}
	default: // *m.data == *n.data
		data := n.data
		n.data = m.data
		return data
	}
}

func (n *node[E]) walkInOrder(cb walkCB[E]) {
	if n == nil {
		return
	}
	if n.left != nil {
		n.left.walkInOrder(cb)
	}
	cb(n.data)
	if n.right != nil {
		n.right.walkInOrder(cb)
	}
}

func (n *node[E]) walkPostOrder(cb walkCB[E]) {
	if n == nil {
		return
	}
	if n.left != nil {
		n.left.walkPostOrder(cb)
	}
	if n.right != nil {
		n.right.walkPostOrder(cb)
	}
	cb(n.data)
}

func (n *node[E]) walkPreOrder(cb walkCB[E]) {
	if n == nil {
		return
	}
	cb(n.data)
	if n.left != nil {
		n.left.walkPreOrder(cb)
	}
	if n.right != nil {
		n.right.walkPreOrder(cb)
	}
}

type Tree[E constraints.Ordered] struct {
	root *node[E]
}

// Upsert adds a value to the tree, replacing and returning the previous one if any.
// If none existed, it returns nil.
func (t *Tree[E]) Upsert(e ...*E) []*E {
	res := make([]*E, 0, len(e))
	for _, oneE := range e {
		n := &node[E]{data: oneE}

		switch {
		case t == nil, e == nil:
			res = append(res, nil)
		case t.root == nil:
			t.root = n
			res = append(res, nil)
		default:
			res = append(res, t.root.upsert(n))
		}
	}
	return res
}

func (t *Tree[E]) Delete(e *E) {
	if t == nil || e == nil {
		return
	}
	t.root.delete(e)
	// two children: promote the leftmost child of the right tree as the root.
	// If it had a right child (can't have a left child since it is rightmost), promote it.
}

// IndexOf returns the position of the value among those in the tree.
// If the value cannot be found, it will return 0, false, otherwise the position
// starting at 0, and true.
func (t *Tree[E]) IndexOf(e *E) (int, bool) {
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
func (t *Tree[E]) WalkInOrder(cb walkCB[E]) {
	if t == nil {
		return
	}
	t.root.walkInOrder(cb)
}

// WalkPostOrder in useful for deleting subtrees.
func (t *Tree[E]) WalkPostOrder(fn func(e *E)) {
	if t == nil {
		return
	}
	t.root.walkPostOrder(fn)
}

// WalkPreOrder is useful to clone the tree.
func (t *Tree[E]) WalkPreOrder(cb walkCB[E]) {
	if t == nil {
		return
	}
	t.root.walkPreOrder(cb)
}

func (t *Tree[E]) Clone() *Tree[E] {
	clone := &Tree[E]{}
	t.WalkPreOrder(func(e *E) {
		clone.Upsert(e)
	})
	return clone
}
