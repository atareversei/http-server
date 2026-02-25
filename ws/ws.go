package ws

import (
	"crypto/sha1"
	"encoding/base64"
	"net"
)

type Connection struct {
	// loggingEnabled indicates if logging is active.
	loggingEnabled bool

	// previousLogger holds the logger used before disabling logging.
	previousLogger Logger

	// Logger is the active logger for server messages.
	Logger Logger

	// parser is used to parse the incoming WebSocket frames.
	parser *parser

	// http provides useful utility functions for HTTP protocol requests and responses.
	http httpToolkit
}

func New() *Connection {
	return &Connection{
		Logger: &DefaultLogger{},
		parser: newParser(),
		http:   httpToolkit{},
	}
}

func (c *Connection) Start(conn net.Conn, req HTTPRequest) {
	if req.Method != HTTPMethodGET {
		c.Logger.Error("method mismatch", nil)
		c.http.response400(conn)
		return
	}

	upgrade, upgradeOk := req.header(HeaderKeyUpgrade)
	connection, connectionOk := req.header(HeaderKeyConnection)

	if !upgradeOk || !connectionOk {
		c.Logger.Error("headers do not comply with RFC 6455", nil)
		c.http.response400(conn)
		return
	}

	if upgrade == HeaderValueWSUpgrade && connection == HeaderValueConnection_Upgrade {
		clientKey, ok := req.header(HeaderKeyWSKey)
		if !ok {
			c.Logger.Error("headers do not comply with RFC 6455", nil)
			c.http.response400(conn)
			return
		}

		acceptValue := calculateAcceptValue(clientKey)
		c.Logger.Success("websocket connection initialized")
		c.http.response101(conn, acceptValue)
		c.manageWSConnection(conn)
	} else {
		c.Logger.Error("error while establishing websocket connection", nil)
		c.http.response400(conn)
	}
}

func calculateAcceptValue(clientKey string) string {
	data := clientKey + GUID
	hash := sha1.Sum([]byte(data))
	return base64.StdEncoding.EncodeToString(hash[:])
}
