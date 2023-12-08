package datas

import (
	"reflect"
	"testing"
)

func Test_credentials_Type(t *testing.T) {
	tests := []struct {
		name string
		c    credentials
		want DataType
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Type(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("credentials.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCredentials(tt.args.login, tt.args.password); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCredentials() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_credentials_Value(t *testing.T) {
	tests := []struct {
		name string
		c    credentials
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Value(); got != tt.want {
				t.Errorf("credentials.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_credentials_SetValue(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		c       *credentials
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.SetValue(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("credentials.SetValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestData_SetCredentials(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		d       *Data
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.d.SetCredentials(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Data.SetCredentials() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestData_CredentialsValues(t *testing.T) {
	tests := []struct {
		name         string
		d            *Data
		wantLogin    string
		wantPassword string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLogin, gotPassword := tt.d.CredentialsValues()
			if gotLogin != tt.wantLogin {
				t.Errorf("Data.CredentialsValues() gotLogin = %v, want %v", gotLogin, tt.wantLogin)
			}
			if gotPassword != tt.wantPassword {
				t.Errorf("Data.CredentialsValues() gotPassword = %v, want %v", gotPassword, tt.wantPassword)
			}
		})
	}
}
