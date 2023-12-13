package datas

type text struct {
	metaData
	data string
}

func (d *Data) SetText(value string) error {
	d.editNow()
	d.Value = value
	return nil
}

func NewText(data string) *text {
	t := text{
		data: data,
	}
	t.metaData = newMetaData()
	return &t
}

func (t text) Value() string {
	return t.data
}
