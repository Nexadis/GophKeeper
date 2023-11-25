package datas

import (
	"time"
)

type DataType int

const (
	BankCardType DataType = iota
	BinaryType
	CredentialsType
	TextType
)

type metaData struct {
	id        int
	userID    int
	desc      string
	createdAt time.Time
	editedAt  time.Time
}

func newMetaData() metaData {
	now := time.Now()
	return metaData{
		createdAt: now,
		editedAt:  now,
	}
}

func (dt DataType) String() string {
	var t string
	switch dt {
	case BankCardType:
		t = "Bank Card"
	case BinaryType:
		t = "Binary"
	case CredentialsType:
		t = "Credentials"
	case TextType:
		t = "Text"
	default:
		t = "Unknown type"
	}
	return t
}

func (md metaData) ID() int {
	return md.id
}

func (md *metaData) SetID(id int) {
	md.editNow()
	md.id = id
}

func (md metaData) UserID() int {
	return md.userID
}

func (md *metaData) SetUserID(id int) {
	md.editNow()
	md.userID = id
}

func (md metaData) Description() string {
	return md.desc
}

func (md *metaData) SetDescription(desc string) {
	md.editNow()
	md.desc = desc
}

func (md metaData) CreatedAt() time.Time {
	return md.createdAt
}

func (md metaData) EditedAt() time.Time {
	return md.editedAt
}

func (md *metaData) editNow() {
	md.editedAt = time.Now()
}
