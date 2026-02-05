package ws

import "net"

type Server struct{}

func (s *Server) Start(conn *net.Conn) {}
