package set_test

import (
	"fmt"
	"iter"
	"strconv"
	"strings"
	"testing"

	"github.com/fgm/container"
	"github.com/fgm/container/set"
)

type unit = struct{}

type testIntSet interface {
	container.Set[int]
	container.Countable
	fmt.Stringer
}

// nilSet provides a typed nil that can be used as a receiver.
func nilSet() testIntSet {
	var ns *set.BasicMap[int]
	return ns
}

func TestBasicMap(t *testing.T) {
	tests := []struct {
		name string
		fn   func(*testing.T)
	}{
		{"nil receiver", testNilReceiver},
		{"empty set", testEmptySet},
		{"single element", testSingleElement},
		{"multiple elements", testMultipleElements},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.fn)
	}
}

func testNilReceiver(t *testing.T) {
	s := nilSet()
	if s.Len() != 0 {
		t.Errorf("nil set should have length 0")
	}
	if s.Contains(1) {
		t.Errorf("nil set should not contain anything")
	}
	if s.Add(1) {
		t.Errorf("nil set Add should return false")
	}
	if s.Remove(1) {
		t.Errorf("nil set Remove should return false")
	}
	if s.Clear() != 0 {
		t.Errorf("nil set Clear should return 0")
	}

	count := 0
	for range s.Items() {
		count++
	}
	if count != 0 {
		t.Errorf("nil set should not iterate over any items")
	}
}

func testEmptySet(t *testing.T) {
	s := set.NewBasicMap[int](0)
	c, ok := s.(container.Countable)
	if !ok {
		t.Fatalf("expected a countable implementation")
	}
	if c.Len() != 0 {
		t.Errorf("empty set should have length 0")
	}
	if s.Contains(1) {
		t.Errorf("empty set should not contain anything")
	}
	if s.Remove(1) {
		t.Errorf("empty set Remove should return false")
	}
	if s.Clear() != 0 {
		t.Errorf("empty set Clear should return 0")
	}
}

func testSingleElement(t *testing.T) {
	s := set.NewBasicMap[int](1)
	if found := s.Add(1); found {
		t.Errorf("Add to empty set should return false")
	}
	if !s.Contains(1) {
		t.Errorf("set should contain added element")
	}
	c, ok := s.(container.Countable)
	if !ok {
		t.Fatalf("expected a countable implementation")
	}
	if c.Len() != 1 {
		t.Errorf("set should have length 1")
	}
	if found := s.Add(1); !found {
		t.Errorf("Add of existing element should return true")
	}
	if !s.Remove(1) {
		t.Errorf("Remove of existing element should return true")
	}
	if s.Remove(1) {
		t.Errorf("Remove of non-existing element should return false")
	}
}

func testMultipleElements(t *testing.T) {
	s := set.NewBasicMap[int](3)
	elements := []int{1, 2, 3}

	for _, e := range elements {
		if found := s.Add(e); found {
			t.Errorf("Add of new element should return false")
		}
	}

	c, ok := s.(container.Countable)
	if !ok {
		t.Fatalf("expected a countable implementation")
	}
	if c.Len() != 3 {
		t.Errorf("set should have length 3")
	}

	count := 0
	for item := range s.Items() {
		if !s.Contains(item) {
			t.Errorf("set should contain iterated item %v", item)
		}
		count++
	}
	if count != 3 {
		t.Errorf("should iterate over 3 items, got %d", count)
	}

	if cleared := s.Clear(); cleared != 3 {
		t.Errorf("Clear should return 3, got %d", cleared)
	}
	if c.Len() != 0 {
		t.Errorf("set should be empty after Clear")
	}
}

func TestBasicMap_SetOperations(t *testing.T) {
	tests := []struct {
		name string
		fn   func(*testing.T)
	}{
		{"union", testUnion},
		{"intersection", testIntersection},
		{"difference", testDifference},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.fn)
	}
}

