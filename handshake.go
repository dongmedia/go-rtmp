package gortmp

import (
	"fmt"
	"io"
	"net"
)

type HandshakeService struct {
	RtmpVersion byte
	C0c1        []byte // C0+C1
	S0s1s2      []byte // S0+S1+S2
	C2          []byte
}

func NewHandshakeService() HandshakeService {
	return HandshakeService{
		RtmpVersion: 0x03,
	}
}

func (h *HandshakeService) Do(conn net.Conn) error {
	// C0+C1
	c0c1 := make([]byte, 1537)
	if _, err := io.ReadFull(conn, c0c1); err != nil {
		return err
	}

	if c0c1[0] != h.RtmpVersion {
		return fmt.Errorf("given rtmp version (%v) is unsupported rtmp version: %v", c0c1[0], h.RtmpVersion)
	}

	// S0+S1+S2
	s0s1s2 := make([]byte, 3073)
	s0s1s2[0] = h.RtmpVersion
	copy(s0s1s2[1:], c0c1[1:])

	if _, err := conn.Write(s0s1s2); err != nil {
		return fmt.Errorf("write handshake err: %v", err)
	}

	// C2
	c2 := make([]byte, 1536)
	_, err := io.ReadFull(conn, c2)
	return err
}
