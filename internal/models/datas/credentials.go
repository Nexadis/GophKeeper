package datas

import (
	"errors"
	"fmt"
	"strings"
)

const (
	credentialsFormat = "%s;%s"
	credentialssep    = ";"
)

var ErrCredsInvalidFormat = errors.New("invalid creds format")

type credentials struct {
	metaData
	login    string
	password string
}

// NewCredentials - создаёт структуру для работы с логином и паролем
func NewCredentials(login string, password string) *credentials {
	c := credentials{
		login:    login,
		password: password,
	}
	c.metaData = newMetaData()
	return &c
}

// Value - возвращает логин и пароль в виде одной строки в заданном формате
func (c credentials) Value() string {
	return fmt.Sprintf(credentialsFormat, c.login, c.password)
}

// SetValue - валидирует значение из строки, если она соответствует заданному формату, то парсит её и заполняет структуру
func (c *credentials) SetValue(value string) error {
	values := strings.Split(value, credentialssep)
	if len(values) != 2 {
		return fmt.Errorf("%w: %q", ErrCredsInvalidFormat, value)
	}
	c.login = values[0]
	c.password = values[1]

	return nil
}

// SetCredentials - проверяет строку согласно формату и если всё правильно, то записывает её в поле Value
func (d *Data) SetCredentials(value string) error {
	d.editNow()
	c := credentials{}
	err := c.SetValue(value)
	if err != nil {
		return fmt.Errorf("%w: %q", ErrCredsInvalidFormat, value)
	}
	d.Value = value
	return nil
}

// CredentialsValues - позволяет получить отдельные значения логина и пароля из сохраненного Value
func (d *Data) CredentialsValues() (login string, password string) {
	c := credentials{}
	c.SetValue(d.Value)
	login = c.login
	password = c.password
	return
}
