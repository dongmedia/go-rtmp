package gortmp

import (
	"bytes"

	"github.com/dongmedia/go-rtmp/amf"
)

func writeOnStatusPlayStart(w *ChunkWriter, streamID uint32) error {
	var p bytes.Buffer
	_ = amf.EncodeString(&p, "onStatus")
	_ = amf.EncodeNumber(&p, 0)
	_ = amf.EncodeNull(&p)
	_ = amf.EncodeNull(&p) // 최소 (필요시 object로 확장)

	return w.WriteMessage(5, 0, 20, streamID, p.Bytes())
}

func sendSequenceHeadersIfAny(s *Stream, sub *Subscriber) {
	// 구독자에게 “먼저” 시퀀스 헤더를 보내야 디코더가 정상 시작합니다.
	if len(s.lastAACSeq) > 0 {
		_ = sub.W.WriteMessage(4, 0, 8, s.ID, s.lastAACSeq)
	}
	if len(s.lastAVCSeq) > 0 {
		_ = sub.W.WriteMessage(6, 0, 9, s.ID, s.lastAVCSeq)
	}
}
