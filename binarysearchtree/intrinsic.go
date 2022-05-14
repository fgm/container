package binarysearchtree

import (
	"cmp"

	"github.com/fgm/container"
)

type node[E cmp.Ordered] struct {
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

func (n *node[E]) walkInOrder(cb container.WalkCB[E]) {
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

func (n *node[E]) walkPostOrder(cb container.WalkCB[E]) {
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

func (n *node[E]) walkPreOrder(cb container.WalkCB[E]) {
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

// Intrinsic holds nodes which are their own ordering key.
type Intrinsic[E cmp.Ordered] struct {
	root *node[E]
}

// Len returns the number of nodes in the tree, for the container.Countable interface.
// Complexity is O(n).
func (t *Intrinsic[E]) Len() int {
	l := 0
	t.WalkPostOrder(func(_ *E) { l++ })
	return l
}

func (t *Intrinsic[E]) Elements() []E {
	var sl []E
	t.WalkPreOrder(func(e *E) { sl = append(sl, *e) })
	return sl
}

// Upsert adds a value to the tree, replacing and returning the previous one if any.
// If none existed, it returns nil.
func (t *Intrinsic[E]) Upsert(e ...*E) []*E {
	results := make([]*E, 0, len(e))
	var result *E
	for _, oneE := range e {
		n := &node[E]{data: oneE}

		switch {
		case t == nil, e == nil:
			result = nil
		case t.root == nil:
			t.root = n
			result = nil
		default:
			result = t.root.upsert(n)
		}
		results = append(results, result)
	}
	return results
}

func (t *Intrinsic[E]) Delete(e *E) {
	if t == nil || e == nil {
		return
	}
	t.root.delete(e)
}

// IndexOf returns the position of the value among those in the tree.
// If the value cannot be found, it will return 0, false, otherwise the position
// starting at 0, and true.
func (t *Intrinsic[E]) IndexOf(e *E) (int, bool) {
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

// WalkInOrder is useful for search and listing nodes in order.
func (t *Intrinsic[E]) WalkInOrder(cb container.WalkCB[E]) {
	if t == nil {
		return
	}
	t.root.walkInOrder(cb)
}

// WalkPostOrder in useful for deleting subtrees.
func (t *Intrinsic[E]) WalkPostOrder(cb container.WalkCB[E]) {
	if t == nil {
		return
	}
	t.root.walkPostOrder(cb)
}

// WalkPreOrder is useful to clone the tree.
func (t *Intrinsic[E]) WalkPreOrder(cb container.WalkCB[E]) {
	if t == nil {
		return
	}
	t.root.walkPreOrder(cb)
}

func (t *Intrinsic[E]) Clone() container.BinarySearchTree[E] {
	clone := &Intrinsic[E]{}
	t.WalkPreOrder(func(e *E) {
		clone.Upsert(e)
	})
	return clone
}
