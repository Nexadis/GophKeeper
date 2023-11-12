package datas

type Binary struct {
	metaData
	Data []byte
}

func (b Binary) Type() DataType {
	return BinaryType
}
