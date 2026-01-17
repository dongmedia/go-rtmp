package message

import (
	"bytes"
	"fmt"

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

	name, err := amf.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("amf decode name from request payload err: %v", err)
	}
	tx, err := amf.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("amf decode transaction from request payload err: %v", err)
	}

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
