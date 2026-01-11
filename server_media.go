package gortmp

import "github.com/dongmedia/go-rtmp/message"

func chunkToMediaPacket(ch *Chunk) *message.MediaPacket {
	return &message.MediaPacket{
		Type:      message.MediaType(ch.TypeID),
		Timestamp: ch.Timestamp,
		StreamID:  ch.StreamID,
		Payload:   ch.Payload,
	}
}
