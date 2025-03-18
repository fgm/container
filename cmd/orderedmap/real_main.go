package main

import (
	"fmt"
	"io"
	"strconv"

	"github.com/fgm/container/orderedmap"
)

type in struct {
	key   string
	value int
}

func realMain(w1, w2 io.Writer) int {
	const size = 8
	input := make([]in, size)
	for i := 1; i <= size; i++ {
		input[i-1] = in{key: strconv.Itoa(i), value: i}
	}
	fmt.Fprintln(w1, "Go map:")
	m := make(map[string]int, size)
	for _, e := range input {
		m[e.key] = e.value
	}
	delete(m, "5")
	m["1"] = 11
	for k, v := range m {
		fmt.Fprintln(w1, k, v)
	}

	fmt.Fprintln(w2, "OrderedMap:")
	var om = orderedmap.NewSlice[string, int](size, true)
	for _, e := range input {
		om.Store(e.key, e.value)
	}
	om.Delete("5")
	om.Store("1", 11)
	om.Range(func(k string, v int) bool {
		fmt.Fprintln(w2, k, v)
		return true
	})
	return 0
}
