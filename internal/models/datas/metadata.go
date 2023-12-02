package datas

import (
	"errors"
	"time"
)

var ErrInvalidType = errors.New("invalid type")

type IData interface {
	ID() int
	SetID(id int)
	UserID() int
	SetUserID(id int)
	Description() string
	SetDescription(desc string)
	CreatedAt() time.Time
	EditedAt() time.Time
	SetValue(value string) error
	Value() string
}

type DataType int

const (
	BankCardType DataType = iota
	BinaryType
	CredentialsType
	TextType
)

func ParseDataType(d string) (DataType, error) {
	switch d {
	case BankCardType.String():
		return BankCardType, nil
	case BinaryType.String():
		return BinaryType, nil
	case CredentialsType.String():
		return CredentialsType, nil
	case TextType.String():
		return TextType, nil
	default:
		return 0, ErrInvalidType
	}
}

var Types []string = []string{
	BankCardType.String(),
	BinaryType.String(),
	CredentialsType.String(),
	TextType.String(),
}

// Data - контейнер для всех типов данных
type Data struct {
	ID          int       `json:"id,omitempty"`
	UserID      int       `json:"uid,omitempty"`
	Type        DataType  `json:"type"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	EditedAt    time.Time `json:"edited_at,omitempty"`
	Value       string    `json:"value"`
}

func NewData(t DataType, value string) (*Data, error) {
	now := time.Now()
	d := Data{
		Type:      t,
		CreatedAt: now,
		EditedAt:  now,
	}
	err := d.SetValue(value)
	return &d, err
}

func (d *Data) SetValue(value string) error {
	switch d.Type {
	case BinaryType:
		return d.SetBinary(value)
	case TextType:
		return d.SetText(value)
	case BankCardType:
		return d.SetBankCard(value)
	case CredentialsType:
		return d.SetCredentials(value)
	default:
		return ErrInvalidType
	}
}

func (d *Data) editNow() {
	d.EditedAt = time.Now()
}

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
