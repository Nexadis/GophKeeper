package datas

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCredentials(t *testing.T) {
	type args struct {
		login    string
		password string
	}
	tests := []struct {
		name string
		args args
		want *credentials
	}{
		{
			"Create credentials",
			args{
				"login",
				"password",
			},
			&credentials{
				login:    "login",
				password: "password",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCredentials(tt.args.login, tt.args.password)
			assert.Equal(
				t,
				tt.want.login,
				got.login,
			)
			assert.Equal(
				t,
				tt.want.password,
				got.password,
			)
		})
	}
}

func Test_credentials_Value(t *testing.T) {
	type fields struct {
		metaData metaData
		login    string
		password string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"Create credentials",
			fields{
				metaData{},
				"login",
				"password",
			},
			"login;password",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := credentials{
				metaData: tt.fields.metaData,
				login:    tt.fields.login,
				password: tt.fields.password,
			}
			if got := c.Value(); got != tt.want {
				t.Errorf("credentials.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestData_SetCredentials(t *testing.T) {
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
		{
			"Not splitable",
			fields{},
			args{"invalid credentials"},
			true,
		},
		{
			"Too much splits",
			fields{},
			args{"login;password;something"},
			true,
		},
		{
			"Valid value",
			fields{
				Value: "login;password",
			},
			args{"login;password"},
			false,
		},
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
			if err := d.SetCredentials(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Data.SetCredentials() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestData_CredentialsValues(t *testing.T) {
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
		name         string
		fields       fields
		wantLogin    string
		wantPassword string
	}{
		{
			"Parse Value",
			fields{
				Value: "login;password",
			},
			"login",
			"password",
		},
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
			gotLogin, gotPassword := d.CredentialsValues()
			if gotLogin != tt.wantLogin {
				t.Errorf("Data.CredentialsValues() gotLogin = %v, want %v", gotLogin, tt.wantLogin)
			}
			if gotPassword != tt.wantPassword {
				t.Errorf(
					"Data.CredentialsValues() gotPassword = %v, want %v",
					gotPassword,
					tt.wantPassword,
				)
			}
		})
	}
}
