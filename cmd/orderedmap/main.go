package main

import (
	"fmt"
	"strconv"

	"github.com/fgm/container/orderedmap"
)

type in struct {
	key   string
	value int
}

func main() {
	const size = 8
	input := make([]in, size)
	for i := 1; i <= size; i++ {
		input[i-1] = in{key: strconv.Itoa(i), value: i}
	}

	fmt.Println("Go map:")
	m := make(map[string]int, size)
	for _, e := range input {
		m[e.key] = e.value
	}
	delete(m, "5")
	m["1"] = 11
	for k, v := range m {
		fmt.Println(k, v)
	}

	fmt.Println("OrderedMap:")
	var om = orderedmap.NewSlice[string, int](size)
	for _, e := range input {
		om.Store(e.key, e.value)
	}
	om.Delete("5")
	om.Store("1", 11)

	om.Range(func(k string, v int) bool {
		fmt.Println(k, v)
		return true
	})
}
