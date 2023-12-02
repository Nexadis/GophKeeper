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

func (c credentials) Type() DataType {
	return CredentialsType
}

func NewCredentials(login string, password string) *credentials {
	c := credentials{
		login:    login,
		password: password,
	}
	c.metaData = newMetaData()
	return &c
}

func (c credentials) Value() string {
	return fmt.Sprintf(credentialsFormat, c.login, c.password)
}

func (c *credentials) SetValue(value string) error {
	values := strings.Split(value, credentialssep)
	if len(values) != 2 {
		return fmt.Errorf("%w: %q", ErrCredsInvalidFormat, value)
	}
	c.login = values[0]
	c.password = values[1]

	return nil
}

func (d *Data) SetCredentials(value string) error {
	d.editNow()
	values := strings.Split(value, credentialssep)
	if len(values) != 2 {
		return fmt.Errorf("%w: %q", ErrCredsInvalidFormat, value)
	}
	d.Value = value
	return nil
}

func (d *Data) CredentialsValues() (login string, password string) {
	fmt.Sscanf(d.Value, credentialsFormat, login, password)
	return
}
