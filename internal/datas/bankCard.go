package datas

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	dateFormat     = "02/06"
	bankCardFormat = "%s %s %s %d"
)

var (
	ErrCardInvalidNumber = "bankCard.validateNumber %s: %w"
	ErrCardInvalidExpire = "bankCard.validateExpire %s: %w"
	ErrCardInvalidCVV    = "bankCard.validateCVV '%s' invalid number"
	ErrCardInvalidFormat = "bankCard.SetValue %w"
)

type bankCard struct {
	metaData
	number     string
	cardHolder string
	expire     time.Time
	cvv        int
}

func (bk bankCard) Type() DataType {
	return BinaryType
}

func NewBankCard(number, cardHolder, expire string) bankCard {
	bc := bankCard{}
	bc.metaData = newMetaData()
	return bc
}

func (bc bankCard) Value() string {
	return fmt.Sprintf(
		bankCardFormat,
		bc.number,
		bc.cardHolder,
		bc.expire,
		bc.cvv,
	)
}

func (bc *bankCard) SetValue(value string) error {
	bc.editNow()
	var number, cardHolder, expire string
	var cvv int
	_, err := fmt.Sscanf(value, bankCardFormat, number, cardHolder, expire, cvv)
	if err != nil {
		return fmt.Errorf(ErrCardInvalidFormat, err)
	}
	err = bc.validateNumber(number)
	if err != nil {
		return err
	}
	t, err := bc.parseExpire(expire)
	if err != nil {
		return err
	}
	err = bc.validateCVV(cvv)
	if err != nil {
		return err
	}
	bc.number = number
	bc.cardHolder = cardHolder
	bc.expire = t
	bc.cvv = cvv
	return nil
}

func (bc bankCard) validateNumber(number string) error {
	trimmedNum := strings.TrimSpace(number)
	cardnum := strings.Join(strings.Split(trimmedNum, " "), "")
	_, err := strconv.Atoi(cardnum)
	if err != nil {
		return fmt.Errorf(ErrCardInvalidNumber, number, err)
	}
	return nil
}

func (bc bankCard) parseExpire(expire string) (time.Time, error) {
	t, err := time.Parse(dateFormat, expire)
	if err != nil {
		return time.Time{}, fmt.Errorf(ErrCardInvalidExpire, expire, err)
	}
	return t, nil
}

func (bc bankCard) validateCVV(CVV int) error {
	if CVV <= 99 || CVV >= 1000 {
		return fmt.Errorf(ErrCardInvalidCVV, CVV)
	}
	return nil
}
