package gortmp

import (
	"net"

	"github.com/dongmedia/go-rtmp/message"
)

type Conn struct {
	conn      net.Conn
	handshake HandshakeService
}

func NewConn(c net.Conn) *Conn {
	return &Conn{conn: c}
}

func (c *Conn) Serve() {
	defer c.conn.Close()

	if err := c.handshake.Do(c.conn); err != nil {
		return
	}

	rd := NewChunkReader(c.conn)
	var streamID uint32 = 1

	for {
		ch, err := rd.Read()
		if err != nil {
			return
		}

		if ch.TypeID != 20 {
			continue
		}

		cmd, err := message.DecodeCommand(ch.Payload)
		if err != nil {
			continue
		}

		switch cmd.Name {

		case "connect":
			writeConnectSuccess(c.conn, cmd.TransactionID)

		case "createStream":
			writeCreateStreamResult(c.conn, cmd.TransactionID, streamID)

		case "publish":
			writePublishStart(c.conn, streamID)
		}
	}
}
