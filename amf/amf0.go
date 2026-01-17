package amf

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
)

func Decode(r io.Reader) (any, error) {
	var t [1]byte
	if _, err := io.ReadFull(r, t[:]); err != nil {
		return nil, err
	}

	switch t[0] {
	case 0x02: // string
		var lb [2]byte
		if _, err := io.ReadFull(r, lb[:]); err != nil {
			return nil, err
		}
		l := binary.BigEndian.Uint16(lb[:])
		buf := make([]byte, l)
		if _, err := io.ReadFull(r, buf); err != nil {
			return nil, err
		}
		return string(buf), nil

	case 0x00: // number (IEEE754 float64, big-endian)
		var b [8]byte
		if _, err := io.ReadFull(r, b[:]); err != nil {
			return nil, err
		}
		u := binary.BigEndian.Uint64(b[:])
		return math.Float64frombits(u), nil

	case 0x05: // null
		return nil, nil
	}

	return nil, errors.New("unsupported amf0 type")
}

func EncodeString(w io.Writer, s string) error {
	_, err := w.Write([]byte{0x02, byte(len(s) >> 8), byte(len(s))})
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(s))
	return err
}

func EncodeNumber(w io.Writer, f float64) error {
	_, err := w.Write([]byte{0x00})
	if err != nil {
		return err
	}
	u := math.Float64bits(f)
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], u)
	_, err = w.Write(b[:])
	return err
}

func EncodeNull(w io.Writer) error {
	_, err := w.Write([]byte{0x05})
	return err
}