func testUnion(t *testing.T) {
	cases := []struct {
		name     string
		s1       container.Set[int]
		s2       container.Set[int]
		expected int
	}{
		// nil without a type is not a valid receiver type.
		{"nilSet + nil", nilSet(), nil, 0},
		{"nilSet + nilSet", nilSet(), nilSet(), 0},
		{"nilSet + empty", nilSet(), set.NewBasicMap[int](0), 0},
		{"nilSet + empty", nilSet(), set.NewBasicMap[int](0), 0},
		{"empty + nil", set.NewBasicMap[int](0), nil, 0},
		{"empty + nilSet", set.NewBasicMap[int](0), nilSet(), 0},
		{"empty + empty", set.NewBasicMap[int](0), set.NewBasicMap[int](0), 0},
		{"nilSet + nonempty", nilSet(), createSet(t, 1, 2), 2},
		{"nonempty + nil", createSet(t, 1, 2), nil, 2},
		{"nonempty + nilSet", createSet(t, 1, 2), nilSet(), 2},
		{"disjoint", createSet(t, 1, 2), createSet(t, 3, 4), 4},
		{"overlapping", createSet(t, 1, 2), createSet(t, 2, 3), 3},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s1, s2 := tc.s1, tc.s2
			result := s1.Union(s2)
			c, ok := result.(container.Countable)
			if !ok {
				t.Fatalf("expected a countable implementation")
			}
			if c.Len() != tc.expected {
				t.Errorf("expected length %d, got %d", tc.expected, c.Len())
			}
		})
	}
}

func testIntersection(t *testing.T) {
	cases := []struct {
		name     string
		s1       container.Set[int]
		s2       container.Set[int]
		expected int
	}{
		// nil without a type is not a valid receiver type.
		{"nilSet + nil", nilSet(), nil, 0},
		{"nilSet + nilSet", nilSet(), nilSet(), 0},
		{"nilSet + empty", nilSet(), set.NewBasicMap[int](0), 0},
		{"empty + nil", set.NewBasicMap[int](0), nil, 0},
		{"empty + nilSet", set.NewBasicMap[int](0), nilSet(), 0},
		{"empty + empty", set.NewBasicMap[int](0), set.NewBasicMap[int](0), 0},
		{"nilSet + nonempty", nilSet(), createSet(t, 1, 2), 0},
		{"nonempty + nil", createSet(t, 1, 2), nil, 0},
		{"nonempty + nilSet", createSet(t, 1, 2), nilSet(), 0},
		{"disjoint", createSet(t, 1, 2), createSet(t, 3, 4), 0},
		{"overlapping", createSet(t, 1, 2), createSet(t, 2, 3), 1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.s1.Intersection(tc.s2)
			c, ok := result.(container.Countable)
			if !ok {
				t.Fatalf("expected a countable implemenration")
			}
			if c.Len() != tc.expected {
				t.Errorf("expected length %d, got %d", tc.expected, c.Len())
			}
		})
	}
}

func testDifference(t *testing.T) {
	cases := []struct {
		name     string
		s1       container.Set[int]
		s2       container.Set[int]
		expected int
	}{
		// nil without a type is not a valid receiver type.
		{"nilSet - nilSet", nilSet(), nilSet(), 0},
		{"nilSet - empty", nilSet(), set.NewBasicMap[int](0), 0},
		{"empty - nil", set.NewBasicMap[int](0), nil, 0},
		{"empty - nilSet", set.NewBasicMap[int](0), nilSet(), 0},
		{"empty - empty", set.NewBasicMap[int](0), set.NewBasicMap[int](0), 0},
		{"nonempty - nil", createSet(t, 1, 2), nil, 2},
		{"nonempty - nilSet", createSet(t, 1, 2), nilSet(), 2},
		{"nonempty - empty", createSet(t, 1, 2), set.NewBasicMap[int](0), 2},
		{"disjoint", createSet(t, 1, 2), createSet(t, 3, 4), 2},
		{"overlapping", createSet(t, 1, 2), createSet(t, 2, 3), 1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.s1.Difference(tc.s2)
			c, ok := result.(container.Countable)
			if !ok {
				t.Fatalf("expected a countable implementation")
			}
			if c.Len() != tc.expected {
				t.Errorf("expected length %d, got %d", tc.expected, c.Len())
			}
		})
	}
}

// Helper function to create a set with given elements
func createSet(tb testing.TB, elements ...int) *set.BasicMap[int] {
	tb.Helper()
	s := set.NewBasicMap[int](len(elements))
	for _, e := range elements {
		s.Add(e)
	}
	tr, ok := s.(*set.BasicMap[int])
	if !ok {
		tb.Fatalf("expected a set.BasicMap implementation")
	}
	return tr
}

