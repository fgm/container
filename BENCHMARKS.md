Summary:

- Slice with preallocation wins in all benchmarks.
- List wins no benchmark 

goos: darwin
goarch: amd64
pkg: github.com/fgm/container
cpu: Intel(R) Core(TM) i5-1038NG7 CPU @ 2.00GHz

BenchmarkSliceQueue_Dequeue_prealloc-8   	403397954	         3.014 ns/op
BenchmarkSliceQueue_Dequeue_raw-8        	393923754	         3.118 ns/op
BenchmarkListQueue_Dequeue-8             	419303047	         3.208 ns/op

BenchmarkSliceQueue_Enqueue_prealloc-8   	373755117	         3.106 ns/op
BenchmarkSliceQueue_Enqueue_raw-8        	188607312	        25.11  ns/op
BenchmarkListQueue_Enqueue-8             	29348926	        43.42  ns/op

BenchmarkSliceStack_Pop_prealloc-8       	494476706	         2.369 ns/op
BenchmarkSliceStack_Pop_raw-8            	479408811	        10.78  ns/op
BenchmarkListStack_Pop-8                 	380633768	         4.632 ns/op

BenchmarkSliceStack_Push_prealloc-8      	363288776	         3.067 ns/op
BenchmarkSliceStack_Push_raw-8           	85524145	        13.67  ns/op
BenchmarkListStack_Push-8                	29843348	        40.64  ns/op
