package http

import (
	"fmt"
	"io"
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
	// headers holds other header fields
	headers map[string]string
	// headersSent
	headersSent bool
	// conn holds the TCP connection information.
	// required for writing a response.
	conn io.ReadWriteCloser
	// encoder
	encoder Encoder
}

// newResponse creates a new response struct that has
// useful receiver functions.
func newResponse(conn io.ReadWriteCloser, req Request) Response {
	return Response{
		server:  "basliq labs",
		headers: make(map[string]string),
		version: req.Version(),
		method:  req.Method(),
		conn:    conn,
		encoder: selectEncoder(req.Header("Accept-Encoding")),
	}
}

// SetStatus is used to set a status code and status message.
func (res *Response) SetStatus(status StatusCode) {
	res.SetStatusWithMessage(status, status.String())
}

func (res *Response) SetStatusWithMessage(status StatusCode, message string) {
	res.statusCode = int(status)
	res.statusMessage = message
}

// Write is used to write a response. This function receives an argument
// that will be written in body.
func (res *Response) Write(data []byte) {
	if !res.headersSent {
		if res.statusCode == 0 || res.statusMessage == "" {
			res.SetStatus(StatusOk)
		}
		res.contentLength = len(data)
		res.WriteHeader()
	}
	res.conn.Write(res.encoder.Encode([]byte(data)))
}

// SetHeader takes a key and value pair that will be written in the headers.
func (res *Response) SetHeader(key, value string) {
	res.headers[key] = value
}

func (res *Response) WriteHeader() {
	res.headersSent = true
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
	res.conn.Write([]byte(builder.String()))
}

func selectEncoder(value string, enable bool) Encoder {
	if !enable || value == "" {
		return PLAIN
	}

	encArr := strings.Split(value, ",")
	// TODO: implement priority list for encoders
	for _, enc := range encArr {
		e, err := IsEncodingValid(strings.TrimSpace(enc))
		if err == nil {
			return e
		}
	}
	return PLAIN
}
