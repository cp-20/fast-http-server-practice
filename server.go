package fast_http_server_practice

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"sync"
	"time"
)

var (
	methodGet             = []byte("GET")
	methodPost            = []byte("POST")
	responseOK            = []byte("HTTP/1.1 200 OK\r\n")
	lineBreak             = []byte("\r\n")
	headerContentLength   = []byte("Content-Length: ")
	headerContentType     = []byte("Content-Type: ")
	headerConnectionClose = []byte("Connection: close")
)

type requestContextState struct {
	resMain           *[]byte
	resHeadersBuilder *bytes.Buffer
	resBody           *[]byte
	method            *[]byte
	closeConn         bool
}

type RequestCtx struct {
	state *requestContextState
}

func (c *RequestCtx) IsGet() bool {
	return bytes.Equal(*c.state.method, methodGet)
}

func (c *RequestCtx) IsPost() bool {
	return bytes.Equal(*c.state.method, methodPost)
}

func (c *RequestCtx) Method() string {
	return string(*c.state.method)
}

func (c *RequestCtx) Success(contentType string, body []byte) {
	*c.state.resMain = responseOK
	c.state.resHeadersBuilder.Write(headerContentType)
	c.state.resHeadersBuilder.WriteString(contentType)
	c.state.resHeadersBuilder.Write(lineBreak)
	c.state.resHeadersBuilder.Write(headerContentLength)
	c.state.resHeadersBuilder.WriteString(strconv.Itoa(len(body)))
	c.state.resHeadersBuilder.Write(lineBreak)
	*c.state.resBody = body
}

func (c *RequestCtx) SetConnectionClose() {
	c.state.resHeadersBuilder.Write(headerConnectionClose)
	c.state.resHeadersBuilder.Write(lineBreak)
	c.state.closeConn = true
}

type RequestHandler func(ctx *RequestCtx)

type Server struct {
	Handler     RequestHandler
	Concurrency int
}

func handleReadError(err error) error {
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return nil
	} else if err == io.EOF {
		return nil
	}
	panic(err)
}

func getMethod(firstLine []byte) []byte {
	if len(firstLine) < 4 {
		return methodGet
	}
	if firstLine[3] == ' ' {
		return firstLine[:3]
	}
	if firstLine[4] == ' ' {
		return firstLine[:4]
	}
	if firstLine[5] == ' ' {
		return firstLine[:5]
	}
	if firstLine[6] == ' ' {
		return firstLine[:6]
	}
	if firstLine[7] == ' ' {
		return firstLine[:7]
	}
	if firstLine[8] == ' ' {
		return firstLine[:8]
	}
	return firstLine[:9]
}

var readerPool = sync.Pool{
	New: func() interface{} {
		return bufio.NewReader(nil)
	},
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func handleConn(conn net.Conn, handler RequestHandler) error {
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	reader := readerPool.Get().(*bufio.Reader)
	reader.Reset(conn)
	defer readerPool.Put(reader)

	// buffer for response
	resBuf := bufferPool.Get().(*bytes.Buffer)
	resBuf.Reset()
	defer bufferPool.Put(resBuf)

	// buffer for response headers
	resHeaderBuilder := bufferPool.Get().(*bytes.Buffer)
	defer bufferPool.Put(resHeaderBuilder)

	method := make([]byte, 0, 8)

	var rcs requestContextState
	reqCtx := RequestCtx{
		state: &rcs,
	}

	var resMain []byte
	var resBody []byte
	rcs.resMain = &resMain
	rcs.resHeadersBuilder = resHeaderBuilder
	rcs.resBody = &resBody
	rcs.method = &method

	for {
		resHeaderBuilder.Reset()

		firstLine, err := reader.ReadSlice('\n')
		if err != nil {
			return handleReadError(err)
		}

		method = append(method[:0], getMethod(firstLine)...)

		// headers := make([][]byte, 0, 3)
		contentLength := 0
		for {
			line, err := reader.ReadSlice('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				return handleReadError(err)
			}
			if len(line) == 2 && line[0] == '\r' && line[1] == '\n' {
				break
			}
			// headers = append(headers, line)
			if bytes.HasPrefix(line, headerContentLength) {
				// Parse Content-Length header
				// "Content-Length: <slice here>\r\n"
				val := line[16 : len(line)-2]
				length := 0
				for _, b := range val {
					if b == ' ' || b == '\t' {
						continue
					}
					if b < '0' || b > '9' {
						return fmt.Errorf("invalid Content-Length header: %s", val)
					}
					length = length*10 + int(b-'0')
				}
				contentLength = length
			}
		}

		var body []byte
		if contentLength > 0 {
			body = make([]byte, contentLength)
			_, err := io.ReadFull(reader, body)
			if err != nil {
				return handleReadError(err)
			}
		}

		handler(&reqCtx)

		resBuf.Reset()
		resBuf.Write(resMain)
		resBuf.Write(lineBreak)
		resBuf.Write(reqCtx.state.resHeadersBuilder.Bytes())
		resBuf.Write(lineBreak)
		resBuf.Write(resBody)
		_, err = conn.Write(resBuf.Bytes())
		if err != nil {
			panic(err)
		}

		if reqCtx.state.closeConn {
			break
		}
	}

	return nil
}

func (s *Server) Serve(ln net.Listener) error {

	for {
		conn, err := ln.Accept()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			panic(err)
		}

		go handleConn(conn, s.Handler)
	}
}

func init() {
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()
}
