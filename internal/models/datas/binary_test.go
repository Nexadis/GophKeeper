package datas

import (
	"reflect"
	"testing"
)

func Test_binary_Type(t *testing.T) {
	tests := []struct {
		name string
		b    binary
		want DataType
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Type(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("binary.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBinary(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want *binary
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBinary(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBinary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_binary_Value(t *testing.T) {
	tests := []struct {
		name string
		b    binary
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Value(); got != tt.want {
				t.Errorf("binary.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_binary_SetValue(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		b       *binary
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.SetValue(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("binary.SetValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestData_SetBinary(t *testing.T) {
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
			if err := tt.d.SetBinary(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Data.SetBinary() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
