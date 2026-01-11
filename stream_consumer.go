package gortmp

import (
	"log"

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
				log.Println("[AAC] sequence header", len(pkt.Payload))
			} else {
				log.Println("[AAC] frame", len(pkt.Payload))
			}
		}
	}()

	go func() {
		for v := range s.VideoChan {
			pkt, err := message.ParseH264(v.Payload)
			if err != nil {
				continue
			}

			if pkt.Type == message.H264SequenceHeader {
				parseAVCConfig(pkt.Payload, s)
				log.Println("[H264] SPS/PPS")
			} else {
				log.Printf("[H264] frame key=%v size=%d\n",
					pkt.IsKeyFrame, len(pkt.Payload))
			}
		}
	}()
}
