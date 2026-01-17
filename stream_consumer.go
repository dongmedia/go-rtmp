package gortmp

import (
	"fmt"

	"github.com/dongmedia/go-rtmp/message"
)

func ConsumeStream(s *Stream) {
	go func() {
		for a := range s.AudioChan {
			pkt, err := message.ParseAAC(a.Payload)
			if err != nil {
				select {
				case s.ErrChan <- fmt.Errorf("parse AAC: %w", err):
				default:
				}
				continue
			}

			if pkt.Type == message.AACSequenceHeader {
				s.AACConfig = pkt.Payload
			}
		}
	}()

	go func() {
		for v := range s.VideoChan {
			pkt, err := message.ParseH264(v.Payload)
			if err != nil {
				select {
				case s.ErrChan <- fmt.Errorf("parse H264: %w", err):
				default:
				}
				continue
			}

			if pkt.Type == message.H264SequenceHeader {
				if err := parseAVCConfig(pkt.Payload, s); err != nil {
					s.ErrChan <- fmt.Errorf("parse avc config err: %v", err)
				}
			}
		}
	}()
}
