package ws

import (
	"net"
)

type Server struct{}

func (s *Server) Start(conn *net.Conn, req HTTPRequest) {
	if req.Method != HTTPMethodGET {
		return
	}

	upgrade, upgradeOk := req.Header(HeaderKeyUpgrade)
	connection, connectionOk := req.Header(HeaderKeyConnection)

	if !upgradeOk || !connectionOk {
		return
	}

	if upgrade == HeaderValueWSUpgrade && connection == HeaderValueConnection_Upgrade {
	}
}
