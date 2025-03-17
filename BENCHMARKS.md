# Benchmarks

## Summary

- Slice with preallocation wins in most benchmarks, on both amd64 and arm64.
- The ranking column totals the relative ranks in each category. Lower is better.
    - Preallocating improves aggregate results in all cases, so use it if you can.
    - Although the list with internal pool comes close, it never outperforms the
      simple slice-based implementation.
    - Ranking is identical across Go 1.18 vs Go 1.24 and amd64 vs arm64.
    - 
## Data
### ARM64

go: 1.24.1
goos: darwin
goarch: amd64
pkg: github.com/fgm/container
cpu: Apple M3 Pro

| Storage         | Queue.Enqueue | Queue.Dequeue | Stack.Push | Stack.Pop |      Ranking |
|:----------------|--------------:|--------------:|-----------:|----------:|-------------:|
| Slice prealloc  |     2.0 ns/op |     1.4 ns/op |  2.0 ns/op | 1.1 ns/op |  1+1+1+2 = 5 |
| Slice raw       |     3.8 ns/op |     2.2 ns/op |  3.4 ns/op | 0.8 ns/op | 3+3+3+1 = 10 |
| ListIP prealloc |     2.2 ns/op |     2.5 ns/op |  2.0 ns/op | 2.6 ns/op | 2+4+1+5 = 12 |
| List            |    23.6 ns/op |     1.4 ns/op | 27.9 ns/op | 1.6 ns/op | 4+1+6+4 = 15 |
| ListIP raw      |    23.7 ns/op |     6.0 ns/op | 23.8 ns/op | 1.3 ns/op | 5+5+5+3 = 18 |
| ListSP prealloc |    28.1 ns/op |     7.6 ns/op |  9.4 ns/op | 7.6 ns/op | 6+6+4+6 = 22 |
| ListSP raw      |    36.5 ns/op |    13.2 ns/op | 36.4 ns/op | 7.9 ns/op | 7+7+7+7 = 28 |

### AMD64

go: 1.18.0
goos: darwin
goarch: amd64
pkg: github.com/fgm/container
cpu: Intel(R) Core(TM) i7-4980HQ CPU @ 2.80GHz

| Storage         | Queue.Enqueue | Queue.Dequeue | Stack.Push | Stack.Pop |      Ranking |
|:----------------|--------------:|--------------:|-----------:|----------:|-------------:|
| Slice prealloc  |       9 ns/op |      10 ns/op |    8 ns/op |   7 ns/op |  1+1+1+2 = 5 |
| Slice raw       |      25 ns/op |      11 ns/op |   20 ns/op |   5 ns/op | 3+3+3+1 = 10 |
| ListIP prealloc |      17 ns/op |      14 ns/op |   11 ns/op |  11 ns/op | 2+4+2+4 = 12 |
| List            |     144 ns/op |      10 ns/op |  130 ns/op |  56 ns/op | 5+2+6+6 = 19 |
| ListIP raw      |     139 ns/op |      84 ns/op |  121 ns/op |   9 ns/op | 4+7+5+3 = 19 |
| ListSP prealloc |     533 ns/op |      52 ns/op |   96 ns/op |  18 ns/op | 7+5+5+5 = 22 |
| ListSP raw      |     368 ns/op |      56 ns/op |  385 ns/op |  61 ns/op | 6+6+7+6 = 25 |

## Ranking per op

The `>>` operator indicates a big decrease in performance.

```
- Enqueue: SP > IPP > SR              >> IPR >  L  > SPR > SPP
- Dequeue: SP >  L  > SR  > IPP       >> SPP > SPR > IPR
- Push:    SP > IPP > SR              >> SPP > IPR >  L  > SPR
- Pop:     SR >  SP > IPR > IPP > SPP >>  L  > SPR
```
