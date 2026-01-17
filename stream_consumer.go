package gortmp

import (
	"github.com/dongmedia/go-rtmp/message"
)

func ConsumeStream(s *Stream) {
	go func() {
		for a := range s.AudioChan {
			pkt, err := message.ParseAAC(a.Payload)
			if err != nil {
				continue
			}

			if pkt.Type == message.AACSequenceHeader {
				s.AACConfig = pkt.Payload
				s.lastAACSeq = a.Payload // 원본 payload 저장
			}

			// 구독자에게 RTMP audio payload 그대로 전달
			s.Broadcast(8, a.Timestamp, a.Payload)
		}
	}()

	go func() {
		for v := range s.VideoChan {
			pkt, err := message.ParseH264(v.Payload)
			if err != nil {
				continue
			}

			if pkt.Type == message.H264SequenceHeader {
				_ = parseAVCConfig(pkt.Payload, s)
				s.lastAVCSeq = v.Payload
			}

			// 구독자에게 RTMP video payload 그대로 전달
			s.Broadcast(9, v.Timestamp, v.Payload)
		}
	}()
}
