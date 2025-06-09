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
BenchmarkServerGet1ReqPerConn-16                        20081194               571.9 ns/op           195 B/op          7 allocs/op
BenchmarkServerGet2ReqPerConn-16                        41866302               279.8 ns/op            91 B/op          3 allocs/op
BenchmarkServerGet10ReqPerConn-16                       150649264               71.36 ns/op           18 B/op          0 allocs/op
BenchmarkServerGet10KReqPerConn-16                      367426825               29.95 ns/op            0 B/op          0 allocs/op
BenchmarkServerGet1ReqPerConn10KClients-16              25662541               455.1 ns/op           171 B/op          7 allocs/op
BenchmarkServerGet2ReqPerConn10KClients-16              45048451               246.1 ns/op            85 B/op          3 allocs/op
BenchmarkServerGet10ReqPerConn10KClients-16             170778354               67.67 ns/op           17 B/op          0 allocs/op
BenchmarkServerGet100ReqPerConn10KClients-16            314508996               35.83 ns/op            1 B/op          0 allocs/op
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
