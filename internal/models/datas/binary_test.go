package datas

import (
	"testing"
	"time"
)

func TestData_SetBinary(t *testing.T) {
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
			"Set invalid Hex data",
			fields{},
			args{
				"invalidhex",
			},
			true,
		},
		{
			"Set Hex data",
			fields{
				Value: "123",
			},
			args{
				"313233",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Data{
				Type:        tt.fields.Type,
				Description: tt.fields.Description,
				Value:       tt.fields.Value,
			}
			err := d.SetBinary(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Data.SetBinary() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}

		})
	}
}
