package ws

import (
	"io"
	"net"
)

func (c *Connection) manageWSConnection(conn net.Conn) {
	for {
		n, err := conn.Read(c.parser.buffer.GetWriteBuffer())
		c.parser.buffer.CommitWrite(n)

		for n > 0 && !c.parser.buffer.IsEmpty() {
			msg, pErr := c.parser.parse()

			c.handlers.onMessage(msg)
			if pErr != nil {
				c.Logger.Error("parse error", err)
			}
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
