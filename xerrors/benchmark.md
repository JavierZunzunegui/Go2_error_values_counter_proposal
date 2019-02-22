# Benchmark results

`go test --bench=Benchmark -benchmem ./...`

goos: darwin
goarch: amd64
pkg: github.com/JavierZunzunegui/Go2_error_values_counter_proposal/xerrors

| name | repetitions | time/op | heap bytes/op | heap allocations/op |
| --- | --- | --- | --- | --- |
| BenchmarkString/nonWrapped-8          |   20000000	|   108 ns/op   |   3 B/op	    |   1 allocs/op |
| BenchmarkString/singleWrapped-8       |   10000000	|   202 ns/op	|   32 B/op	    |   1 allocs/op |
| BenchmarkString/doubleWrapped-8       |   5000000	    |   281 ns/op	|   48 B/op	    |   1 allocs/op |
| BenchmarkDetailString/nonWrapped-8    |  	20000000    |   107 ns/op	|   3 B/op	    |   1 allocs/op |
| BenchmarkDetailString/singleWrapped-8 |	1000000	    |   1043 ns/op	|   256 B/op	|   3 allocs/op |
| BenchmarkDetailString/doubleWrapped-8 | 	1000000	    |   2048 ns/op	|   496 B/op    |   5 allocs/op |
