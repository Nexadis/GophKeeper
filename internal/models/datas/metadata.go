package datas

import (
	"errors"
	"time"
)

var ErrInvalidType = errors.New("invalid type")

type DataType int

const (
	BankCardType DataType = iota
	BinaryType
	CredentialsType
	TextType
)

// ParseDataType - Возвращает тип данных по названию
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
		return -1, ErrInvalidType
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

// NewData - создаёт структуру Data для заданного типа данных, при этом валидируя value на соответсвие формату записи этих данных
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

// SetValue - валидирут полученные данные согласно формату типа, если всё впорядке, то записывает их
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

// String - возвращает название типа данных
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

func (md *metaData) editNow() {
	md.editedAt = time.Now()
}
