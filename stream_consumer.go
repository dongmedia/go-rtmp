package gortmp

import "log"

func ConsumeStream(s *Stream) {
	go func() {
		for a := range s.AudioChan {
			log.Printf("[AUDIO] ts=%d size=%d\n", a.Timestamp, len(a.Payload))
		}
	}()

	go func() {
		for v := range s.VideoChan {
			log.Printf("[VIDEO] ts=%d size=%d\n", v.Timestamp, len(v.Payload))
		}
	}()
}
