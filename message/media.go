package message

type MediaType uint8

const (
	MediaAudio MediaType = 8
	MediaVideo MediaType = 9
)

type MediaPacket struct {
	Type      MediaType
	Timestamp uint32
	StreamID  uint32
	Payload   []byte
}
