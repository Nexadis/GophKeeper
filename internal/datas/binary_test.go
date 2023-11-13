package datas

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBinaryType(t *testing.T) {
	b := NewBinary([]byte("123"))
	assert.NotEmpty(t, b)
	assert.Equal(t, BinaryType, b.Type())
}

func TestBinarySetValue(t *testing.T) {
	b := NewBinary([]byte("123"))
	val := b.Value()
	err := b.SetValue(val)
	assert.NoError(t, err)
	assert.Equal(t, val, b.Value())
	err = b.SetValue("invaliddata")
	assert.Error(t, err)
}
