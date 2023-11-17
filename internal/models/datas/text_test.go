package datas

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextType(t *testing.T) {
	b := NewText("my text")
	assert.NotEmpty(t, b)
	assert.Equal(t, TextType, b.Type())
}

func TestTextSetValue(t *testing.T) {
	b := NewText("my text")
	assert.NotEmpty(t, b)
	val := b.Value()
	err := b.SetValue(val)
	assert.NoError(t, err)
	assert.Equal(t, val, b.Value())
	err = b.SetValue("new text")
	assert.NoError(t, err)
}
