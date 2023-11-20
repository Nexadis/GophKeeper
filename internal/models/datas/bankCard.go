package datas

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	dateLayout     = "02/06"
	bankCardsep    = ";"
	bankCardFormat = "%s;%s;%s;%d"
)

var (
	ErrCardInvalidNumber = errors.New("invalid number")
	ErrCardInvalidExpire = errors.New("invalid expire date")
	ErrCardInvalidCVV    = errors.New("invalid cvv")
	ErrCardInvalidFormat = errors.New("can't parse card info")
)

type bankCard struct {
	metaData
	number     string
	cardHolder string
	expire     time.Time
	cvv        int
}

func (bk bankCard) Type() DataType {
	return BankCardType
}

func NewBankCard(number, cardHolder, expire string, cvv int) (*bankCard, error) {
	bc := bankCard{}
	number, err := bc.validateNumber(number)
	if err != nil {
		return nil, err
	}
	err = bc.validateCVV(cvv)
	if err != nil {
		return nil, err
	}
	exp, err := bc.parseExpire(expire)
	if err != nil {
		return nil, err
	}
	bc.metaData = newMetaData()
	bc.number = number
	bc.cardHolder = cardHolder
	bc.expire = exp
	bc.cvv = cvv
	return &bc, nil
}

func (bc bankCard) Value() string {
	return fmt.Sprintf(
		bankCardFormat,
		bc.number,
		bc.cardHolder,
		bc.expire.Format(dateLayout),
		bc.cvv,
	)
}

func (bc *bankCard) SetValue(value string) error {
	bc.editNow()
	values := strings.Split(value, bankCardsep)
	if len(values) != 4 {
		return fmt.Errorf("%w: %q", ErrCardInvalidFormat, value)
	}
	number := values[0]
	cardHolder := values[1]
	expire := values[2]
	cvv, err := strconv.Atoi(values[3])
	if err != nil {
		return fmt.Errorf("%w: %q: %q", ErrCardInvalidCVV, err, value)
	}

	number, err = bc.validateNumber(number)
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

func (bc bankCard) validateNumber(number string) (string, error) {
	trimmedNum := strings.TrimSpace(number)
	cardnum := strings.Join(strings.Split(trimmedNum, " "), "")
	_, err := strconv.Atoi(cardnum)
	if err != nil {
		return "", fmt.Errorf("%w %q", ErrCardInvalidNumber, err)
	}
	if len(cardnum) != 16 {
		return "", fmt.Errorf("%w: %q", ErrCardInvalidNumber, number)
	}
	return cardnum, nil
}

func (bc bankCard) parseExpire(expire string) (time.Time, error) {
	t, err := time.Parse(dateLayout, expire)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %q", ErrCardInvalidExpire, expire)
	}
	return t, nil
}

func (bc bankCard) validateCVV(CVV int) error {
	if CVV <= 99 || CVV >= 1000 {
		return fmt.Errorf("%w: %q", ErrCardInvalidCVV, CVV)
	}
	return nil
}
