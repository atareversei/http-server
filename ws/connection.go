package ws

import (
	"io"
	"net"
)

func (c *Connection) manageWSConnection(conn net.Conn) {
	for {
		n, err := conn.Read(c.parser.buffer)

		if n > 0 {
			c.parser.parse(c.parser.buffer[:n])
		}

		if err != nil {
			if err == io.EOF {
				c.Logger.Info("connection closed by peer")
				return
			}

			c.Logger.Error("read error", err)
		}
	}
}
