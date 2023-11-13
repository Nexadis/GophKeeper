package datas

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCredentialsType(t *testing.T) {
	b := NewCredentials("login", "password")
	assert.NotEmpty(t, b)
	assert.Equal(t, CredentialsType, b.Type())
}

func TestCredentialsSetValue(t *testing.T) {
	b := NewCredentials("login", "password")
	assert.NotEmpty(t, b)
	val := b.Value()
	err := b.SetValue(val)
	assert.NoError(t, err)
	assert.Equal(t, val, b.Value())
	err = b.SetValue("invaliddata")
	assert.Error(t, err)
}
