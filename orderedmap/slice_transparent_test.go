package orderedmap

import (
	"errors"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSlice_mustIndexOf(t *testing.T) {
	tests := [...]struct {
		name   string
		stable bool
	}{
		{"stable", true},
		{"recency-based", false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			// Catch expected panic
			defer func() { _ = recover() }()

			var om = NewSlice[string, int](1, test.stable)
			om.store["one"] = 1
			om.mustIndexOf("one")
			t.Fatalf("mustIndexOf did not panic on missing order key")
		})
	}
}

func TestSlice_Range_inconsistent(t *testing.T) {
	tests := [...]struct {
		name   string
		stable bool
	}{
		{"stable", true},
		{"recency-based", false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			// Catch expected panic
			defer func() { _ = recover() }()

			var om = NewSlice[string, int](1, test.stable)
			om.Store("one", 1)
			delete(om.store, "one")
			om.Range(func(_ string, _ int) bool { return false })
			t.Fatalf("Range did not panic on missing map entry")
		})
	}
}

// TestSlice_Range_Mutators evaluates how Range[Mutable] behaves with callbacks that may delete keys in the map.
func TestSlice_Range_Mutators(t *testing.T) {
	t.Parallel()
	const size = 4
	type testType = Slice[int, int]
	src := func() testType {
		base := NewSlice[int, int](size, false)
		for i := range size {
			base.Store(i, i)
		}
		return *base

	}
	// transform represents a Slice[int, int] to a []int, in which the first element
	// is the int value of .stable and the values are k,v pairs.
	transform := func(s testType) []int {
		var res = make([]int, 1, 1+2*size)
		if s.stable {
			res[0] = 1
		}
		for _, e := range s.order {
			res = append(res, e, s.store[e])
		}
		return res
	}
	fnNop := func(input *testType) func(key, value int) bool {
		return func(key, value int) bool {
			return true
		}
	}
	fnDeleteOdd := func(input *testType) func(key, value int) bool {
		return func(key, value int) bool {
			if key%2 == 1 {
				input.Delete(key)
			}
			return true
		}
	}
	fnInsertOne := func(input *testType) func(key, value int) bool {
		return func(key, value int) bool {
			if key == 1 {
				input.Store(10, 10)
			}
			return true
		}
	}
	fnMutateOne := func(input *testType) func(key, value int) bool {
		return func(key, value int) bool {
			if key == 1 {
				input.Store(1, 10)
			}
			return true
		}
	}
	tests := [...]struct {
		name            string
		useRangeMutable bool
		fn              func(*testType) func(key, value int) bool
		expectedPanic   error
		expected        []int
	}{
		{"nop/plain range", false, fnNop, nil, []int{0, 0, 0, 1, 1, 2, 2, 3, 3}},
		{"nop/mutable", true, fnNop, nil, []int{0, 0, 0, 1, 1, 2, 2, 3, 3}},
		// Mutating a value moves it to the end of the map.
		{"mutateOne/plain range", false, fnMutateOne, nil, []int{0, 0, 0, 2, 2, 3, 3, 1, 10}},
		{"mutateOne/mutable", false, fnMutateOne, nil, []int{0, 0, 0, 2, 2, 3, 3, 1, 10}},
		// Store a value moves it to the end of the map.
		{"insertOne/plain range", false, fnInsertOne, nil, []int{0, 0, 0, 1, 1, 2, 2, 3, 3, 10, 10}},
		{"insertOne/mutable", false, fnInsertOne, nil, []int{0, 0, 0, 1, 1, 2, 2, 3, 3, 10, 10}},
		// This is why we need to some form of mutation support: observe the unexpected results.
		{"deleteOdd/plain range", false, fnDeleteOdd, errors.New(""), []int{0, 0, 0, 2, 2}},
		{"deleteOdd/mutable", true, fnDeleteOdd, nil, []int{0, 0, 0, 2, 2}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			input := src()
			var actualPanic error
			func() {
				defer func() {
					if err := recover(); err != nil {
						actualPanic = err.(error)
					}
				}()
				if test.useRangeMutable {
					input.RangeMutable(test.fn(&input))
				} else {
					input.Range(test.fn(&input))
				}
			}()
			actual := transform(input)
			t1, t2 := reflect.TypeOf(actualPanic), reflect.TypeOf(test.expectedPanic)
			if t1 != t2 {
				t.Fatalf("unexpected panic: %T, expected: %T", actualPanic, test.expectedPanic)
			}
			if !cmp.Equal(actual, test.expected) {
				t.Logf("expected: %d/%d %v\n                      actual:   %d/%d %v",
					len(test.expected), cap(test.expected), test.expected,
					len(actual), cap(actual), actual)
				t.Fatalf("unexpected result: %s", cmp.Diff(test.expected, actual))
			}
		})
	}
}
