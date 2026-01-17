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
		0x02,                   // fmt=0, csid=2 (protocol control)
		0x00, 0x00, 0x00,       // timestamp
		0x00, 0x00, 0x04,       // message length = 4
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
		0x02,                   // fmt=0, csid=2 (protocol control)
		0x00, 0x00, 0x00,       // timestamp
		0x00, 0x00, 0x05,       // message length = 5
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
		0x02,                   // fmt=0, csid=2
		0x00, 0x00, 0x00,       // timestamp
		0x00, 0x00, 0x06,       // message length = 6
		0x04,                   // type id = 4 (user control)
		0x00, 0x00, 0x00, 0x00, // stream id = 0
		0x00, 0x00, // event type = 0 (stream begin)
		byte(streamID >> 24), byte(streamID >> 16), byte(streamID >> 8), byte(streamID),
	}
	_, err := conn.Write(chunk)
	return err
}

func writeConnectSuccess(conn net.Conn, tx uint64) error {
	var amf bytes.Buffer

	// _result
	if _, err := amf.Write([]byte{0x02, 0x00, 0x07}); err != nil {
		return fmt.Errorf("write amf result err: %v", err)
	}

	if _, err := amf.WriteString("_result"); err != nil {
		return fmt.Errorf("write amf result string err: %v", err)
	}

	// transaction id
	if err := amf.WriteByte(0x00); err != nil {
		return fmt.Errorf("write amf transaction id err: %v", err)
	}

	if err := binary.Write(&amf, binary.BigEndian, float64(tx)); err != nil {
		return fmt.Errorf("write binary data err: %v", err)
	}

	// null
	if _, err := amf.Write([]byte{0x05}); err != nil {
		return fmt.Errorf("write amf null bytes err: %v", err)
	}

	payload := amf.Bytes()

	// chunk (fmt0, csid=3)
	if _, err := conn.Write([]byte{0x03}); err != nil {
		return fmt.Errorf("write chunk err: %v", err)
	}

	if _, err := conn.Write([]byte{
		0x00, 0x00, 0x00,
		byte(len(payload) >> 16), byte(len(payload) >> 8), byte(len(payload)),
		0x14,
		0x00, 0x00, 0x00, 0x00,
	}); err != nil {
		return fmt.Errorf("write chunk payload data err: %v", err)
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
