package datas

type BankCard struct {
	metaData
	Number string
}

func (bk BankCard) Type() DataType {
	return BinaryType
}
