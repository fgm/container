package orderedmap

import (
	"testing"
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

			var om = NewSlice[string, int](1, test.stable).(*Slice[string, int])
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

			var om = NewSlice[string, int](1, test.stable).(*Slice[string, int])
			om.Store("one", 1)
			delete(om.store, "one")
			om.Range(func(_ string, _ int) bool { return false })
			t.Fatalf("Range did not panic on missing map entry")
		})
	}
}
