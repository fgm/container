Summary:

- Slice with preallocation wins in all benchmarks.
- List2 (with pool) wins no benchmark

goos: darwin
goarch: amd64
pkg: github.com/fgm/container
cpu: Intel(R) Core(TM) i5-1038NG7 CPU @ 2.00GHz

| Storage         | Queue.Enqueue | Queue.Dequeue |   Stack.Push |    Stack.Pop |
|:----------------|--------------:|--------------:|-------------:|-------------:|
| Slice prealloc  |   3.149 ns/op |   3.404 ns/op |  2.455 ns/op |  3.301 ns/op |
| Slice raw-8     |  18.92  ns/op |  66.00  ns/op | 14.70  ns/op | 58.58  ns/op |
| List            |  26.68  ns/op |  42.84  ns/op |  4.184 ns/op | 41.42  ns/op |
| List2 prealloc  |  52.11  ns/op |  61.49  ns/op | 20.39  ns/op | 18.11  ns/op |
| List2 raw       |  74.28  ns/op |  68.60  ns/op | 29.97  ns/op | 65.25  ns/op |
