package message

type AACPacketType uint8

const (
	AACSequenceHeader AACPacketType = 0
	AACRaw            AACPacketType = 1
)

type H264PacketType uint8

const (
	H264SequenceHeader H264PacketType = 0
	H264NALU           H264PacketType = 1
)

type AACPacket struct {
	Type    AACPacketType
	Payload []byte
}

type H264Packet struct {
	IsKeyFrame bool
	Type       H264PacketType
	Payload    []byte
}
