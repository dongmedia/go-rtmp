package gortmp

import (
	"net"

	"github.com/dongmedia/go-rtmp/message"
)

type Conn struct {
	conn      net.Conn
	rd        *ChunkReader
	wr        *ChunkWriter
	reg       *Registry
	stream    *Stream
	handshake HandshakeService
}

func NewConn(c net.Conn, reg *Registry) *Conn {
	return &Conn{conn: c, reg: reg}
}

func (c *Conn) Serve() {
	defer c.conn.Close()

	if err := c.handshake.Do(c.conn); err != nil {
		return
	}

	c.rd = NewChunkReader(c.conn)
	c.wr = NewChunkWriter(c.conn)

	// OBS 안정성을 위해 기본 제어 메시지 송신
	_ = writeWindowAckSize(c.wr, 5000000)
	_ = writeSetPeerBandwidth(c.wr, 5000000, 2)
	_ = writeSetChunkSize(c.wr, 4096)
	c.wr.SetOutChunkSize(4096)

	var streamID uint32 = 1

	for {
		ch, err := c.rd.ReadMessage()
		if err != nil {
			return
		}

		switch ch.TypeID {

		case 1: // Set Chunk Size (client -> server)
			if len(ch.Payload) >= 4 {
				sz := uint32(ch.Payload[0])<<24 | uint32(ch.Payload[1])<<16 | uint32(ch.Payload[2])<<8 | uint32(ch.Payload[3])
				c.rd.SetInChunkSize(sz)
			}

		case 20: // Command
			cmd, err := message.DecodeCommand(ch.Payload)
			if err != nil {
				continue
			}

			switch cmd.Name {

			case "connect":
				// writeConnectSuccess(c.conn, cmd.TransactionID)
				writeConnectSuccess(c.wr, cmd.TransactionID)

			case "createStream":
				writeCreateStreamResult(c.wr, cmd.TransactionID, streamID)
				// c.stream = NewStream(streamID)
				// ConsumeStream(c.stream)
				// writeCreateStreamResult(c.conn, cmd.TransactionID, streamID)

			case "publish":
				name := extractStreamName(cmd.Args)
				if name == "" {
					name = "default"
				}

				s := NewStream(streamID, name)
				c.stream = s
				c.reg.Upsert(name, s)
				ConsumeStream(s)

			case "play":
				name := extractStreamName(cmd.Args)
				if name == "" {
					name = "default"
				}

				s := c.reg.Get(name)
				if s == nil {
					// 없는 스트림 재생 요청: 실제로는 onStatus(NetStream.Play.StreamNotFound) 등을 보내야 합니다.
					continue
				}

				sub := &Subscriber{W: c.wr}
				s.AddSubscriber(sub)

				_ = writeOnStatusPlayStart(c.wr, s.ID)
				sendSequenceHeadersIfAny(s, sub)
				// writePublishStart(c.conn, streamID)
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

func extractStreamName(args []any) string {
	// 흔한 형태: [null, "streamName", "live"]
	for _, a := range args {
		if s, ok := a.(string); ok && s != "" {
			return s
		}
	}
	return ""
}
