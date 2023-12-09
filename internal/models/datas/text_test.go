package datas

import (
	"testing"
)

func TestNewText(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name string
		args args
		want *text
	}{
		{
			"New valid text",
			args{"text"},
			&text{data: "text"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewText(tt.args.data); got.data != tt.want.data {
				t.Errorf("NewText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_text_Value(t *testing.T) {
	type fields struct {
		metaData metaData
		data     string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"Show text",
			fields{
				metaData{},
				"text",
			},
			"text",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := text{
				metaData: tt.fields.metaData,
				data:     tt.fields.data,
			}
			if got := tr.Value(); got != tt.want {
				t.Errorf("text.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}
