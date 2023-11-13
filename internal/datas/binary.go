package datas

import (
	"encoding/base64"
	"fmt"
)

var ErrBinaryInvalidValue = "binary.SetValue: %w"

type binary struct {
	metaData
	data []byte
}

func (b binary) Type() DataType {
	return BinaryType
}

func NewBinary(data []byte) binary {
	b := binary{}
	copy(b.data, data)
	b.metaData = newMetaData()
	return b
}

func (b binary) Value() string {
	return base64.StdEncoding.EncodeToString(b.data)
}

func (b binary) SetValue(value string) error {
	b.editNow()
	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return fmt.Errorf(ErrBinaryInvalidValue, err)
	}

	b.data = data
	return nil
}
