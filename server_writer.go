package gortmp

import (
	"bytes"
	"encoding/binary"
	"net"
)

func writeConnectSuccess(conn net.Conn, tx uint64) {
	var amf bytes.Buffer

	// _result
	amf.Write([]byte{0x02, 0x00, 0x07})
	amf.WriteString("_result")

	// transaction id
	amf.WriteByte(0x00)
	binary.Write(&amf, binary.BigEndian, float64(tx))

	// null
	amf.Write([]byte{0x05})

	payload := amf.Bytes()

	// chunk (fmt0, csid=3)
	conn.Write([]byte{0x03})
	conn.Write([]byte{
		0x00, 0x00, 0x00,
		byte(len(payload) >> 16), byte(len(payload) >> 8), byte(len(payload)),
		0x14,
		0x00, 0x00, 0x00, 0x00,
	})
	conn.Write(payload)
}
