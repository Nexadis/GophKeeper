package datas

type text struct {
	metaData
	Data string
}

func (t text) Type() DataType {
	return TextType
}

func NewText() text {
	t := text{}
	t.metaData = newMetaData()
	return t
}
