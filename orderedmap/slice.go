package orderedmap

import (
	"fmt"

	"github.com/fgm/container"
)

type Slice[K comparable, V any] struct {
	order []K
	store map[K]V
}

// mustIndexOf may only be used with a key known to be present in the map.
func (s Slice[K, V]) mustIndexOf(k K) int {
	for i, ck := range s.order {
		if ck == k {
			return i
		}
	}
	// Should never happen: someone probably caused a race condition.
	panic(fmt.Errorf("structure inconsistency: key %v not found", k))
}

func (s *Slice[K, V]) Delete(k K) {
	_, loaded := s.store[k]
	if !loaded {
		return
	}
	delete(s.store, k)
	index := s.mustIndexOf(k)
	s.order = append(s.order[:index], s.order[index+1:]...)
}

func (s Slice[K, V]) Load(key K) (V, bool) {
	v, loaded := s.store[key]
	return v, loaded
}

func (s Slice[K, V]) Range(f func(key K, value V) bool) {
	for _, k := range s.order {
		v, loaded := s.store[k]
		if !loaded {
			panic(fmt.Errorf("structure inconsistency: key %v not found", k))
		}
		if !f(k, v) {
			break
		}
	}
}

func (s *Slice[K, V]) Store(k K, v V) {
	_, loaded := s.store[k]
	if !loaded {
		s.order = append(s.order, k)
		s.store[k] = v
		return
	}

	index := s.mustIndexOf(k)
	s.order = append(s.order[:index], s.order[index+1:]...)
	s.order = append(s.order, k)
	s.store[k] = v
}

func NewSlice[K comparable, V any](sizeHint int) container.OrderedMap[K, V] {
	s := &Slice[K, V]{
		order: make([]K, 0, sizeHint),
		store: make(map[K]V, sizeHint),
	}

	return s
}
