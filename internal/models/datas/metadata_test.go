package datas

import (
	"reflect"
	"testing"
	"time"
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDataType(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDataType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewData(tt.args.t, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestData_SetValue(t *testing.T) {
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
			if err := tt.d.SetValue(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Data.SetValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestData_editNow(t *testing.T) {
	tests := []struct {
		name string
		d    *Data
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.d.editNow()
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

func Test_metaData_ID(t *testing.T) {
	tests := []struct {
		name string
		md   metaData
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.md.ID(); got != tt.want {
				t.Errorf("metaData.ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_metaData_SetID(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name string
		md   *metaData
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.md.SetID(tt.args.id)
		})
	}
}

func Test_metaData_UserID(t *testing.T) {
	tests := []struct {
		name string
		md   metaData
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.md.UserID(); got != tt.want {
				t.Errorf("metaData.UserID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_metaData_SetUserID(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name string
		md   *metaData
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.md.SetUserID(tt.args.id)
		})
	}
}

func Test_metaData_Description(t *testing.T) {
	tests := []struct {
		name string
		md   metaData
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.md.Description(); got != tt.want {
				t.Errorf("metaData.Description() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_metaData_SetDescription(t *testing.T) {
	type args struct {
		desc string
	}
	tests := []struct {
		name string
		md   *metaData
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.md.SetDescription(tt.args.desc)
		})
	}
}

func Test_metaData_CreatedAt(t *testing.T) {
	tests := []struct {
		name string
		md   metaData
		want time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.md.CreatedAt(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("metaData.CreatedAt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_metaData_EditedAt(t *testing.T) {
	tests := []struct {
		name string
		md   metaData
		want time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.md.EditedAt(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("metaData.EditedAt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_metaData_editNow(t *testing.T) {
	tests := []struct {
		name string
		md   *metaData
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.md.editNow()
		})
	}
}
