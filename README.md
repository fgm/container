# Go containers

[![GoDoc](https://pkg.go.dev/badge/github.com/fgm/container)](https://pkg.go.dev/github.com/fgm/container)
[![Go Report Card](https://goreportcard.com/badge/github.com/fgm/container)](https://goreportcard.com/report/github.com/fgm/container)
[![github](https://github.com/fgm/container/actions/workflows/workflow.yml/badge.svg)](https://github.com/fgm/container/actions/workflows/workflow.yml)
[![codecov](https://codecov.io/gh/fgm/container/branch/main/graph/badge.svg?token=8YYX1B720M)](https://codecov.io/gh/fgm/container)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/fgm/container/badge)](https://securityscorecards.dev/viewer/?uri=github.com/fgm/container)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/10245/badge)](https://www.bestpractices.dev/projects/10245)

This module contains minimal type-safe Ordered Map, Queue, Set and Stack implementations
using Go generics.

The Ordered Map supports both stable (in-place) updates and recency-based ordering,
making it suitable both for highest performance (in-place), and for LRU caches (recency).

## Contents

See the available types by underlying storage

| Type          | Slice | Map | List | List+sync.Pool | List+int. pool | Recommended          |
|---------------|:-----:|:---:|:----:|:--------------:|:--------------:|----------------------|
| OrderedMap    |   Y   |     |      |                |                | Slice with size hint |
| Queue         |   Y   |     |  Y   |       Y        |       Y        | Slice with size hint |
| WaitableQueue |   Y   |     |      |                |                | Slice with size hint |
| Set           |       |  Y  |      |                |                | Map with size hint   |
| Stack         |   Y   |     |  Y   |       Y        |       Y        | Slice with size hint |


**CAVEAT**: In order to optimize performance, except for WaitableQueue,
all of these implementations are unsafe for concurrent execution,
so they need protection in concurrency situations.

WaitableQueue being designed for concurrent code, on the other hand, is concurrency-safe.

Generally speaking, in terms of performance:

- Slice > list+internal pool > plain List > list+sync.Pool
- Preallocated > not preallocated

See [BENCHARKS.md](BENCHMARKS.md) for details.

## Usage

See complete listings in:

- [`cmd/orderedmap`](cmd/orderedmap/real_main.go)
- [`cmd/queuestack`](cmd/queuestack/real_main.go)
- [`cmd/set`](cmd/set/real_main.go)

### Ordered Map

```go
stable := true
om := orderedmap.NewSlice[Key, Value](sizeHint, stable) // OrderedMap and Countable
om.Store(k, v)
om.Range(func (k K, v V) bool { fmt.Println(k, v); return true })
v, loaded := om.Load(k)
if !loaded {
        fmt.Fprintf(w, "No entry for key %v\n", k)
}
om.Delete(k) // Idempotent: does not fail on nonexistent keys.
```

### Classic Queues without flow control

```go
var e Element
q := queue.NewSliceQueue[Element](sizeHint)
q.Enqueue(e)
if lq, ok := q.(container.Countable); ok {
        fmt.Fprintf(w, "elements in queue: %d\n", lq.Len())
}
for i := 0; i < 2; i++ {
        e, ok := q.Dequeue()
        fmt.Fprintf(w, "Element: %v, ok: %t\n", e, ok)
}
```

### WaitableQueue: a concurrent queue with flow control

```go
var e Element
q, _ := queue.NewWaitableQueue[Element](sizeHint, lowWatermark, highWatermark)
go func() {
        wqs := q.Enqueue(e)
        if lq, ok := q.(container.Countable); ok {
                fmt.Fprintf(w, "elements in queue: %d, status: %s\n", lq.Len(), wqs)
        }
}
<-q.WaitChan()                    // Wait for elements to be available to dequeue
for i := 0; i < 2; i++ {          // Then dequeue them
        e, ok, wqs := q.Dequeue() // Non-blocking, ok will be true for the first and false for the second 
        fmt.Fprintf(w, "Element: %v, ok: %t, status: %s\n", e, ok, wqs)
}
q.Close() // Only needed if consumers may still be waiting on <-q.WaitChan
```

### Sets

```go
var e Element
s := set.NewBasicMap[Element](sizeHint)
s.Add(e)
s.Add(e)
if cs, ok := q.(container.Countable); ok {
        fmt.Fprintf(w, "elements in set: %d\n", cs.Len()) // 1
}
for e := range s.Items() {
        fmt.Fprintln(w, e)
}

```

### Stacks

```go
s := stack.NewSliceStack[Element](sizeHint)
s.Push(e)
if ls, ok := s.(container.Countable); ok {
        fmt.Printf("elements in stack: %d\n", ls.Len())
}
for i := 0; i < 2; i++ {
        e, ok := s.Pop()
        fmt.Printf("Element: %v, ok: %t\n", e, ok)
}
```

### Running tests
#### Normal tests: unit and benchmarks

The complete test coverage requires running not only the unit tests, but also
the benchmarks, like:

```
    go test -race -run=. -bench=. -coverprofile=cover.out -covermode=atomic ./...
```

This will also run the fuzz tests in unit test mode, without triggering the fuzzing logic.


#### Fuzz tests

Fuzz tests are not run by CI, but you can run them on-demand during development with:

```
    go test -run='^$' -fuzz='^\QFuzzBasicMapAdd\E$'   -fuzztime=20s ./set
    go test -run='^$' -fuzz='^\QFuzzBasicMapItems\E$' -fuzztime=20s ./set
    go test -run='^$' -fuzz='^\QFuzzBasicMapUnion\E$' -fuzztime=20s ./set
``` 

- Adjust `-fuzztime` duration as relevant: 20 seconds is just a smoke test.
- Be sure to escape the `^$` and `\Q\E` characters in the `-fuzz` argument in your shell.
