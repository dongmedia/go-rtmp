package gortmp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

func writeWindowAckSize(conn net.Conn, size uint32) error {
	// Type 5: Window Acknowledgement Size
	chunk := []byte{
		0x02,             // fmt=0, csid=2 (protocol control)
		0x00, 0x00, 0x00, // timestamp
		0x00, 0x00, 0x04, // message length = 4
		0x05,                   // type id = 5
		0x00, 0x00, 0x00, 0x00, // stream id = 0
		byte(size >> 24), byte(size >> 16), byte(size >> 8), byte(size),
	}
	_, err := conn.Write(chunk)
	return err
}

func writeSetPeerBandwidth(conn net.Conn, size uint32, limitType byte) error {
	// Type 6: Set Peer Bandwidth
	chunk := []byte{
		0x02,             // fmt=0, csid=2 (protocol control)
		0x00, 0x00, 0x00, // timestamp
		0x00, 0x00, 0x05, // message length = 5
		0x06,                   // type id = 6
		0x00, 0x00, 0x00, 0x00, // stream id = 0
		byte(size >> 24), byte(size >> 16), byte(size >> 8), byte(size),
		limitType, // 0=hard, 1=soft, 2=dynamic
	}
	_, err := conn.Write(chunk)
	return err
}

func writeStreamBegin(conn net.Conn, streamID uint32) error {
	// Type 4: User Control Message (Stream Begin = 0)
	chunk := []byte{
		0x02,             // fmt=0, csid=2
		0x00, 0x00, 0x00, // timestamp
		0x00, 0x00, 0x06, // message length = 6
		0x04,                   // type id = 4 (user control)
		0x00, 0x00, 0x00, 0x00, // stream id = 0
		0x00, 0x00, // event type = 0 (stream begin)
		byte(streamID >> 24), byte(streamID >> 16), byte(streamID >> 8), byte(streamID),
	}
	_, err := conn.Write(chunk)
	return err
}

// writeCommandResult sends a simple _result response for commands like releaseStream, FCPublish
func writeCommandResult(conn net.Conn, tx uint64) error {
	var amf bytes.Buffer

	// _result
	if _, err := amf.Write([]byte{0x02, 0x00, 0x07}); err != nil {
		return fmt.Errorf("write amf result data err: %v", err)
	}

	if _, err := amf.WriteString("_result"); err != nil {
		return fmt.Errorf("write amf result string err: %v", err)
	}

	// transaction id
	if err := amf.WriteByte(0x00); err != nil {
		return fmt.Errorf("write amf transaction id byte err: %v", err)
	}
	if err := binary.Write(&amf, binary.BigEndian, float64(tx)); err != nil {
		return fmt.Errorf("write binary transaction err: %v", err)
	}

	// null (command object)
	if _, err := amf.Write([]byte{0x05}); err != nil {
		return fmt.Errorf("write amf command object null err: %v", err)
	}

	// undefined (response)
	if _, err := amf.Write([]byte{0x06}); err != nil {
		return fmt.Errorf("write amf response err: %v", err)
	}

	payload := amf.Bytes()

	if _, err := conn.Write([]byte{0x03}); err != nil {
		return fmt.Errorf("write connection payload err: %v", err)
	} // fmt=0, csid=3
	if _, err := conn.Write([]byte{
		0x00, 0x00, 0x00,
		byte(len(payload) >> 16), byte(len(payload) >> 8), byte(len(payload)),
		0x14,
		0x00, 0x00, 0x00, 0x00,
	}); err != nil {
		return fmt.Errorf("write connection payload data err: %v", err)
	}

	_, err := conn.Write(payload)
	return err
}

