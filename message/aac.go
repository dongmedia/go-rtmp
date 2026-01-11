package message

import "errors"

func ParseAAC(payload []byte) (*AACPacket, error) {
	if len(payload) < 2 {
		return nil, errors.New("aac payload too short")
	}

	soundFormat := payload[0] >> 4
	if soundFormat != 10 { // AAC
		return nil, errors.New("not aac")
	}

	pktType := AACPacketType(payload[1])

	return &AACPacket{
		Type:    pktType,
		Payload: payload[2:],
	}, nil
}
