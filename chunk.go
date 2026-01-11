package gortmp

import (
	"encoding/binary"
	"io"
)

type Chunk struct {
	Fmt       uint8
	CSID      uint32
	Timestamp uint32
	TypeID    uint8
	StreamID  uint32
	Payload   []byte
}

type ChunkReader struct {
	r io.Reader
}

func NewChunkReader(r io.Reader) *ChunkReader {
	return &ChunkReader{r: r}
}

func (rd *ChunkReader) Read() (*Chunk, error) {
	// Basic Header (1 byte only: fmt=0, csid=3)
	bh := make([]byte, 1)
	if _, err := io.ReadFull(rd.r, bh); err != nil {
		return nil, err
	}

	fmt := bh[0] >> 6
	csid := uint32(bh[0] & 0x3f)

	// Message Header (fmt0 = 11 bytes)
	mh := make([]byte, 11)
	if _, err := io.ReadFull(rd.r, mh); err != nil {
		return nil, err
	}

	timestamp := uint32(mh[0])<<16 | uint32(mh[1])<<8 | uint32(mh[2])
	length := uint32(mh[3])<<16 | uint32(mh[4])<<8 | uint32(mh[5])
	typeID := mh[6]
	streamID := binary.LittleEndian.Uint32(mh[7:11])

	payload := make([]byte, length)
	if _, err := io.ReadFull(rd.r, payload); err != nil {
		return nil, err
	}

	return &Chunk{
		Fmt:       fmt,
		CSID:      csid,
		Timestamp: timestamp,
		TypeID:    typeID,
		StreamID:  streamID,
		Payload:   payload,
	}, nil
}
