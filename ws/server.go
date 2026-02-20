package ws

import (
	"crypto/sha1"
	"encoding/base64"
	"net"
)

type Server struct{}

func (s *Server) Start(conn net.Conn, req HTTPRequest) {
	if req.Method != HTTPMethodGET {
		HTTPResponse400(conn)
		return
	}

	upgrade, upgradeOk := req.Header(HeaderKeyUpgrade)
	connection, connectionOk := req.Header(HeaderKeyConnection)

	if !upgradeOk || !connectionOk {
		HTTPResponse400(conn)
		return
	}

	if upgrade == HeaderValueWSUpgrade && connection == HeaderValueConnection_Upgrade {
		clientKey, ok := req.Header(HeaderKeyWSKey)
		if !ok {
			HTTPResponse400(conn)
			return
		}

		acceptValue := calculateAcceptValue(clientKey)
		HTTPResponse101(conn, acceptValue)
		manageWSConnection(conn)
	} else {
		HTTPResponse400(conn)
	}
}

func calculateAcceptValue(clientKey string) string {
	data := clientKey + GUID
	hash := sha1.Sum([]byte(data))
	return base64.StdEncoding.EncodeToString(hash[:])
}
