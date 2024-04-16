package user

type User struct {
	Id              int64
	UserDocument    string
	CreditCardToken string
	Value           int64
}

func NewDefaultUser() *User {
    return &User{}
}
