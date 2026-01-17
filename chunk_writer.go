package gortmp

import (
	"encoding/binary"
	"io"
)

type ChunkWriter struct {
	w        io.Writer
	outChunk uint32
}

func NewChunkWriter(w io.Writer) *ChunkWriter {
	return &ChunkWriter{w: w, outChunk: 4096}
}

func (wr *ChunkWriter) SetOutChunkSize(sz uint32) {
	if sz > 0 {
		wr.outChunk = sz
	}
}

func (wr *ChunkWriter) WriteMessage(csid uint32, timestamp uint32, typeID uint8, streamID uint32, payload []byte) error {
	// fmt0 basic header (csid < 64 only, 최소 구현: csid 3/4/5 정도만 쓸 것)
	bh := byte(csid & 0x3f) // fmt0 => upper 2 bits 00
	if _, err := wr.w.Write([]byte{bh}); err != nil {
		return err
	}

	// message header 11 bytes
	mh := make([]byte, 11)
	mh[0] = byte(timestamp >> 16)
	mh[1] = byte(timestamp >> 8)
	mh[2] = byte(timestamp)
	mh[3] = byte(len(payload) >> 16)
	mh[4] = byte(len(payload) >> 8)
	mh[5] = byte(len(payload))
	mh[6] = typeID
	binary.LittleEndian.PutUint32(mh[7:11], streamID)

	if _, err := wr.w.Write(mh); err != nil {
		return err
	}

	// chunking
	remain := payload
	for len(remain) > 0 {
		n := int(wr.outChunk)
		if len(remain) < n {
			n = len(remain)
		}
		if _, err := wr.w.Write(remain[:n]); err != nil {
			return err
		}
		remain = remain[n:]
		if len(remain) > 0 {
			// continuation chunk: fmt3 (11)
			// fmt3 => 0b11xxxxxx
			if _, err := wr.w.Write([]byte{0xC0 | byte(csid&0x3f)}); err != nil {
				return err
			}
		}
	}
	return nil
}
