# fast http server (for practice)

acknowledge for https://github.com/valyala/fasthttp

## Benchmark

### My framework

```
$ task bench-fast
goos: linux
goarch: amd64
pkg: github.com/cp-20/fast-http-server-practice
cpu: AMD Ryzen 9 6900HS with Radeon Graphics        
BenchmarkServerGet1ReqPerConn-16                         5672046              1996 ns/op             170 B/op          7 allocs/op
BenchmarkServerGet2ReqPerConn-16                        12485984              1056 ns/op              85 B/op          3 allocs/op
BenchmarkServerGet10ReqPerConn-16                       43329940               252.6 ns/op            17 B/op          0 allocs/op
BenchmarkServerGet10KReqPerConn-16                      206215317               51.26 ns/op            0 B/op          0 allocs/op
BenchmarkServerGet1ReqPerConn10KClients-16              11965674               947.5 ns/op           173 B/op          7 allocs/op
BenchmarkServerGet2ReqPerConn10KClients-16              25569140               475.2 ns/op            85 B/op          3 allocs/op
BenchmarkServerGet10ReqPerConn10KClients-16             100000000              116.5 ns/op            17 B/op          0 allocs/op
BenchmarkServerGet100ReqPerConn10KClients-16            249859675               44.87 ns/op            1 B/op          0 allocs/op
```

### net/http

```
$ task bench-net
goos: linux
goarch: amd64
pkg: github.com/cp-20/fast-http-server-practice
cpu: AMD Ryzen 9 6900HS with Radeon Graphics        
BenchmarkNetHTTPServerGet1ReqPerConn-16                  1898720              5455 ns/op            3346 B/op         36 allocs/op
BenchmarkNetHTTPServerGet2ReqPerConn-16                  3003394              4511 ns/op            2895 B/op         28 allocs/op
BenchmarkNetHTTPServerGet10ReqPerConn-16                 2829384              3632 ns/op            2540 B/op         23 allocs/op
BenchmarkNetHTTPServerGet10KReqPerConn-16                4662751              2561 ns/op            2405 B/op         21 allocs/op
BenchmarkNetHTTPServerGet1ReqPerConn10KClients-16        2410336              4694 ns/op            3394 B/op         36 allocs/op
BenchmarkNetHTTPServerGet2ReqPerConn10KClients-16        3363481              3273 ns/op            2904 B/op         28 allocs/op
BenchmarkNetHTTPServerGet10ReqPerConn10KClients-16       3928376              2956 ns/op            2585 B/op         23 allocs/op
BenchmarkNetHTTPServerGet100ReqPerConn10KClients-16      6754627              1688 ns/op            2430 B/op         21 allocs/op
```
