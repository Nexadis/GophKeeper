package datas

import "fmt"

const credentialsFormat = "%s:%s"

var ErrCredsInvalidFormat = "credentials.SetValue: %w"

type credentials struct {
	metaData
	login    string
	password string
}

func (c credentials) Type() DataType {
	return CredentialsType
}

func NewCredentials(login string, password string) credentials {
	c := credentials{
		login:    login,
		password: password,
	}
	c.metaData = newMetaData()
	return c
}

func (c credentials) Value() string {
	return fmt.Sprintf(credentialsFormat, c.login, c.password)
}

func (c *credentials) SetValue(value string) error {
	c.editNow()
	var login, password string
	_, err := fmt.Sscanf(value, credentialsFormat, login, password)
	if err != nil {
		return fmt.Errorf(ErrCredsInvalidFormat, err)
	}
	c.login = login
	c.password = password

	return nil
}
