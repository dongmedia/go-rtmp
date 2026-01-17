package gortmp

import (
	"bytes"
	"encoding/binary"
	"net"

	"github.com/dongmedia/go-rtmp/amf"
)

func writeConnectSuccess(w *ChunkWriter, tx float64) error {
	var p bytes.Buffer
	_ = amf.EncodeString(&p, "_result")
	_ = amf.EncodeNumber(&p, tx)

	// props (간단히 null 처리해도 되지만, OBS 호환을 위해 object를 넣는 것이 안전함)
	// 여기서는 최소로 null + info object 형태로 단순화
	_ = amf.EncodeNull(&p)

	// info object (최소)
	// AMF0 object를 제대로 구현하면 좋지만, 여기서는 아주 최소로 “null”로도 통과하는 환경이 많습니다.
	// 다만 실패 시 object 구현을 추가해야 합니다.
	_ = amf.EncodeNull(&p)

	return w.WriteMessage(3, 0, 20, 0, p.Bytes())
}

func writeCreateStreamResult(w *ChunkWriter, tx float64, streamID uint32) error {
	var p bytes.Buffer
	_ = amf.EncodeString(&p, "_result")
	_ = amf.EncodeNumber(&p, tx)
	_ = amf.EncodeNull(&p)
	_ = amf.EncodeNumber(&p, float64(streamID))
	return w.WriteMessage(3, 0, 20, 0, p.Bytes())
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

// func writeConnectSuccess(conn net.Conn, tx uint64) {
// 	var amf bytes.Buffer

// 	// _result
// 	amf.Write([]byte{0x02, 0x00, 0x07})
// 	amf.WriteString("_result")

// 	// transaction id
// 	amf.WriteByte(0x00)
// 	binary.Write(&amf, binary.BigEndian, float64(tx))

// 	// null
// 	amf.Write([]byte{0x05})

// 	payload := amf.Bytes()

// 	// chunk (fmt0, csid=3)
// 	conn.Write([]byte{0x03})
// 	conn.Write([]byte{
// 		0x00, 0x00, 0x00,
// 		byte(len(payload) >> 16), byte(len(payload) >> 8), byte(len(payload)),
// 		0x14,
// 		0x00, 0x00, 0x00, 0x00,
// 	})
// 	conn.Write(payload)
// }

// func writeCreateStreamResult(conn net.Conn, tx uint64, streamID uint32) {
// 	var amf bytes.Buffer

// 	// _result
// 	amf.Write([]byte{0x02, 0x00, 0x07})
// 	amf.WriteString("_result")

// 	// transaction id
// 	amf.WriteByte(0x00)
// 	binary.Write(&amf, binary.BigEndian, float64(tx))

// 	// null
// 	amf.Write([]byte{0x05})

// 	// stream id
// 	amf.WriteByte(0x00)
// 	binary.Write(&amf, binary.BigEndian, float64(streamID))

// 	payload := amf.Bytes()

// 	conn.Write([]byte{0x03})
// 	conn.Write([]byte{
// 		0x00, 0x00, 0x00,
// 		byte(len(payload) >> 16), byte(len(payload) >> 8), byte(len(payload)),
// 		0x14,
// 		0x00, 0x00, 0x00, 0x00,
// 	})
// 	conn.Write(payload)
// }