func writeConnectSuccess(conn net.Conn, tx uint64) error {
	var amf bytes.Buffer

	// _result (string marker + length + string)
	if _, err := amf.Write([]byte{0x02, 0x00, 0x07}); err != nil {
		return fmt.Errorf("write amf result marker err: %v", err)
	}
	if _, err := amf.WriteString("_result"); err != nil {
		return fmt.Errorf("write amf result string err: %v", err)
	}

	// transaction id (number marker + float64)
	if err := amf.WriteByte(0x00); err != nil {
		return fmt.Errorf("write amf tx marker err: %v", err)
	}
	if err := binary.Write(&amf, binary.BigEndian, float64(tx)); err != nil {
		return fmt.Errorf("write amf tx err: %v", err)
	}

	// Properties object
	if err := amf.WriteByte(0x03); err != nil {
		return fmt.Errorf("write props object marker err: %v", err)
	}
	// fmsVer
	if _, err := amf.Write([]byte{0x00, 0x06}); err != nil {
		return fmt.Errorf("write fmsVer key len err: %v", err)
	}
	if _, err := amf.WriteString("fmsVer"); err != nil {
		return fmt.Errorf("write fmsVer key err: %v", err)
	}
	if _, err := amf.Write([]byte{0x02, 0x00, 0x0d}); err != nil {
		return fmt.Errorf("write fmsVer val marker err: %v", err)
	}
	if _, err := amf.WriteString("FMS/3,0,1,123"); err != nil {
		return fmt.Errorf("write fmsVer val err: %v", err)
	}
	// capabilities
	if _, err := amf.Write([]byte{0x00, 0x0c}); err != nil {
		return fmt.Errorf("write capabilities key len err: %v", err)
	}
	if _, err := amf.WriteString("capabilities"); err != nil {
		return fmt.Errorf("write capabilities key err: %v", err)
	}
	if err := amf.WriteByte(0x00); err != nil {
		return fmt.Errorf("write capabilities val marker err: %v", err)
	}
	if err := binary.Write(&amf, binary.BigEndian, float64(31)); err != nil {
		return fmt.Errorf("write capabilities val err: %v", err)
	}
	// object end
	if _, err := amf.Write([]byte{0x00, 0x00, 0x09}); err != nil {
		return fmt.Errorf("write props object end err: %v", err)
	}

	// Information object
	if err := amf.WriteByte(0x03); err != nil {
		return fmt.Errorf("write info object marker err: %v", err)
	}
	// level
	if _, err := amf.Write([]byte{0x00, 0x05}); err != nil {
		return fmt.Errorf("write level key len err: %v", err)
	}
	if _, err := amf.WriteString("level"); err != nil {
		return fmt.Errorf("write level key err: %v", err)
	}
	if _, err := amf.Write([]byte{0x02, 0x00, 0x06}); err != nil {
		return fmt.Errorf("write level val marker err: %v", err)
	}
	if _, err := amf.WriteString("status"); err != nil {
		return fmt.Errorf("write level val err: %v", err)
	}
	// code
	if _, err := amf.Write([]byte{0x00, 0x04}); err != nil {
		return fmt.Errorf("write code key len err: %v", err)
	}
	if _, err := amf.WriteString("code"); err != nil {
		return fmt.Errorf("write code key err: %v", err)
	}
	if _, err := amf.Write([]byte{0x02, 0x00, 0x1d}); err != nil {
		return fmt.Errorf("write code val marker err: %v", err)
	}
	if _, err := amf.WriteString("NetConnection.Connect.Success"); err != nil {
		return fmt.Errorf("write code val err: %v", err)
	}
	// description
	if _, err := amf.Write([]byte{0x00, 0x0b}); err != nil {
		return fmt.Errorf("write description key len err: %v", err)
	}
	if _, err := amf.WriteString("description"); err != nil {
		return fmt.Errorf("write description key err: %v", err)
	}
	if _, err := amf.Write([]byte{0x02, 0x00, 0x15}); err != nil {
		return fmt.Errorf("write description val marker err: %v", err)
	}
	if _, err := amf.WriteString("Connection succeeded."); err != nil {
		return fmt.Errorf("write description val err: %v", err)
	}
	// objectEncoding
	if _, err := amf.Write([]byte{0x00, 0x0e}); err != nil {
		return fmt.Errorf("write objectEncoding key len err: %v", err)
	}
	if _, err := amf.WriteString("objectEncoding"); err != nil {
		return fmt.Errorf("write objectEncoding key err: %v", err)
	}
	if err := amf.WriteByte(0x00); err != nil {
		return fmt.Errorf("write objectEncoding val marker err: %v", err)
	}
	if err := binary.Write(&amf, binary.BigEndian, float64(0)); err != nil {
		return fmt.Errorf("write objectEncoding val err: %v", err)
	}
	// object end
	if _, err := amf.Write([]byte{0x00, 0x00, 0x09}); err != nil {
		return fmt.Errorf("write info object end err: %v", err)
	}

	payload := amf.Bytes()

	// chunk header (fmt0, csid=3)
	if _, err := conn.Write([]byte{0x03}); err != nil {
		return fmt.Errorf("write chunk header err: %v", err)
	}
	if _, err := conn.Write([]byte{
		0x00, 0x00, 0x00, // timestamp
		byte(len(payload) >> 16), byte(len(payload) >> 8), byte(len(payload)),
		0x14,                   // type id = 20 (command)
		0x00, 0x00, 0x00, 0x00, // stream id = 0
	}); err != nil {
		return fmt.Errorf("write chunk msg header err: %v", err)
	}
	if _, err := conn.Write(payload); err != nil {
		return fmt.Errorf("write chunk payload err: %v", err)
	}

	return nil
}

