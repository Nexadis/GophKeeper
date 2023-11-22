package httpserver

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type BankCard struct {
	ID          int    `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
	Number      string `json:"number"`
	CardHolder  string `json:"cardholder"`
	Expire      string `json:"expire"`
	CVV         int    `json:"cvv"`
}

type Credentials struct {
	ID          int    `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
	Login       string `json:"login"`
	Password    string `json:"password"`
}

type Text struct {
	ID          int    `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
	Text        string `json:"text"`
}

type Binary struct {
	ID          int    `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
	Data        []byte `json:"data"`
}
