# Go containers

[![GoDoc](https://pkg.go.dev/badge/github.com/fgm/container)](https://pkg.go.dev/github.com/fgm/container)
[![Go Report Card](https://goreportcard.com/badge/github.com/fgm/container)](https://goreportcard.com/report/github.com/fgm/container)
[![github](https://github.com/fgm/container/actions/workflows/workflow.yml/badge.svg)](https://github.com/fgm/container/actions/workflows/workflow.yml)
[![codecov](https://codecov.io/gh/fgm/container/branch/main/graph/badge.svg?token=8YYX1B720M)](https://codecov.io/gh/fgm/container)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/fgm/container/badge)](https://securityscorecards.dev/viewer/?uri=github.com/fgm/container)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/10245/badge)](https://www.bestpractices.dev/projects/10245)

This module contains minimal type-safe Ordered Map, Queue and Stack implementations
using Go generics.

The Ordered Map supports both stable (in-place) updates and recency-based ordering,
making it suitable both for highest performance (in-place), and for LRU caches (recency).

## Contents

See the available types by underlying storage

| Type       | Slice | List | List+sync.Pool | List+int. pool | Recommended          |
|------------|:-----:|:----:|:--------------:|:--------------:|----------------------|
| OrderedMap |   Y   |      |                |                | Slice with size hint |
| Queue      |   Y   |  Y   |       Y        |       Y        | Slice with size hint |
| Stack      |   Y   |  Y   |       Y        |       Y        | Slice with size hint |

**CAVEAT**: In order to optimize performance, 
all of these implementations are unsafe for concurrent execution,
so they need protection in concurrency situations.

Generally speaking, in terms of performance:

- Slice > list+internal pool > plain List > list+sync.Pool
- Preallocated > not preallocated

See [BENCHARKS.md](BENCHMARKS.md) for details.

## Usage

See complete listings in:

- [`cmd/orderedmap/example.go`](cmd/orderedmap/example.go)
- [`cmd/queuestack/example.go`](cmd/queuestack/example.go)

### Ordered Map

```go
stable := true
om := orderedmap.NewSlice[Key,Value](sizeHint, stable)
om.Store(k, v)
om.Range(func(k K, v V) bool { fmt.Println(k, v); return true })
v, loaded := om.Load(k)
if !loaded { fmt.Printf("No entry for key %v\n", k)}
om.Delete(k) // Idempotent: does not fail on nonexistent keys.
```

### Queues

```go
var e Element
q := queue.NewSliceQueue[Element](sizeHint)
q.Enqueue(e)
if lq, ok := q.(container.Countable); ok {
    fmt.Printf("elements in queue: %d\n", lq.Len())
}
for i := 0; i < 2; i++ {
    e, ok := q.Dequeue()
    fmt.Printf("Element: %v, ok: %t\n", e, ok)
}
```

### Sets

```go
var e Element
s := set.NewTrivial[Element](sizeHint)
s.Add(e)
s.Add(e)
if cs, ok := q.(container.Countable); ok {
    fmt.Printf("elements in set: %d\n", cs.Len()) // 1
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

The complete test coverage requires running not only the unit tests, but also
the benchmarks, like:
```
    go test -race -run=. -bench=. -coverprofile=cover.out -covermode=atomic ./...
```
