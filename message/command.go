package message

import (
	"bytes"

	"github.com/dongmedia/go-rtmp/amf"
)

type Command struct {
	Name          string
	TransactionID uint64
}

func DecodeCommand(payload []byte) (*Command, error) {
	r := bytes.NewReader(payload)

	name, _ := amf.Decode(r)
	tx, _ := amf.Decode(r)

	return &Command{
		Name:          name.(string),
		TransactionID: tx.(uint64),
	}, nil
}
