package datas

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMetadata(t *testing.T) {
	d := newMetaData()
	require.NotEmpty(t, d)
	require.NotEmpty(t, d.CreatedAt())
	require.NotEmpty(t, d.EditedAt())
}

func TestEditMetadata(t *testing.T) {
	d := newMetaData()
	assert.NotEmpty(t, d)
	before := d.editedAt
	d.editNow()
	require.NotEqual(t, before, d.editedAt)
}

func TestDataTypes(t *testing.T) {
	assert.NotEmpty(t, BankCardType.String())
	assert.NotEmpty(t, BinaryType.String())
	assert.NotEmpty(t, CredentialsType.String())
	assert.NotEmpty(t, TextType.String())
}

func TestID(t *testing.T) {
	id := 1
	d := newMetaData()
	assert.NotEmpty(t, d)
	d.SetID(id)
	assert.Equal(t, d.ID(), id)
}

func TestUserID(t *testing.T) {
	id := 2
	d := newMetaData()
	assert.NotEmpty(t, d)
	d.SetUserID(id)
	assert.Equal(t, d.UserID(), id)
}

func TestDescription(t *testing.T) {
	desc := "some desc"
	d := newMetaData()
	assert.NotEmpty(t, d)
	d.SetDescription(desc)
	assert.Equal(t, d.Description(), desc)
}