func writePublishStart(conn net.Conn, streamID uint32) error {
	var amf bytes.Buffer

	// onStatus
	if _, err := amf.Write([]byte{0x02, 0x00, 0x08}); err != nil {
		return fmt.Errorf("write amf on status data err: %v", err)
	}
	if _, err := amf.WriteString("onStatus"); err != nil {
		return fmt.Errorf("write amf on status string err: %v", err)
	}

	// transaction id = 0
	if _, err := amf.Write([]byte{0x00}); err != nil {
		return fmt.Errorf("write amf transaction id err: %v", err)
	}

	if err := binary.Write(&amf, binary.BigEndian, float64(0)); err != nil {
		return fmt.Errorf("write binary data err: %v", err)
	}

	// null
	if _, err := amf.Write([]byte{0x05}); err != nil {
		return fmt.Errorf("write amf null data err: %v", err)
	}

	// object
	if _, err := amf.Write([]byte{
		0x03,
		0x00, 0x04, 'c', 'o', 'd', 'e',
		0x02, 0x00, 0x17,
	}); err != nil {
		return fmt.Errorf("write amf object data err: %v", err)
	}
	if _, err := amf.WriteString("NetStream.Publish.Start"); err != nil {
		return fmt.Errorf("write amf net stream publish startstring err: %v", err)
	}

	if _, err := amf.Write([]byte{
		0x00, 0x0b, 'd', 'e', 's', 'c', 'r', 'i', 'p', 't', 'i', 'o', 'n',
		0x02, 0x00, 0x0f,
	}); err != nil {
		return fmt.Errorf("write amf data err: %v", err)
	}
	if _, err := amf.WriteString("Publish succeeded"); err != nil {
		return fmt.Errorf("write amf publish success string err: %v", err)
	}

	// object end
	if _, err := amf.Write([]byte{0x00, 0x00, 0x09}); err != nil {
		return fmt.Errorf("write amf object end data err: %v", err)
	}

	payload := amf.Bytes()

	if _, err := conn.Write([]byte{0x03}); err != nil {
		return fmt.Errorf("write connection data err: %v", err)
	}

	if _, err := conn.Write([]byte{
		0x00, 0x00, 0x00,
		byte(len(payload) >> 16), byte(len(payload) >> 8), byte(len(payload)),
		0x14,
		byte(streamID), 0x00, 0x00, 0x00,
	}); err != nil {
		return fmt.Errorf("write connection payload data with stream id err: %v", err)
	}
	if _, err := conn.Write(payload); err != nil {
		return fmt.Errorf("write connection payload err: %v", err)
	}

	return nil
}

func writeCreateStreamResult(conn net.Conn, tx uint64, streamID uint32) error {
	var amf bytes.Buffer

	// _result
	if _, err := amf.Write([]byte{0x02, 0x00, 0x07}); err != nil {
		return fmt.Errorf("write amf data err: %v", err)
	}
	if _, err := amf.WriteString("_result"); err != nil {
		return fmt.Errorf("write amf result string err: %v", err)
	}

	// transaction id
	if err := amf.WriteByte(0x00); err != nil {
		return fmt.Errorf("write amf byte err: %v", err)
	}
	if err := binary.Write(&amf, binary.BigEndian, float64(tx)); err != nil {
		return fmt.Errorf("write binary data err: %v", err)
	}

	// null
	if _, err := amf.Write([]byte{0x05}); err != nil {
		return fmt.Errorf("write amf null data err: %v", err)
	}

	// stream id
	if err := amf.WriteByte(0x00); err != nil {
		return fmt.Errorf("write stream id byte err: %v", err)
	}

	if err := binary.Write(&amf, binary.BigEndian, float64(streamID)); err != nil {
		return fmt.Errorf("write binary err: %v", err)
	}

	payload := amf.Bytes()

	if _, err := conn.Write([]byte{0x03}); err != nil {
		return fmt.Errorf("write connection data err: %v", err)
	}

	if _, err := conn.Write([]byte{
		0x00, 0x00, 0x00,
		byte(len(payload) >> 16), byte(len(payload) >> 8), byte(len(payload)),
		0x14,
		0x00, 0x00, 0x00, 0x00,
	}); err != nil {
		return fmt.Errorf("write connection payload data err: %v", err)
	}

	if _, err := conn.Write(payload); err != nil {
		return fmt.Errorf("write connection payload err: %v", err)
	}

	return nil
}
