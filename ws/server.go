package ws

import (
	"crypto/sha1"
	"encoding/base64"
	"net"
)

type Server struct {
	// loggingEnabled indicates if logging is active.
	loggingEnabled bool

	// previousLogger holds the logger used before disabling logging.
	previousLogger Logger

	// Logger is the active logger for server messages.
	Logger Logger
}

func New() *Server {
	return &Server{
		Logger: &DefaultLogger{},
	}
}

func (s *Server) Start(conn net.Conn, req HTTPRequest) {
	if req.Method != HTTPMethodGET {
		s.Logger.Error("method mismatch", nil)
		HTTPResponse400(conn)
		return
	}

	upgrade, upgradeOk := req.Header(HeaderKeyUpgrade)
	connection, connectionOk := req.Header(HeaderKeyConnection)

	if !upgradeOk || !connectionOk {
		s.Logger.Error("headers do not comply with RFC 6455", nil)
		HTTPResponse400(conn)
		return
	}

	if upgrade == HeaderValueWSUpgrade && connection == HeaderValueConnection_Upgrade {
		clientKey, ok := req.Header(HeaderKeyWSKey)
		if !ok {
			s.Logger.Error("headers do not comply with RFC 6455", nil)
			HTTPResponse400(conn)
			return
		}

		acceptValue := calculateAcceptValue(clientKey)
		s.Logger.Success("websocket connection initialized")
		HTTPResponse101(conn, acceptValue)
		s.manageWSConnection(conn)
	} else {
		s.Logger.Error("error while establishing websocket connection", nil)
		HTTPResponse400(conn)
	}
}

func calculateAcceptValue(clientKey string) string {
	data := clientKey + GUID
	hash := sha1.Sum([]byte(data))
	return base64.StdEncoding.EncodeToString(hash[:])
}
