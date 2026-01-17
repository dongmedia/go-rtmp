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
		// Check for stream consumer errors (non-blocking)
		if c.stream != nil {
			select {
			case err := <-c.stream.ErrChan:
				return fmt.Errorf("stream consumer err: %w", err)
			default:
			}
		}

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
				if err := writeConnectSuccess(c.conn, cmd.TransactionID); err != nil {
					return fmt.Errorf("write connect success err: %v", err)
				}

			case "createStream":
				c.stream = NewStream(streamID)
				ConsumeStream(c.stream)
				if err := writeCreateStreamResult(c.conn, cmd.TransactionID, streamID); err != nil {
					return fmt.Errorf("write create stream result err: %v", err)
				}

			case "publish":
				if err := writePublishStart(c.conn, streamID); err != nil {
					return fmt.Errorf("write publish start err: %v", err)
				}
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
