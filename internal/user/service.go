package user

type repo interface {
	Create(user *User) error
	Read(id int64) (*User, error)
	Delete(id int64) error
	Update(id int64, user *User) (*User, error)
}

type service struct {
	repository repo
}

func NewService(rp repo) *service {
	return &service{rp}
}

func (svc *service) CreateUser(userDocument string, creditCardToken string, value int64) error {
	user := NewUser(userDocument, creditCardToken, value)
	return svc.repository.Create(user)
}

func (svc *service) GetUserByID(id int64) (*User, error) {
	return svc.repository.Read(id)
}

func (svc *service) DeleteUserById(id int64) error {
	return svc.repository.Delete(id)
}

func (svc *service) UpdateUserById(id int64, user *User) (*User, error) {
	return svc.repository.Update(id, user)
}
