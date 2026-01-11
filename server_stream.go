package gortmp

import "github.com/dongmedia/go-rtmp/message"

type Stream struct {
	ID uint32

	AudioChan chan *message.MediaPacket
	VideoChan chan *message.MediaPacket

	AACConfig []byte
	SPS       []byte
	PPS       []byte
}

func NewStream(id uint32) *Stream {
	return &Stream{
		ID:        id,
		AudioChan: make(chan *message.MediaPacket, 1024),
		VideoChan: make(chan *message.MediaPacket, 1024),
	}
}
