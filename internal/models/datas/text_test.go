package datas

import (
	"reflect"
	"testing"
)

func TestData_SetText(t *testing.T) {
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
			if err := tt.d.SetText(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Data.SetText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_text_Type(t *testing.T) {
	tests := []struct {
		name string
		tr   text
		want DataType
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Type(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("text.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewText(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name string
		args args
		want *text
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewText(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_text_Value(t *testing.T) {
	tests := []struct {
		name string
		tr   text
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Value(); got != tt.want {
				t.Errorf("text.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_text_SetValue(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		tr      *text
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.tr.SetValue(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("text.SetValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
