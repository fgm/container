# Go containers

This module contains minimal type-safe Queue and Stack implementations using
Go 1.18 generics.


## Contents

See the available types by underlying storage 

| Container | Slice | plain List | List with pool |     Recommended      |
|:---------:|:-----:|:----------:|:--------------:|:--------------------:|
|   Queue   |   Y   |     Y      |       Y        | Slice with size hint |
|   Stack   |   Y   |     Y      |       Y        | Slice with size hint |

Generally speaking, in terms of performance: 

- Slice > plain List > list with pool
- Preallocated > not preallocated

See [BENCHARKS.md](BENCHMARKS.md) for details.


## Use

See complete listing in [`cmd/example.go`](cmd/example.go)

### Queues

```go
queue := container.NewSliceQueue[Element](sizeHint) // resp. NewListQueue
queue.Enqueue(e)
if lq, ok := queue.(container.Countable); ok {
    fmt.Printf("elements in queue: %d\n", lq.Len())
}
for i := 0; i < 2; i++ {
    e, ok := queue.Dequeue()
    fmt.Printf("Element: %v, ok: %t\n", e, ok)
}
```

### Stacks

```go
stack := container.NewSliceStack[Element](sizeHint) // resp. NewListStack
stack.Push(e)
if ls, ok := stack.(container.Countable); ok {
    fmt.Printf("elements in stack: %d\n", ls.Len())
}
for i := 0; i < 2; i++ {
    e, ok := stack.Pop()
    fmt.Printf("Element: %v, ok: %t\n", e, ok)
}
```
