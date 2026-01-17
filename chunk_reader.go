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
	r         io.Reader
	inChunkSz uint32

	// per-csid state
	prev map[uint32]*ChunkHeader
	buf  map[uint32][]byte
}

type ChunkHeader struct {
	Timestamp uint32
	Length    uint32
	TypeID    uint8
	StreamID  uint32
}

func NewChunkReader(r io.Reader) *ChunkReader {
	return &ChunkReader{
		r:         r,
		inChunkSz: 128, // default
		prev:      map[uint32]*ChunkHeader{},
		buf:       map[uint32][]byte{},
	}
}

func (rd *ChunkReader) SetInChunkSize(sz uint32) {
	if sz > 0 {
		rd.inChunkSz = sz
	}
}

func (rd *ChunkReader) ReadMessage() (*Chunk, error) {
	fmt, csid, err := rd.readBasicHeader()
	if err != nil {
		return nil, err
	}

	h, err := rd.readMessageHeader(fmt, csid)
	if err != nil {
		return nil, err
	}

	// read payload chunk piece
	need := h.Length - uint32(len(rd.buf[csid]))
	take := rd.inChunkSz
	if need < take {
		take = need
	}

	part := make([]byte, take)
	if _, err := io.ReadFull(rd.r, part); err != nil {
		return nil, err
	}
	rd.buf[csid] = append(rd.buf[csid], part)

	if uint32(len(rd.buf[csid])) < h.Length {
		// not complete yet; continue reading next chunk(s)
		return rd.ReadMessage()
	}

	payload := rd.buf[csid]
	delete(rd.buf, csid)

	return &Chunk{
		Fmt:       fmt,
		CSID:      csid,
		Timestamp: h.Timestamp,
		TypeID:    h.TypeID,
		StreamID:  h.StreamID,
		Payload:   payload,
	}, nil
}

func (rd *ChunkReader) readBasicHeader() (uint8, uint32, error) {
	var b0 [1]byte
	if _, err := io.ReadFull(rd.r, b0[:]); err != nil {
		return 0, 0, err
	}

	fmt := b0[0] >> 6
	csid := uint32(b0[0] & 0x3f)

	if csid == 0 {
		var b [1]byte
		if _, err := io.ReadFull(rd.r, b[:]); err != nil {
			return 0, 0, err
		}
		csid = 64 + uint32(b[0])
	} else if csid == 1 {
		var b [2]byte
		if _, err := io.ReadFull(rd.r, b[:]); err != nil {
			return 0, 0, err
		}
		csid = 64 + uint32(b[0]) + uint32(b[1])*256
	}

	return fmt, csid, nil
}

func (rd *ChunkReader) readMessageHeader(fmt uint8, csid uint32) (*ChunkHeader, error) {
	prev := rd.prev[csid]
	if prev == nil {
		prev = &ChunkHeader{}
	}

	switch fmt {
	case 0:
		mh := make([]byte, 11)
		if _, err := io.ReadFull(rd.r, mh); err != nil {
			return nil, err
		}
		ts := uint32(mh[0])<<16 | uint32(mh[1])<<8 | uint32(mh[2])
		l := uint32(mh[3])<<16 | uint32(mh[4])<<8 | uint32(mh[5])
		tid := mh[6]
		sid := binary.LittleEndian.Uint32(mh[7:11])

		prev.Timestamp = ts
		prev.Length = l
		prev.TypeID = tid
		prev.StreamID = sid

	case 1:
		mh := make([]byte, 7)
		if _, err := io.ReadFull(rd.r, mh); err != nil {
			return nil, err
		}
		dts := uint32(mh[0])<<16 | uint32(mh[1])<<8 | uint32(mh[2])
		l := uint32(mh[3])<<16 | uint32(mh[4])<<8 | uint32(mh[5])
		tid := mh[6]

		prev.Timestamp += dts
		prev.Length = l
		prev.TypeID = tid
		// StreamID same as prev

	case 2:
		mh := make([]byte, 3)
		if _, err := io.ReadFull(rd.r, mh); err != nil {
			return nil, err
		}
		dts := uint32(mh[0])<<16 | uint32(mh[1])<<8 | uint32(mh[2])
		prev.Timestamp += dts
		// Length/TypeID/StreamID same

	case 3:
		// all same as prev
	}

	rd.prev[csid] = prev
	return prev, nil
}
