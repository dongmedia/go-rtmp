package message

import "errors"

func ParseH264(payload []byte) (*H264Packet, error) {
	if len(payload) < 5 {
		return nil, errors.New("h264 payload too short")
	}

	frameType := payload[0] >> 4
	codecID := payload[0] & 0x0f

	if codecID != 7 { // AVC
		return nil, errors.New("not h264")
	}

	pktType := H264PacketType(payload[1])
	isKey := frameType == 1

	return &H264Packet{
		IsKeyFrame: isKey,
		Type:       pktType,
		Payload:    payload[5:], // skip CTS
	}, nil
}
