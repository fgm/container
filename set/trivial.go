package set

import (
	"fmt"
	"iter"
	"strings"

	"github.com/fgm/container"
)

type unit = struct{}

// Trivial is the textbook Go implementation of a Go set using generics.
//
// It is not concurrency-safe.
// For performance optimization, its union/intersection/difference operations
// may return the receiver or the argument instead of cloning.
type Trivial[E comparable] struct {
	items map[E]unit
}

// String returns an idiomatic unordered representation of the set items.
func (s *Trivial[E]) String() string {
	// Shortcut empty case.
	if s == nil || len(s.items) == 0 {
		return "{}"
	}

	var b strings.Builder
	b.WriteByte('{')

	// Use a separate counter to avoid trailing comma
	// Use a separate counter to avoid trailing comma
	i := 0
	for item := range s.items {
		if i > 0 {
			b.WriteString(", ")
		}
		_, _ = fmt.Fprintf(&b, "%v", item)
		i++
	}

	b.WriteByte('}')
	return b.String()
}

// Len returns the number of items in the Set.
func (s *Trivial[E]) Len() int {
	if s == nil {
		return 0
	}
	return len(s.items)
}

// Add adds an item to the set. Returns true if the item was already present.
func (s *Trivial[E]) Add(item E) (found bool) {
	if s == nil {
		return false
	}
	found = s.Contains(item)
	s.items[item] = unit{}
	return found
}

// Remove removes an item from the set.
// It does not fail if the item was not present, and returns true if it was.
func (s *Trivial[E]) Remove(item E) (found bool) {
	if s == nil {
		return false
	}
	found = s.Contains(item)
	delete(s.items, item)
	return found
}

// Contains returns true if the item is present in the set.
func (s *Trivial[E]) Contains(item E) bool {
	if s == nil {
		return false
	}
	_, exists := s.items[item]
	return exists
}

// Clear removes all items from the set and returns the number of items removed.
func (s *Trivial[E]) Clear() (count int) {
	if s == nil {
		return 0
	}

	count = s.Len()
	s.items = make(map[E]unit)
	return count
}

// Items returns an unordered iterator over the set'set elements.
func (s *Trivial[E]) Items() iter.Seq[E] {
	if s == nil {
		return func(yield func(E) bool) {}
	}

	return func(yield func(E) bool) {
		for item := range s.items {
			if !yield(item) {
				break
			}
		}
	}
}

// Union returns a new set containing elements present in either set.
//
// Note that it may return one of its arguments without creating a clone.
func (s *Trivial[E]) Union(other container.Set[E]) container.Set[E] {
	// Shortcut degenerate cases.
	if s == nil && other == nil {
		return NewTrivial[E](0)
	}
	if s == nil {
		return other
	}
	if other == nil {
		return s
	}
	if s.Len() == 0 {
		return other
	}
	if other, ok := other.(container.Countable); ok && other.Len() == 0 {
		return s
	}

	// Non-degenerate case. It will be at least as long as the receiver.
	result := &Trivial[E]{items: make(map[E]unit, s.Len())}

	// Add all items from this set
	for item := range s.items {
		result.Add(item)
	}

	// Add all items from other set
	for item := range other.Items() {
		result.Add(item)
	}

	return result
}

// Intersection returns a new set containing elements present in both sets.
func (s *Trivial[E]) Intersection(other container.Set[E]) container.Set[E] {
	// Shortcut degenerate cases.
	if s == nil || other == nil {
		return NewTrivial[E](0)
	}
	if s.Len() == 0 {
		return s
	}
	if countable, ok := other.(container.Countable); ok && countable.Len() == 0 {
		return other
	}

	// Non-degenerate case with size optimization
	var result *Trivial[E]
	if other, ok := other.(container.Countable); ok {
		result = &Trivial[E]{items: make(map[E]unit, min(s.Len(), other.Len()))}
	} else {
		result = &Trivial[E]{items: make(map[E]unit)}
	}

	// Add items that exist in both sets
	for item := range s.items {
		if other.Contains(item) {
			result.Add(item)
		}
	}

	return result
}

// Difference returns a new set containing elements present in this set but not in the other.
func (s *Trivial[E]) Difference(other container.Set[E]) container.Set[E] {
	// Shortcut degenerate cases.
	if s == nil {
		return NewTrivial[E](0)
	}
	if other == nil || s.Len() == 0 {
		return s
	}
	if other, ok := other.(container.Countable); ok && other.Len() == 0 {
		return s
	}

	// Non-degenerate case.
	result := &Trivial[E]{items: make(map[E]unit, s.Len())}

	// Add items that exist in this set but not in other
	for item := range s.items {
		if !other.Contains(item) {
			result.Add(item)
		}
	}

	return result
}

// NewTrivial returns a ready-for-use container.Set implemented by the Trivial type.
func NewTrivial[E comparable](sizeHint int) container.Set[E] {
	return &Trivial[E]{make(map[E]unit, sizeHint)}
}
