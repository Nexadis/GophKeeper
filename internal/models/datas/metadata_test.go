package datas

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseDataType(t *testing.T) {
	type args struct {
		d string
	}
	tests := []struct {
		name    string
		args    args
		want    DataType
		wantErr bool
	}{
		{
			"Parse BankCard",
			args{"Bank Card"},
			BankCardType,
			false,
		},
		{
			"Parse Binary",
			args{"Binary"},
			BinaryType,
			false,
		},
		{
			"Parse Credentials",
			args{"Credentials"},
			CredentialsType,
			false,
		},
		{
			"Parse Text",
			args{"Text"},
			TextType,
			false,
		},
		{
			"Invalid type",
			args{"Invalid"},
			-1,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDataType(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDataType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseDataType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewData(t *testing.T) {
	type args struct {
		t     DataType
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    *Data
		wantErr bool
	}{
		{
			"Valid Text",
			args{
				TextType,
				"text",
			},
			&Data{
				Type:  TextType,
				Value: "text",
			},
			false,
		},
		{
			"Valid Bank Card",
			args{
				BankCardType,
				"1234 1234 1234 1234;holder;11/11;123",
			},
			&Data{
				Type:  BankCardType,
				Value: "1234 1234 1234 1234;holder;11/11;123",
			},
			false,
		},
		{
			"Valid Binary",
			args{
				BinaryType,
				"3131",
			},
			&Data{
				Type:  BinaryType,
				Value: "3131",
			},
			false,
		},
		{
			"Valid Credentials",
			args{
				CredentialsType,
				"login;password",
			},
			&Data{
				Type:  CredentialsType,
				Value: "login;password",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewData(tt.args.t, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want.Type, got.Type)
			assert.Equal(t, tt.want.Value, got.Value)
		})
	}
}

func TestData_SetValue(t *testing.T) {
	type fields struct {
		ID          int
		UserID      int
		Type        DataType
		Description string
		CreatedAt   time.Time
		EditedAt    time.Time
		Value       string
	}
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Data{
				ID:          tt.fields.ID,
				UserID:      tt.fields.UserID,
				Type:        tt.fields.Type,
				Description: tt.fields.Description,
				CreatedAt:   tt.fields.CreatedAt,
				EditedAt:    tt.fields.EditedAt,
				Value:       tt.fields.Value,
			}
			if err := d.SetValue(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Data.SetValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestData_editNow(t *testing.T) {
	type fields struct {
		ID          int
		UserID      int
		Type        DataType
		Description string
		CreatedAt   time.Time
		EditedAt    time.Time
		Value       string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Data{
				ID:          tt.fields.ID,
				UserID:      tt.fields.UserID,
				Type:        tt.fields.Type,
				Description: tt.fields.Description,
				CreatedAt:   tt.fields.CreatedAt,
				EditedAt:    tt.fields.EditedAt,
				Value:       tt.fields.Value,
			}
			d.editNow()
		})
	}
}

func Test_newMetaData(t *testing.T) {
	tests := []struct {
		name string
		want metaData
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newMetaData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newMetaData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataType_String(t *testing.T) {
	tests := []struct {
		name string
		dt   DataType
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dt.String(); got != tt.want {
				t.Errorf("DataType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_metaData_editNow(t *testing.T) {
	type fields struct {
		id        int
		userID    int
		desc      string
		createdAt time.Time
		editedAt  time.Time
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := &metaData{
				id:        tt.fields.id,
				userID:    tt.fields.userID,
				desc:      tt.fields.desc,
				createdAt: tt.fields.createdAt,
				editedAt:  tt.fields.editedAt,
			}
			md.editNow()
		})
	}
}
