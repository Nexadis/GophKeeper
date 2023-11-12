package datas

type Credentials struct {
	metaData
	Login    string
	Password string
}

func (c Credentials) Type() DataType {
	return CredentialsType
}
