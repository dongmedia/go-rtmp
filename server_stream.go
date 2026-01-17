package gortmp

import (
	"sync"

	"github.com/dongmedia/go-rtmp/message"
)

type Subscriber struct {
	W *ChunkWriter
}

type Stream struct {
	ID   uint32
	Name string

	AudioChan chan *message.MediaPacket
	VideoChan chan *message.MediaPacket

	AACConfig []byte
	SPS       []byte
	PPS       []byte

	mu          sync.RWMutex
	subscribers []*Subscriber

	// 시퀀스 헤더 원본 RTMP payload 캐시 (구독자 최초 전송용)
	lastAACSeq []byte // type8 payload 전체
	lastAVCSeq []byte // type9 payload 전체
}

func NewStream(id uint32, name string) *Stream {
	return &Stream{
		ID:        id,
		Name:      name,
		AudioChan: make(chan *message.MediaPacket, 1024),
		VideoChan: make(chan *message.MediaPacket, 1024),
	}
}

func (s *Stream) AddSubscriber(sub *Subscriber) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.subscribers = append(s.subscribers, sub)
}

func (s *Stream) Broadcast(typeID uint8, ts uint32, payload []byte) {
	s.mu.RLock()
	subs := append([]*Subscriber(nil), s.subscribers...)
	s.mu.RUnlock()

	for _, sub := range subs {
		// audio/video는 통상 csid 4(audio), 6(video)를 많이 사용합니다(강제는 아님).
		csid := uint32(4)
		if typeID == 9 {
			csid = 6
		}
		_ = sub.W.WriteMessage(csid, ts, typeID, s.ID, payload)
	}
}
