package gortmp

import (
	"log"
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

	chunkReader := NewChunkReader(c.conn)

	for {
		ch, err := chunkReader.Read()
		if err != nil {
			return
		}

		if ch.TypeID == 20 { // command AMF0
			cmd, _ := message.DecodeCommand(ch.Payload)
			log.Println("CMD:", cmd.Name)

			if cmd.Name == "connect" {
				writeConnectSuccess(c.conn, cmd.TransactionID)
			}
		}
	}
}
