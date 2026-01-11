package message

import (
	"bytes"

	"github.com/dongmedia/go-rtmp/amf"
)

type Command struct {
	Name          string
	TransactionID uint64
	Args          []any
}

// publish args
func DecodeCommand(payload []byte) (*Command, error) {
	r := bytes.NewReader(payload)

	name, _ := amf.Decode(r)
	tx, _ := amf.Decode(r)

	var args []any
	for r.Len() > 0 {
		v, err := amf.Decode(r)
		if err != nil {
			break
		}
		args = append(args, v)
	}

	return &Command{
		Name:          name.(string),
		TransactionID: tx.(uint64),
		Args:          args,
	}, nil
}
