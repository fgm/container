# Go containers

[![GoDoc](https://pkg.go.dev/badge/github.com/fgm/container)](https://pkg.go.dev/github.com/fgm/container)
[![Go Report Card](https://goreportcard.com/badge/github.com/fgm/container)](https://goreportcard.com/report/github.com/fgm/container)
[![github](https://github.com/fgm/container/actions/workflows/workflow.yml/badge.svg)](https://github.com/fgm/container/actions/workflows/workflow.yml)
[![codecov](https://codecov.io/gh/fgm/container/branch/main/graph/badge.svg?token=8YYX1B720M)](https://codecov.io/gh/fgm/container)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/fgm/container/badge)](https://securityscorecards.dev/viewer/?uri=github.com/fgm/container)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/10245/badge)](https://www.bestpractices.dev/projects/10245)

This module contains minimal type-safe Double-ended queue, Ordered Map, Queue, Set and Stack implementations
using Go generics, as well as a concurrent-safe WaitableQueue.

The Ordered Map supports both stable (in-place) updates and recency-based ordering,
making it suitable both for highest performance (in-place), and for LRU caches (recency).

The double-ended queue is currently a fork of the Go stdlib [container/list](https://pkg.go.dev/container/list) package,
replacing the `any` value with a type parameter for type safety.
This API is not stable, as it differs from the other packages,
and will be deprecated when a simpler one is found, more in line with the other packages.

## Contents

See the available types by underlying storage

| Type               | Slice | Map | List | List+sync.Pool | List+int. pool | Recommended           |
|--------------------|:-----:|:---:|:----:|:--------------:|:--------------:|-----------------------|
| Double-ended queue |       |     |  Y   |                |                | Type-safe stdlib fork |
| OrderedMap         |   Y   |     |      |                |                | Slice with size hint  |
| Queue              |   Y   |     |  Y   |       Y        |       Y        | Slice with size hint  |
| WaitableQueue      |   Y   |     |      |                |                | Slice with size hint  |
| Set                |       |  Y  |      |                |                | Map with size hint    |
| Stack              |   Y   |     |  Y   |       Y        |       Y        | Slice with size hint  |

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

- [`cmd/list`](cmd/list/real_main.go)
- [`cmd/orderedmap`](cmd/orderedmap/real_main.go)
- [`cmd/queuestack`](cmd/queuestack/real_main.go)
- [`cmd/set`](cmd/set/real_main.go)
- [`cmd/waitablequeue`](cmd/waitablequeue/real_main.go)

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

### Development

Since this is a library, it has no install process, but you can build it to ensure correctness, with:

```
    make lint      # Run staticcheck linting checks
    make build     # Build the library. Include test and fuzz-smoke targets.
    make           # Shortcut for make lint && make build
```

#### Running normal tests: unit and benchmarks

For the simple version, without generating coverage reports, run:
```
    make test       # Unit tests, fast
    make coverage   # Coverage report in cover.out. A bit longer.
    make bench      # Benchmarks: run to update BENCHMARKS.md
    make fuzz-smoke # Fuzz tests as smoke tests: 10 seconds only, for CI builds
```

This will also run the fuzz tests in unit test mode, 
without triggering the fuzzing logic.


#### Running fuzz tests

Fuzz tests are not run by CI, but you can run them on-demand during development with:

```
    go test -fuzz='\QFuzzBasicMapAdd\E'   -fuzztime=20s ./set
    go test -fuzz='\QFuzzBasicMapItems\E' -fuzztime=20s ./set
    go test -fuzz='\QFuzzBasicMapUnion\E' -fuzztime=20s ./set
``` 

- Adjust `-fuzztime` duration as relevant: 20 seconds is just a smoke test.
- Be sure to escape the `\Q\E` characters in the `-fuzz` argument in your shell.

## Licensing

- In the directory [container](container) and below,
  this project includes code derived from the Go standard library's container/list package,
  which is licensed under the BSD 3-Clause License.  
  The original copyright and license text are included in the source files.
- The rest of this project is licensed under the Apache License 2.0.
