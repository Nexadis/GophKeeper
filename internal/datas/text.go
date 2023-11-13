package datas

type text struct {
	metaData
	data string
}

func (t text) Type() DataType {
	return TextType
}

func NewText(data string) text {
	t := text{
		data: data,
	}
	t.metaData = newMetaData()
	return t
}

func (t text) Value() string {
	return t.data
}

func (t text) SetValue(value string) error {
	t.editNow()
	t.data = value
	return nil
}
