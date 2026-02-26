package ws

import (
	"fmt"
	"net"
)

const (
	HeaderKeyUpgrade              = "upgrade"
	HeaderKeyConnection           = "connection"
	HeaderKeyOrigin               = "origin"
	HeaderKeyWSVersion            = "sec-webSocket-version"
	HeaderKeyWSKey                = "sec-websocket-key"
	HeaderValueWSUpgrade          = "websocket"
	HeaderValueConnection_Upgrade = "Upgrade"
	HeaderValueWSVersion          = "13"
	HTTPMethodGET                 = "GET"
	GUID                          = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
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
