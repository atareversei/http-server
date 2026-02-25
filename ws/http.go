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

// header returns the headers of the request.
func (req *HTTPRequest) header(key string) (string, bool) {
	h, ok := req.Headers[key]
	return h, ok
}

type httpToolkit struct{}

func (ht *httpToolkit) response101(conn net.Conn, acceptValue string) {
	response := fmt.Sprintf(
		"HTTP/1.1 101 Switching Protocols\r\n"+
			"Upgrade: websocket\r\n"+
			"Connection: Upgrade\r\n"+
			"Sec-WebSocket-Accept: %s\r\n"+
			"Content-Length: 0\r\n"+
			"\r\n",
		acceptValue,
	)

	conn.Write([]byte(response))
}

func (ht *httpToolkit) response400(conn net.Conn) {
	response := fmt.Sprintf(
		"HTTP/1.1 400 Bad Request\r\n" +
			"Content-Length: 0" +
			"\r\n",
	)

	conn.Write([]byte(response))
	conn.Close()
}