func TestBasicMap_String(t *testing.T) {
	tests := []struct {
		name     string
		set      testIntSet
		elements []int
	}{
		{"nil set", nilSet(), nil},
		{"empty set", set.NewBasicMap[int](0).(testIntSet), nil},
		{"single element", createSet(t, 42), []int{42}},
		{"multiple elements", createSet(t, 1, 2, 3), []int{1, 2, 3}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := test.set.String()

			// Check empty/nil tests
			if test.elements == nil {
				if actual != "{}" {
					t.Errorf("expected empty set \"{}\", got %q", actual)
				}
				return
			}

			// For non-empty sets:
			// 1. Must start with { and end with }
			if !strings.HasPrefix(actual, "{") || !strings.HasSuffix(actual, "}") {
				t.Errorf("expected actual wrapped in {}, got %q", actual)
				return
			}

			// 2. Extract elements between { and }
			elements := actual[1 : len(actual)-1]
			if elements == "" {
				t.Errorf("unexpected empty content in non-empty set %q", actual)
				return
			}

			// 3. Split elements and trim spaces
			parts := strings.Split(elements, ",")
			nums := make([]int, 0, len(parts))
			for _, p := range parts {
				n, err := strconv.Atoi(strings.TrimSpace(p))
				if err != nil {
					t.Errorf("invalid number format in %q", p)
					return
				}
				nums = append(nums, n)
			}

			// 4. Verify we got all expected elements regardless of order
			if len(nums) != len(test.elements) {
				t.Errorf("expected %d elements, got %d", len(test.elements), len(nums))
				return
			}

			// Convert expected elements to map for easy lookup
			expected := make(map[int]bool)
			for _, e := range test.elements {
				expected[e] = true
			}

			// Verify each found element was expected
			for _, n := range nums {
				if !expected[n] {
					t.Errorf("unexpected element %d in actual", n)
				}
				delete(expected, n)
			}

			// Verify we found all expected elements
			if len(expected) > 0 {
				missing := make([]int, 0, len(expected))
				for e := range expected {
					missing = append(missing, e)
				}
				t.Errorf("missing elements: %v", missing)
			}
		})
	}
}

func TestBasicMap_ItemsEarlyTermination(t *testing.T) {
	// Create a set with multiple elements
	s := set.NewBasicMap[int](5)
	for i := 1; i <= 5; i++ {
		s.Add(i)
	}

	// Count how many items we actually iterate through
	count := 0
	for _ = range s.Items() {
		count++
		// Break after seeing 3 items by returning false from the yield function
		if count == 3 {
			break
		}
	}

	if count != 3 {
		t.Errorf("expected iteration to stop after 3 items, but got %d items", count)
	}
}

func TestBasicMap_IntersectionNonCountable(t *testing.T) {
	// Create a mock set that implements Set[int] but not Countable
	mock := &mockNonCountableSet{elements: map[int]bool{1: true, 2: true}}

	// Create a normal set with some overlapping elements
	s := set.NewBasicMap[int](3)
	s.Add(1)
	s.Add(3)

	// Test intersection with non-countable set
	result := s.Intersection(mock)
	c, ok := result.(container.Countable)
	if !ok {
		t.Fatal("expected a countable implementation")
	}

	// Verify the result
	if c.Len() != 1 {
		t.Errorf("expected intersection size 1, got %d", c.Len())
	}
	if !result.Contains(1) {
		t.Error("expected intersection to contain 1")
	}
	if result.Contains(2) || result.Contains(3) {
		t.Error("intersection contains unexpected elements")
	}
}

// mockNonCountableSet implements Set[int] but not Countable
type mockNonCountableSet struct {
	elements map[int]bool
}

func (m *mockNonCountableSet) Add(item int) bool {
	exists := m.elements[item]
	m.elements[item] = true
	return exists
}

func (m *mockNonCountableSet) Remove(item int) bool {
	return false // not needed for this test
}

func (m *mockNonCountableSet) Contains(item int) bool {
	return m.elements[item]
}

func (m *mockNonCountableSet) Clear() int {
	return 0 // not needed for this test
}

func (m *mockNonCountableSet) Items() iter.Seq[int] {
	return nil // not needed for this test
}

func (m *mockNonCountableSet) String() string {
	return "mock set" // not needed for this test
}

func (m *mockNonCountableSet) Union(other container.Set[int]) container.Set[int] {
	return nil // not needed for this test
}

func (m *mockNonCountableSet) Intersection(other container.Set[int]) container.Set[int] {
	result := &mockNonCountableSet{elements: make(map[int]bool)}
	for k := range m.elements {
		if other.Contains(k) {
			result.elements[k] = true
		}
	}
	return result
}

func (m *mockNonCountableSet) Difference(other container.Set[int]) container.Set[int] {
	return nil // not needed for this test
}

