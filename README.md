```
BenchmarkFastRouter_10000000Goroutines_2Params-8        10000000               530 ns/op             419 B/op          3 allocs/op
BenchmarkFastRouter_8Goroutines_2Params-8                5000000               301 ns/op             416 B/op          3 allocs/op
BenchmarkFastRouter_1Goroutine_2Params-8                 2000000               792 ns/op             416 B/op          3 allocs/op
BenchmarkFastRouter_10000000Goroutines_26Prams-8        10000000               912 ns/op             418 B/op          3 allocs/op
BenchmarkFastRouter_10000000Goroutines_NoPrams-8        10000000               484 ns/op             416 B/op          3 allocs/op

BenchmarkHttpRouter_10000000Goroutines_2Params-8        10000000               748 ns/op             487 B/op          4 allocs/op
BenchmarkHttpRouter_8Goroutines_2Params-8                5000000               289 ns/op             480 B/op          4 allocs/op
BenchmarkHttpRouter_1Goroutine_2Params-8                 2000000               787 ns/op             480 B/op          4 allocs/op
BenchmarkHttpRouter_10000000Goroutines_26Params-8       10000000              2027 ns/op            1328 B/op          4 allocs/op
BenchmarkHttpRouter_10000000Goroutines_NoParams-8       10000000               525 ns/op             416 B/op          3 allocs/op

BenchmarkFastRouter_10000000Goroutines_2Params-8        10000000              7068 ns/op            1800 B/op         14 allocs/op
BenchmarkFastRouter_8Goroutines_2Params-8                1000000              1810 ns/op            1424 B/op         13 allocs/op
BenchmarkFastRouter_1Goroutine_2Params-8                  500000              4241 ns/op            1424 B/op         13 allocs/op
BenchmarkFastRouter_10000000Goroutines_26Params-8        10000000             51068 ns/op            4027 B/op         16 allocs/op
BenchmarkFastRouter_10000000Goroutines_NoParams-8        10000000              2336 ns/op            1108 B/op         12 allocs/op
```
