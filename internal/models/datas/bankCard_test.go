package datas

import (
	"reflect"
	"testing"
	"time"
)

func TestData_SetBankCard(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		d       *Data
		args    args
		wantErr bool
	}{
		{
			"Empty value",
			&Data{},
			args{""},
			true,
		},
		{
			"Invalid CVV, not number",
			&Data{},
			args{"1234;name;time;CVV"},
			true,
		},
		{
			"Invalid Number",
			&Data{},
			args{"1;name;time;123"},
			true,
		},
		{
			"Invalid time",
			&Data{},
			args{"1234 1234 1234 1234;name;time;123"},
			true,
		},
		{
			"Invalid CVV, too long",
			&Data{},
			args{"1234 1234 1234 1234;name;10/14;1234"},
			true,
		},
		{
			"Valid Bank Card",
			&Data{},
			args{"1234 1234 1234 1234;name;10/24;123"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.d.SetBankCard(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Data.SetBankCard() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestData_BankCardValues(t *testing.T) {
	tests := []struct {
		name           string
		d              *Data
		wantNumber     string
		wantCardHolder string
		wantExpire     string
		wantCvv        int
	}{{
		"Valid Bank Card in data",
		&Data{
			Type:  BankCardType,
			Value: "1234 1234 1234 1234;Name;11/11;123",
		},
		"1234 1234 1234 1234",
		"Name",
		"11/11",
		123,
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNumber, gotCardHolder, gotExpire, gotCvv := tt.d.BankCardValues()
			if gotNumber != tt.wantNumber {
				t.Errorf("Data.BankCardValues() gotNumber = %v, want %v", gotNumber, tt.wantNumber)
			}
			if gotCardHolder != tt.wantCardHolder {
				t.Errorf(
					"Data.BankCardValues() gotCardHolder = %v, want %v",
					gotCardHolder,
					tt.wantCardHolder,
				)
			}
			if gotExpire != tt.wantExpire {
				t.Errorf("Data.BankCardValues() gotExpire = %v, want %v", gotExpire, tt.wantExpire)
			}
			if gotCvv != tt.wantCvv {
				t.Errorf("Data.BankCardValues() gotCvv = %v, want %v", gotCvv, tt.wantCvv)
			}
		})
	}
}

func Test_validateNumber(t *testing.T) {
	type args struct {
		number string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateNumber(tt.args.number)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("validateNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseExpire(t *testing.T) {
	type args struct {
		expire string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseExpire(tt.args.expire)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseExpire() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseExpire() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateCVV(t *testing.T) {
	type args struct {
		CVV int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateCVV(tt.args.CVV); (err != nil) != tt.wantErr {
				t.Errorf("validateCVV() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_bankCard_Type(t *testing.T) {
	tests := []struct {
		name string
		bk   bankCard
		want DataType
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bk.Type(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bankCard.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBankCard(t *testing.T) {
	type args struct {
		number     string
		cardHolder string
		expire     string
		cvv        int
	}
	tests := []struct {
		name    string
		args    args
		want    *bankCard
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBankCard(tt.args.number, tt.args.cardHolder, tt.args.expire, tt.args.cvv)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBankCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBankCard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bankCard_Value(t *testing.T) {
	tests := []struct {
		name string
		bc   bankCard
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bc.Value(); got != tt.want {
				t.Errorf("bankCard.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bankCard_SetValue(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		bc      *bankCard
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.bc.SetValue(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("bankCard.SetValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_bankCard_validateNumber(t *testing.T) {
	type args struct {
		number string
	}
	tests := []struct {
		name    string
		bc      bankCard
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.bc.validateNumber(tt.args.number)
			if (err != nil) != tt.wantErr {
				t.Errorf("bankCard.validateNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("bankCard.validateNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bankCard_parseExpire(t *testing.T) {
	type args struct {
		expire string
	}
	tests := []struct {
		name    string
		bc      bankCard
		args    args
		want    time.Time
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.bc.parseExpire(tt.args.expire)
			if (err != nil) != tt.wantErr {
				t.Errorf("bankCard.parseExpire() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bankCard.parseExpire() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bankCard_validateCVV(t *testing.T) {
	type args struct {
		CVV int
	}
	tests := []struct {
		name    string
		bc      bankCard
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.bc.validateCVV(tt.args.CVV); (err != nil) != tt.wantErr {
				t.Errorf("bankCard.validateCVV() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
