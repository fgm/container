package orderedmap

import (
	"testing"
)

func TestSlice_mustIndexOf(t *testing.T) {
	// Catch expected panic
	defer func() { _ = recover() }()

	var om = NewSlice[string, int](1).(*Slice[string, int])
	om.store["one"] = 1
	om.mustIndexOf("one")
	t.Fatalf("mustIndexOf did not panic on missing order key")
}

func TestSlice_Range_inconsistent(t *testing.T) {
	// Catch expected panic
	defer func() { _ = recover() }()

	var om = NewSlice[string, int](1).(*Slice[string, int])
	om.Store("one", 1)
	delete(om.store, "one")
	om.Range(func(_ string, _ int) bool { return false })
	t.Fatalf("Range did not panic on missing map entry")
}
