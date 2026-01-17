package gortmp

import (
	"fmt"
	"io"
	"net"

	"github.com/dongmedia/go-rtmp/message"
)

type Conn struct {
	conn      net.Conn
	stream    *Stream
	handshake HandshakeService
}

func NewConn(c net.Conn) *Conn {
	return &Conn{conn: c}
}

func (c *Conn) Serve() error {
	defer c.conn.Close()

	if err := c.handshake.Do(c.conn); err != nil {
		return fmt.Errorf("serve handshake err: %v", err)
	}

	rd := NewChunkReader(c.conn)
	var streamID uint32 = 1

	for {
		ch, err := rd.Read()
		if err != nil {
			if err == io.EOF {
				return nil // clean disconnect
			}
			return fmt.Errorf("read chunk err: %v", err)
		}

		switch ch.TypeID {

		case 20: // Command
			cmd, err := message.DecodeCommand(ch.Payload)
			if err != nil {
				continue
			}

			switch cmd.Name {

			case "connect":
				writeConnectSuccess(c.conn, cmd.TransactionID)

			case "createStream":
				c.stream = NewStream(streamID)
				ConsumeStream(c.stream)
				writeCreateStreamResult(c.conn, cmd.TransactionID, streamID)

			case "publish":
				writePublishStart(c.conn, streamID)
			}

		case 8, 9: // Audio / Video
			if c.stream == nil {
				continue
			}

			pkt := chunkToMediaPacket(ch)

			if pkt.Type == message.MediaAudio {
				c.stream.AudioChan <- pkt
			} else {
				c.stream.VideoChan <- pkt
			}
		}
	}
}
