package datas

import (
	"encoding/hex"
)

var ErrBinaryInvalidValue = "invalid value"

// SetBinary - записывает в Value структуры Data бинарные данные в виде hex-строки
func (d *Data) SetBinary(value string) error {
	d.editNow()
	_, err := hex.DecodeString(value)
	if err != nil {
		return err
	}
	d.Value = value
	return nil
}
