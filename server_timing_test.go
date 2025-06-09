package fast_http_server_practice

// ref: https://github.com/valyala/fasthttp/blob/master/server_timing_test.go

import (
	"io"
	"net"
	"net/http"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var defaultClientsCount = runtime.NumCPU() / 2

func BenchmarkServerGet1ReqPerConn(b *testing.B) {
	benchmarkServerGet(b, defaultClientsCount, 1)
}

func BenchmarkServerGet2ReqPerConn(b *testing.B) {
	benchmarkServerGet(b, defaultClientsCount, 2)
}

func BenchmarkServerGet10ReqPerConn(b *testing.B) {
	benchmarkServerGet(b, defaultClientsCount, 10)
}

func BenchmarkServerGet10KReqPerConn(b *testing.B) {
	benchmarkServerGet(b, defaultClientsCount, 10000)
}

func BenchmarkNetHTTPServerGet1ReqPerConn(b *testing.B) {
	benchmarkNetHTTPServerGet(b, defaultClientsCount, 1)
}

func BenchmarkNetHTTPServerGet2ReqPerConn(b *testing.B) {
	benchmarkNetHTTPServerGet(b, defaultClientsCount, 2)
}

func BenchmarkNetHTTPServerGet10ReqPerConn(b *testing.B) {
	benchmarkNetHTTPServerGet(b, defaultClientsCount, 10)
}

func BenchmarkNetHTTPServerGet10KReqPerConn(b *testing.B) {
	benchmarkNetHTTPServerGet(b, defaultClientsCount, 10000)
}

func BenchmarkServerGet1ReqPerConn10KClients(b *testing.B) {
	benchmarkServerGet(b, 10000, 1)
}

func BenchmarkServerGet2ReqPerConn10KClients(b *testing.B) {
	benchmarkServerGet(b, 10000, 2)
}

func BenchmarkServerGet10ReqPerConn10KClients(b *testing.B) {
	benchmarkServerGet(b, 10000, 10)
}

func BenchmarkServerGet100ReqPerConn10KClients(b *testing.B) {
	benchmarkServerGet(b, 10000, 100)
}

func BenchmarkNetHTTPServerGet1ReqPerConn10KClients(b *testing.B) {
	benchmarkNetHTTPServerGet(b, 10000, 1)
}

func BenchmarkNetHTTPServerGet2ReqPerConn10KClients(b *testing.B) {
	benchmarkNetHTTPServerGet(b, 10000, 2)
}

func BenchmarkNetHTTPServerGet10ReqPerConn10KClients(b *testing.B) {
	benchmarkNetHTTPServerGet(b, 10000, 10)
}

func BenchmarkNetHTTPServerGet100ReqPerConn10KClients(b *testing.B) {
	benchmarkNetHTTPServerGet(b, 10000, 100)
}

type fakeServerConn struct {
	net.TCPConn
	ln            *fakeListener
	requestsCount int
	pos           int
	closed        uint32
}

func (c *fakeServerConn) Read(b []byte) (int, error) {
	nn := 0
	reqLen := len(c.ln.request)
	for len(b) > 0 {
		if c.requestsCount == 0 {
			if nn == 0 {
				return 0, io.EOF
			}
			return nn, nil
		}
		pos := c.pos % reqLen
		n := copy(b, c.ln.request[pos:])
		b = b[n:]
		nn += n
		c.pos += n
		if n+pos == reqLen {
			c.requestsCount--
		}
	}
	return nn, nil
}

func (c *fakeServerConn) Write(b []byte) (int, error) {
	return len(b), nil
}

var fakeAddr = net.TCPAddr{
	IP:   []byte{1, 2, 3, 4},
	Port: 12345,
}

func (c *fakeServerConn) RemoteAddr() net.Addr {
	return &fakeAddr
}

func (c *fakeServerConn) Close() error {
	if atomic.AddUint32(&c.closed, 1) == 1 {
		// c.ln.ch <- c
		c.ln.activeConns.Add(-1)
		c.ln.connPool.Put(c)
	}
	return nil
}

func (c *fakeServerConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *fakeServerConn) SetWriteDeadline(t time.Time) error {
	return nil
}

type fakeListener struct {
	activeConns     atomic.Int64
	connPool        sync.Pool
	done            chan struct{}
	request         []byte
	requestsCount   int
	requestsPerConn int
	lock            sync.Mutex
	closed          bool
}

func (ln *fakeListener) Accept() (net.Conn, error) {
	ln.lock.Lock()
	if ln.requestsCount == 0 {
		ln.lock.Unlock()
		for ln.activeConns.Load() > 0 {
			time.Sleep(10 * time.Millisecond)
		}
		ln.lock.Lock()
		if !ln.closed {
			close(ln.done)
			ln.closed = true
		}
		ln.lock.Unlock()
		return nil, io.EOF
	}
	requestsCount := min(ln.requestsPerConn, ln.requestsCount)
	ln.requestsCount -= requestsCount
	ln.lock.Unlock()

	c := ln.connPool.Get().(*fakeServerConn)
	ln.activeConns.Add(1)
	c.requestsCount = requestsCount
	c.closed = 0
	c.pos = 0

	return c, nil
}

func (ln *fakeListener) Close() error {
	return nil
}

func (ln *fakeListener) Addr() net.Addr {
	return &fakeAddr
}

func newFakeListener(requestsCount, clientsCount, requestsPerConn int, request string) *fakeListener {
	ln := &fakeListener{
		requestsCount:   requestsCount,
		requestsPerConn: requestsPerConn,
		request:         []byte(request),
		done:            make(chan struct{}),
	}
	ln.connPool.New = func() interface{} {
		return &fakeServerConn{
			ln: ln,
		}
	}
	// for i := 0; i < clientsCount; i++ {
	// 	_ = ln.connPool.Get().(*fakeServerConn)
	// }
	return ln
}

var (
	fakeResponse = []byte("Hello, world!")
	getRequest   = "GET /foobar?baz HTTP/1.1\r\nHost: google.com\r\nUser-Agent: aaa/bbb/ccc/ddd/eee Firefox Chrome MSIE Opera\r\n" +
		"Referer: http://example.com/aaa?bbb=ccc\r\nCookie: foo=bar; baz=baraz; aa=aakslsdweriwereowriewroire\r\n\r\n"
)

func benchmarkServerGet(b *testing.B, clientsCount, requestsPerConn int) {
	var servedRequests atomic.Int64
	s := &Server{
		Handler: func(ctx *RequestCtx) {
			if !ctx.IsGet() {
				b.Fatalf("Unexpected request method: %q", ctx.Method())
			}
			ctx.Success("text/plain", fakeResponse)
			if requestsPerConn == 1 {
				ctx.SetConnectionClose()
			}
			servedRequests.Add(1)
		},
		Concurrency: 16 * clientsCount,
	}
	benchmarkServer(b, s, clientsCount, requestsPerConn, getRequest)
	verifyRequestsServed(b, &servedRequests)
}

var GET = "GET"

func benchmarkNetHTTPServerGet(b *testing.B, clientsCount, requestsPerConn int) {
	var servedRequests atomic.Int64
	s := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.Method != GET {
				b.Fatalf("Unexpected request method: %q", req.Method)
			}
			h := w.Header()
			h.Set("Content-Type", "text/plain")
			if requestsPerConn == 1 {
				h.Set("Connection", "close")
			}
			w.Write(fakeResponse)
			servedRequests.Add(1)
		}),
	}
	benchmarkServer(b, s, clientsCount, requestsPerConn, getRequest)
	verifyRequestsServed(b, &servedRequests)
}

func verifyRequestsServed(b *testing.B, servedRequests *atomic.Int64) {
	requestsServed := int(servedRequests.Load())

	requestsSent := b.N
	for requestsServed < requestsSent {
		b.Fatalf("Unexpected number of requests served %d. Expected %d", requestsServed, requestsSent)
	}
}

type realServer interface {
	Serve(ln net.Listener) error
}

func benchmarkServer(b *testing.B, s realServer, clientsCount, requestsPerConn int, request string) {
	ln := newFakeListener(b.N, clientsCount, requestsPerConn, request)
	ch := make(chan struct{})
	go func() {
		s.Serve(ln)
		ch <- struct{}{}
	}()

	<-ln.done

	select {
	case <-ch:
	case <-time.After(10 * time.Second):
		b.Fatalf("Server.Serve() didn't stop")
	}
}
