package user

type User struct {
	Id              int64  `json:"id"`
	UserDocument    string `json:"user_document"`
	CreditCardToken string `json:"credit_card_token"`
	Value           int64  `json:"value"`
}

func NewDefaultUser() *User {
	return &User{}
}

func NewUser(userDocument string, creditToken string, value int64) *User {
	return &User{
		UserDocument:    userDocument,
		CreditCardToken: creditToken,
		Value:           value,
	}
}
