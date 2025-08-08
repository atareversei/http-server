package http

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// Response is the struct that gets populated during the
// response phase when a TCP request hits the server.
type Response struct {
	method Method
	// version is the HTTP version that the request has
	// been created in.
	version Version
	// statusCode is of HTTP response codes.
	statusCode int
	// statusMessage if of HTTP response messages.
	statusMessage string
	// server is a header field that indicates
	// the server name.
	server string
	// contentLength is a header field that indicates
	// the length of body in bytes if there is any
	contentLength int
	// body holds the body data
	body []byte
	// headers holds other header fields
	headers map[string]string
	// conn holds the TCP connection information.
	// required for writing a response.
	conn net.Conn
}

// newResponse creates a new response struct that has
// useful receiver functions.
func newResponse(conn net.Conn, version Version, method Method) Response {
	return Response{
		server:  "basliq labs",
		headers: make(map[string]string),
		version: version,
		method:  method,
		conn:    conn,
	}
}

// WriteHeader is used to set a status code and status message.
func (res *Response) WriteHeader(status StatusCode) {
	res.statusCode = int(status)
	res.statusMessage = status.String()
}

// Write is used to write a response. This function receives an argument
// that will be written in body.
func (res *Response) Write(data []byte) {
	if res.statusCode == 0 || res.statusMessage == "" {
		res.WriteHeader(StatusOk)
	}
	res.body = data
	res.contentLength = len(data)
	response := res.generate()
	res.conn.Write([]byte(response))
}

// SetHeader takes a key and value pair that will be written in the headers.
func (res *Response) SetHeader(key, value string) {
	res.headers[key] = value
}

// generate is a function that will create the final Response in the
// right format.
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
	if res.method != MethodHead {
		return builder.String()
	}
	builder.WriteString(fmt.Sprintf("\r\n%s", res.body))
	return builder.String()
}
