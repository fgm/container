# Go containers
[![github](https://github.com/fgm/container/actions/workflows/workflow.yml/badge.svg)](https://github.com/fgm/container/actions/workflows/workflow.yml)
[![codecov](https://codecov.io/gh/fgm/container/branch/main/graph/badge.svg?token=8YYX1B720M)](https://codecov.io/gh/fgm/container)
[![Go Report Card](https://goreportcard.com/badge/github.com/fgm/container)](https://goreportcard.com/report/github.com/fgm/container)

This module contains minimal type-safe Queue and Stack implementations using
Go 1.18 generics.


## Contents

See the available types by underlying storage 

| Type  | Slice | List | List+sync.Pool | List+internal pool |     Recommended      |
|:-----:|:-----:|:----:|:--------------:|:------------------:|:--------------------:|
| Queue |   Y   |  Y   |       Y        |                    | Slice with size hint |
| Stack |   Y   |  Y   |       Y        |                    | Slice with size hint |

Generally speaking, in terms of performance: 

- Slice > plain List > list+sync.Pool
- Preallocated > not preallocated

See [BENCHARKS.md](BENCHMARKS.md) for details.

## Use

See complete listing in [`cmd/example.go`](cmd/example.go)

### Queues

```go
q := queue.NewSliceQueue[Element](sizeHint) // resp. NewListQueue
q.Enqueue(e)
if lq, ok := q.(container.Countable); ok {
    fmt.Printf("elements in queue: %d\n", lq.Len())
}
for i := 0; i < 2; i++ {
    e, ok := q.Dequeue()
    fmt.Printf("Element: %v, ok: %t\n", e, ok)
}
```

### Stacks

```go
s := stack.NewSliceStack[Element](sizeHint) // resp. NewListStack
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
