package http

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type Response struct {
	version       string
	statusCode    int
	statusMessage string
	server        string
	contentLength int
	body          []byte
	headers       map[string]string
	conn          net.Conn
}

func NewResponse(conn net.Conn, version string) Response {
	return Response{
		server:  "basliq labs",
		headers: make(map[string]string),
		version: version,
		conn:    conn,
	}
}

func (res *Response) WriteHeader(status StatusCode) {
	res.statusCode = int(status)
	res.statusMessage = status.String()
}

func (res *Response) Write(data []byte) {
	if res.statusCode == 0 || res.statusMessage == "" {
		res.WriteHeader(Ok)
	}
	res.body = data
	res.contentLength = len(data)
	response := res.generate()
	fmt.Println(response)
	res.conn.Write([]byte(response))
}

func (res *Response) SetHeader(key, value string) {
	res.headers[key] = value
}

func (res *Response) generate() string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%s %d %s\r\n", res.version, res.statusCode, res.statusMessage))
	builder.WriteString(fmt.Sprintf("Content-Length: %d\r\n", res.contentLength))
	builder.WriteString(fmt.Sprintf("Server: %s\r\n", res.server))
	for k, v := range res.headers {
		if k == "Content-Length" || k == "Date" {
			continue
		}
		builder.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	builder.WriteString(fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123)))
	builder.WriteString(fmt.Sprintf("\r\n%s", res.body))
	return builder.String()
}
