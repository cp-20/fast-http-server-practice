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
BenchmarkServerGet1ReqPerConn-16                         4172176              2915 ns/op            4369 B/op          9 allocs/op
BenchmarkServerGet2ReqPerConn-16                         7574250              1705 ns/op            2184 B/op          4 allocs/op
BenchmarkServerGet10ReqPerConn-16                       28177276               429.4 ns/op           436 B/op          0 allocs/op
BenchmarkServerGet10KReqPerConn-16                      216451035               50.87 ns/op            0 B/op          0 allocs/op
BenchmarkServerGet1ReqPerConn10KClients-16               4647902              2328 ns/op            4368 B/op          9 allocs/op
BenchmarkServerGet2ReqPerConn10KClients-16              12352939              1153 ns/op            2183 B/op          4 allocs/op
BenchmarkServerGet10ReqPerConn10KClients-16             50485514               226.2 ns/op           436 B/op          0 allocs/op
BenchmarkServerGet100ReqPerConn10KClients-16            190233699               70.12 ns/op           43 B/op          0 allocs/op
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
