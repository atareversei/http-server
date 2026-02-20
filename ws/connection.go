package ws

import (
	"io"
	"net"
)

func (s *Server) manageWSConnection(conn net.Conn) {
	buffer := make([]byte, 4096)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				s.Logger.Info("connection closed by peer")
				return
			}

			s.Logger.Error("read error", err)
		}

		if n > 0 {
			processData(buffer[:n])
		}
	}
}

func processData(data []byte) {}
