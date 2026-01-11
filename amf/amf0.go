package amf

import (
	"encoding/binary"
	"errors"
	"io"
)

func Decode(r io.Reader) (any, error) {
	t := make([]byte, 1)
	if _, err := io.ReadFull(r, t); err != nil {
		return nil, err
	}

	switch t[0] {
	case 0x02: // string
		lenBuf := make([]byte, 2)
		io.ReadFull(r, lenBuf)
		l := binary.BigEndian.Uint16(lenBuf)
		buf := make([]byte, l)
		io.ReadFull(r, buf)
		return string(buf), nil

	case 0x00: // number
		buf := make([]byte, 8)
		io.ReadFull(r, buf)
		return binary.BigEndian.Uint64(buf), nil
	}

	return nil, errors.New("unsupported amf type")
}

func Encode(w io.Writer, v any) error {
	// TODO
	return nil
}