func FuzzBasicMapAdd(f *testing.F) {
	// Add some seed corpus
	f.Add(1, 2, 3)                    // Distinct values
	f.Add(0, 0, 1)                    // Some duplicate values
	f.Add(-1, -1, -1)                 // Only duplicates
	f.Add(2147483647, -2147483648, 0) // Int max/min values

	f.Fuzz(func(t *testing.T, a, b, c int) {
		s := set.NewBasicMap[int](3)

		// Test adding elements
		firstFound := s.Add(a)
		if firstFound {
			t.Errorf("First Add of %d should return false", a)
		}

		secondFound := s.Add(a)
		if !secondFound {
			t.Errorf("Second Add of %d should return true", a)
		}

		s.Add(b)
		s.Add(c)

		// Verify all elements were added correctly
		if !s.Contains(a) {
			t.Errorf("Set should contain %d", a)
		}
		if !s.Contains(b) {
			t.Errorf("Set should contain %d", b)
		}
		if !s.Contains(c) {
			t.Errorf("Set should contain %d", c)
		}

		// Verify length is correct (accounting for duplicates)
		expectedLen := len(map[int]unit{a: {}, b: {}, c: {}})
		countable := s.(container.Countable)
		if countable.Len() != expectedLen {
			t.Errorf("Expected length %d, got %d", expectedLen, countable.Len())
		}
	})
}

func FuzzBasicMapUnion(f *testing.F) {
	// Add some seed corpus
	f.Add(1, 2, 3, 4)     // All distinct
	f.Add(-1, -1, -1, -1) // All duplicates
	f.Add(0, 1, 0, 2)     // Intermixed duplicates

	f.Fuzz(func(t *testing.T, a, b, c, d int) {
		// Create two sets with potentially overlapping elements
		s1 := set.NewBasicMap[int](2)
		s1.Add(a)
		s1.Add(b)

		s2 := set.NewBasicMap[int](2)
		s2.Add(c)
		s2.Add(d)

		// Perform union operation
		result := s1.Union(s2)

		// Verify the result contains all elements from both sets
		if !result.Contains(a) {
			t.Errorf("Union should contain %d from first set", a)
		}
		if !result.Contains(b) {
			t.Errorf("Union should contain %d from first set", b)
		}
		if !result.Contains(c) {
			t.Errorf("Union should contain %d from second set", c)
		}
		if !result.Contains(d) {
			t.Errorf("Union should contain %d from second set", d)
		}

		// Verify size is correct (accounting for duplicates)
		expectedLen := len(map[int]unit{a: {}, b: {}, c: {}, d: {}})
		countable := result.(container.Countable)
		if countable.Len() != expectedLen {
			t.Errorf("Expected union length %d, got %d", expectedLen, countable.Len())
		}
	})
}

func FuzzBasicMapItems(f *testing.F) {
	// Add some seed corpus.
	f.Add(1, 2, 3)                    // Distinct values.
	f.Add(0, 0, 1)                    // Some duplicate values.
	f.Add(-1, -1, -1)                 // Only duplicates.
	f.Add(2147483647, -2147483648, 0) // Int max/min values.

	f.Fuzz(func(t *testing.T, a, b, c int) {
		// Create a set with the fuzzed elements
		s := set.NewBasicMap[int](3)
		s.Add(a)
		s.Add(b)
		s.Add(c)

		// Create a map of expected elements (accounting for duplicates).
		expected := map[int]unit{a: {}, b: {}, c: {}}

		// Collect all items from the iterator.
		actual := make(map[int]unit)
		for item := range s.Items() {
			actual[item] = unit{}
		}

		// Verify that the actual items match the expected elements.
		if len(actual) != len(expected) {
			t.Errorf("Expected %d unique elements, got %d", len(expected), len(actual))
		}

		// Verify each actual element was expected
		for item := range actual {
			if _, ok := expected[item]; !ok {
				t.Errorf("Iterator yielded unexpected element: %v", item)
			}
		}

		// Verify each expected element was actual
		for item := range expected {
			if _, ok := actual[item]; !ok {
				t.Errorf("Iterator failed to yield expected element: %v", item)
			}
		}

		// Test early termination
		count := 0
		for range s.Items() {
			count++
			break // Stop after first item
		}

		// Verify early termination worked correctly
		if count != 1 && len(expected) > 0 {
			t.Errorf("Early termination: Expected 1 item before breaking, got %d", count)
		}
	})
}
