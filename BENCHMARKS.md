# Benchmarks

## Summary

- Slice with preallocation wins in most benchmarks.
- Preallocating improves aggregate results in all cases, so use it if you can.

## Data

goos: darwin
goarch: amd64
pkg: github.com/fgm/container
cpu: Intel(R) Core(TM) i7-4980HQ CPU @ 2.80GHz

| Storage         |    Queue.Enqueue | Queue.Dequeue | Stack.Push |   Stack.Pop |
|:----------------|-----------------:|--------------:|-----------:|------------:|
| Slice prealloc  |          9 ns/op |      10 ns/op |    8 ns/op |     7 ns/op |
| Slice raw       |         25 ns/op |      11 ns/op |   20 ns/op |     5 ns/op |
| List            |        144 ns/op |      10 ns/op |  130 ns/op |    56 ns/op |
| ListSP prealloc |        533 ns/op |      52 ns/op |   96 ns/op |    18 ns/op |
| ListSP raw      |        368 ns/op |      56 ns/op |  385 ns/op |    61 ns/op |
| ListIP prealloc |         17 ns/op |      14 ns/op |   11 ns/op |    11 ns/op |
| ListIP raw      |        139 ns/op |      84 ns/op |  121 ns/op |     9 ns/op | 

## Ranking per op

The `>>` operator indicates a big decrease in performance.

```
- Enqueue: SP > IPP > SR              >> IPR >  L  > SPR > SPP
- Dequeue: SP >  L  > SR  > IPP       >> SPP > SPR > IPR
- Push:    SP > IPP > SR              >> SPP > IPR >  L  > SPR
- Pop:     SR >  SP > IPR > IPP > SPP >>  L  > SPR
```

## Scores

Scores calculated by adding the ranking of implementations for each operation.
Less is better.

| Implementation          | Enqueue | Dequeue | Push | Pop | Total |
|-------------------------|:-------:|:-------:|:----:|:---:|------:|
| Slice, prealloc         |    1    |    1    |  1   |  2  |     5 |
| Slice, raw              |    3    |    3    |  3   |  1  |    10 |
| Internal pool, prealloc |    2    |    4    |  2   |  4  |    12 |
| List                    |    5    |    2    |  6   |  6  |    19 |
| Internal pool, raw      |    4    |    7    |  5   |  3  |    19 |
| Sync pool, prealloc     |    7    |    5    |  5   |  5  |    22 |
| Sync pool, raw          |    6    |    6    |  7   |  6  |    25 |

Although the list with internal pool comes close, it never outperforms the
simple slice-based implementation.
