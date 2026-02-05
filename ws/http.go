package ws

import (
	"fmt"
	"net"
)

type HTTPRequest struct {
	Path    string
	Method  string
	Headers map[string]string
	Params  map[string]string
}

// Header returns the headers of the request.
func (req *HTTPRequest) Header(key string) (string, bool) {
	h, ok := req.Headers[key]
	return h, ok
}

func HTTPResponse(conn net.Conn) {
	// TODO: calculate accept key
	acceptKey := ""

	response := fmt.Sprintf(
		"HTTP/1.1 101 Switching Protocols\r\n"+
			"Upgrade: websocket\r\n"+
			"Connection: Upgrade\r\n"+
			"Sec-WebSocket-Accept: %s\r\n"+
			"\r\n",
		acceptKey,
	)

	conn.Write([]byte(response))
}
