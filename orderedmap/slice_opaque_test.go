package orderedmap_test

import (
	"strconv"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/fgm/container"
	"github.com/fgm/container/orderedmap"
)

// This ensures that OrderedMap maintains compatibility with sync.Map:
// it will fail at compile time if that compatibility is broken.
var _ = (func(orderedMap container.OrderedMap[any, any]) int { return 0 })(&sync.Map{})

type in struct {
	key   string
	value int
}

func TestSlice_Range(t *testing.T) {
	t.Parallel()
	tests := [...]struct {
		name         string
		stable       bool
		expectedKeys []string
		expectedVals []int
	}{
		{
			"stable",
			true,
			[]string{"1", "2", "3", "4", "6", "7", "8"},
			[]int{11, 2, 3, 4, 6, 7, 8},
		},
		{
			"recency-based",
			false,
			[]string{"2", "3", "4", "6", "7", "8", "1"},
			[]int{2, 3, 4, 6, 7, 8, 11},
		},
	}
	const size = 8
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			input := make([]in, size)
			for i := 1; i <= size; i++ {
				input[i-1] = in{key: strconv.Itoa(i), value: i}
			}

			var omc interface {
				container.OrderedMap[string, int]
				container.Countable
			} = orderedmap.NewSlice[string, int](size, test.stable)

			for _, e := range input {
				omc.Store(e.key, e.value)
			}
			// Ensure deletion actually removes existing entries.
			omc.Delete("5")
			if omc.Len() != size-1 {
				t.Fatalf("len is %d, expected %d", omc.Len(), size-1)
			}
			// Ensure deletion does not actually remove anything for nonexistent entries.
			omc.Delete("50")
			if omc.Len() != size-1 {
				t.Fatalf("len is %d, expected %d", omc.Len(), size-1)
			}
			// Ensure updates do not change the entries count.
			omc.Store("1", 11)
			if omc.Len() != size-1 {
				t.Fatalf("len is %d, expected %d", omc.Len(), size-1)
			}

			var keys = make([]string, 0, size)
			var vals = make([]int, 0, size)
			omc.Range(func(k string, v int) bool {
				keys = append(keys, k)
				vals = append(vals, v)
				// Return false on the expected last key to cover the break condition.
				if test.stable {
					return k != "8"
				}
				return k != "1"
			})
			if !cmp.Equal(keys, test.expectedKeys) {
				t.Fatalf("Failed keys comparison:%s", cmp.Diff(keys, test.expectedKeys))
			}
			if !cmp.Equal(vals, test.expectedVals) {
				t.Fatalf("Failed values comparison:%s", cmp.Diff(vals, test.expectedVals))
			}
		})
	}
}

func TestSlice_Store_Load_Delete(t *testing.T) {
	t.Parallel()
	const one = "one"
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
			var om = orderedmap.NewSlice[string, int](1, test.stable)
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
		})
	}
}
