package gortmp

import "encoding/binary"

func writeSetChunkSize(w *ChunkWriter, size uint32) error {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, size)
	// type 1, streamID 0
	return w.WriteMessage(2, 0, 1, 0, b)
}

func writeWindowAckSize(w *ChunkWriter, size uint32) error {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, size)
	// type 5, streamID 0
	return w.WriteMessage(2, 0, 5, 0, b)
}

func writeSetPeerBandwidth(w *ChunkWriter, size uint32, limitType uint8) error {
	b := make([]byte, 5)
	binary.BigEndian.PutUint32(b[:4], size)
	b[4] = limitType // 2=dynamic
	// type 6, streamID 0
	return w.WriteMessage(2, 0, 6, 0, b)
}
