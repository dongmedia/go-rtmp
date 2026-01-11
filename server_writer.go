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

func writePublishStart(conn net.Conn, streamID uint32) {
	var amf bytes.Buffer

	// onStatus
	amf.Write([]byte{0x02, 0x00, 0x08})
	amf.WriteString("onStatus")

	// transaction id = 0
	amf.Write([]byte{0x00})
	binary.Write(&amf, binary.BigEndian, float64(0))

	// null
	amf.Write([]byte{0x05})

	// object
	amf.Write([]byte{
		0x03,
		0x00, 0x04, 'c', 'o', 'd', 'e',
		0x02, 0x00, 0x17,
	})
	amf.WriteString("NetStream.Publish.Start")

	amf.Write([]byte{
		0x00, 0x0b, 'd', 'e', 's', 'c', 'r', 'i', 'p', 't', 'i', 'o', 'n',
		0x02, 0x00, 0x0f,
	})
	amf.WriteString("Publish succeeded")

	// object end
	amf.Write([]byte{0x00, 0x00, 0x09})

	payload := amf.Bytes()

	conn.Write([]byte{0x03})
	conn.Write([]byte{
		0x00, 0x00, 0x00,
		byte(len(payload) >> 16), byte(len(payload) >> 8), byte(len(payload)),
		0x14,
		byte(streamID), 0x00, 0x00, 0x00,
	})
	conn.Write(payload)
}

func writeCreateStreamResult(conn net.Conn, tx uint64, streamID uint32) {
	var amf bytes.Buffer

	// _result
	amf.Write([]byte{0x02, 0x00, 0x07})
	amf.WriteString("_result")

	// transaction id
	amf.WriteByte(0x00)
	binary.Write(&amf, binary.BigEndian, float64(tx))

	// null
	amf.Write([]byte{0x05})

	// stream id
	amf.WriteByte(0x00)
	binary.Write(&amf, binary.BigEndian, float64(streamID))

	payload := amf.Bytes()

	conn.Write([]byte{0x03})
	conn.Write([]byte{
		0x00, 0x00, 0x00,
		byte(len(payload) >> 16), byte(len(payload) >> 8), byte(len(payload)),
		0x14,
		0x00, 0x00, 0x00, 0x00,
	})
	conn.Write(payload)
}
