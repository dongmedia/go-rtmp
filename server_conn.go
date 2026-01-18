package gortmp

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/dongmedia/go-rtmp/message"
)

type Conn struct {
	conn      net.Conn
	stream    *Stream
	handshake HandshakeService
}

func NewConn(c net.Conn) *Conn {
	return &Conn{
		conn:      c,
		handshake: NewHandshakeService(),
	}
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
				log.Printf("[RTMP] decode command err: %v", err)
				continue
			}

			log.Printf("[RTMP] command: %s tx=%d args=%v", cmd.Name, cmd.TransactionID, cmd.Args)

			switch cmd.Name {

			case "connect":
				if err := writeWindowAckSize(c.conn, 2500000); err != nil {
					return fmt.Errorf("write window ack size err: %v", err)
				}
				if err := writeSetPeerBandwidth(c.conn, 2500000, 2); err != nil {
					return fmt.Errorf("write set peer bandwidth err: %v", err)
				}
				if err := writeStreamBegin(c.conn, 0); err != nil {
					return fmt.Errorf("write stream begin err: %v", err)
				}
				if err := writeConnectSuccess(c.conn, cmd.TransactionID); err != nil {
					return fmt.Errorf("write connect success err: %v", err)
				}

			case "releaseStream", "FCPublish":
				if err := writeCommandResult(c.conn, cmd.TransactionID); err != nil {
					return fmt.Errorf("write command result err: %v", err)
				}

			case "createStream":
				c.stream = NewStream(streamID)
				if err := writeCreateStreamResult(c.conn, cmd.TransactionID, streamID); err != nil {
					return fmt.Errorf("write create stream result err: %v", err)
				}

			case "publish":
				if c.stream != nil {
					ConsumeStream(c.stream)
				}
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
