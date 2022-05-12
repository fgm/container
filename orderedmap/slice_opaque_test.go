package orderedmap

import (
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type in struct {
	key   string
	value int
}

func TestSlice_Range(t *testing.T) {
	const size = 8
	input := make([]in, size)
	for i := 1; i <= size; i++ {
		input[i-1] = in{key: strconv.Itoa(i), value: i}
	}
	expectedKeys := []string{"2", "3", "4", "6", "7", "8", "1"}
	expectedVals := []int{2, 3, 4, 6, 7, 8, 11}

	var om = NewSlice[string, int](size)
	for _, e := range input {
		om.Store(e.key, e.value)
	}
	om.Delete("5")
	om.Store("1", 11)

	var keys = make([]string, 0, size)
	var vals = make([]int, 0, size)
	om.Range(func(k string, v int) bool {
		keys = append(keys, k)
		vals = append(vals, v)
		return k != "1" // This is the last key because we added it after Delete.
	})
	if !cmp.Equal(keys, expectedKeys) {
		t.Fatalf("Failed keys comparison:%s", cmp.Diff(keys, expectedKeys))
	}
	if !cmp.Equal(vals, expectedVals) {
		t.Fatalf("Failed values comparison:%s", cmp.Diff(vals, expectedVals))
	}
}

func TestSlice_Store_Load_Delete(t *testing.T) {
	const one = "one"
	var om = NewSlice[string, int](1)
	om.Store(one, 1)
	zero, loaded := om.Load("zero")
	if loaded {
		t.Fatalf("unexpected load success for missing key %s, value is %v", "zero", zero)
	}
	_, loaded = om.Load(one)
	if !loaded {
		t.Fatal("unexpected load failure for present key")
	}
	om.Delete(one)
	om.Delete(one) // Ensure multiple deletes do not cause an error
	actual, loaded := om.Load(one)
	if loaded {
		t.Fatalf("unexpected load success for missing key %s, value is %v", one, actual)
	}
}
