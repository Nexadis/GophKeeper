package datas

import (
	"encoding/hex"
	"fmt"
)

var ErrBinaryInvalidValue = "invalid value"

type binary struct {
	metaData
	data []byte
}

func (b binary) Type() DataType {
	return BinaryType
}

func NewBinary(data []byte) *binary {
	b := binary{}
	copy(b.data, data)
	b.metaData = newMetaData()
	return &b
}

func (b binary) Value() string {
	return hex.EncodeToString(b.data)
}

func (b *binary) SetValue(value string) error {
	b.editNow()
	data, err := hex.DecodeString(value)
	if err != nil {
		return fmt.Errorf(ErrBinaryInvalidValue, err)
	}

	b.data = data
	return nil
}

func (d *Data) SetBinary(value string) error {
	d.editNow()
	_, err := hex.DecodeString(value)
	return err
}
