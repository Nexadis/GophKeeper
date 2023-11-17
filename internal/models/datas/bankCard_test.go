package datas

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type card struct {
	number string
	holder string
	expire string
	cvv    int
}

var tests = []struct {
	name string
	card card
	err  error
}{
	{
		"Valid Card",
		card{
			"1111 2222 3333 4444",
			"My Name",
			"06/28",
			123,
		},
		nil,
	},
	{
		"Invalid number length",
		card{
			"1111 2222 3333 44441",
			"My Name",
			"06/28",
			123,
		},
		ErrCardInvalidNumber,
	},
	{
		"Invalid number",
		card{
			"LAKJ AKJL ALKJ KWER",
			"My Name",
			"06/28",
			123,
		},
		ErrCardInvalidNumber,
	},
	{
		"Invalid expire date",
		card{
			"1111 2222 3333 4444",
			"My Name",
			"98/43",
			123,
		},
		ErrCardInvalidExpire,
	},
	{
		"Invalid expire date",
		card{
			"1111 2222 3333 4444",
			"My Name",
			"09042",
			123,
		},
		ErrCardInvalidExpire,
	},
	{
		"Invalid cvv",
		card{
			"1111 2222 3333 4444",
			"My Name",
			"09/25",
			20,
		},
		ErrCardInvalidCVV,
	},
	{
		"Invalid cvv",
		card{
			"1111 2222 3333 4444",
			"My Name",
			"09/25",
			12020,
		},
		ErrCardInvalidCVV,
	},
}

func TestNewBankCard(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewBankCard(test.card.number, test.card.holder, test.card.expire, test.card.cvv)
			if !errors.Is(err, test.err) {
				assert.Equal(t, test.err, err)
			}
		},
		)
	}
}

func TestBankCardType(t *testing.T) {
	valid := card{
		"1111 2222 3333 4444",
		"My Name",
		"06/28",
		123,
	}
	c, err := NewBankCard(valid.number, valid.holder, valid.expire, valid.cvv)
	assert.NoError(t, err)
	assert.Equal(t, BankCardType, c.Type())
}

func TestBankCardSetValue(t *testing.T) {
	valid := card{
		"1111 2222 3333 4444",
		"My Name",
		"06/28",
		123,
	}
	c, err := NewBankCard(valid.number, valid.holder, valid.expire, valid.cvv)
	assert.NoError(t, err)
	val := c.Value()
	err = c.SetValue(val)
	assert.NoError(t, err)
	assert.Equal(t, val, c.Value())
}
