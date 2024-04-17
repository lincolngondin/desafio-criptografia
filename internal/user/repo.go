package user

import (
	"database/sql"
	"errors"

	"github.com/lincolngondin/desafio-criptografia/internal/crypto"
)

var (
	ErrUserNotFound = errors.New("Id de usuario não encontrado!")
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{db}
}

func (repo *repository) Create(user *User) error {
	userDocumentEncrypted, encryptErr := crypto.Encrypt(user.UserDocument)
	if encryptErr != nil {
		return encryptErr
	}

	userCreditCardTokenEncrypted, encryptErr := crypto.Encrypt(user.CreditCardToken)
	if encryptErr != nil {
		return encryptErr
	}

	_, err := repo.db.Exec("INSERT INTO users (user_document, credit_card_token, value) VALUES ($1, $2, $3);", userDocumentEncrypted, userCreditCardTokenEncrypted, user.Value)
	return err
}

func (repo *repository) Read(id int64) (*User, error) {
	var userDocumentEncrypted []byte
	var creditCardTokenEncrypted []byte
	user := NewDefaultUser()
	row := repo.db.QueryRow("SELECT * FROM users WHERE id = $1;", id)
	err := row.Scan(&user.Id, &userDocumentEncrypted, &creditCardTokenEncrypted, &user.Value)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	user.UserDocument = crypto.Decrypt(userDocumentEncrypted)
	user.CreditCardToken = crypto.Decrypt(creditCardTokenEncrypted)
	return user, nil
}

func (repo *repository) Delete(id int64) error {
	tx, errTx := repo.db.Begin()
	if errTx != nil {
		return errTx
	}
	user := NewDefaultUser()
	var userDocumentEncrypted []byte
	var creditCardTokenEncrypted []byte
	row := tx.QueryRow("SELECT * FROM users WHERE id = $1;", id)
	err := row.Scan(&user.Id, &userDocumentEncrypted, &creditCardTokenEncrypted, &user.Value)
	if err == sql.ErrNoRows {
		tx.Rollback()
		return ErrUserNotFound
	}
	if err != nil {
		tx.Rollback()
		return err
	}

	_, errDelete := tx.Exec("DELETE FROM users WHERE id = $1;", id)
	if errDelete != nil {
		tx.Rollback()
		return errDelete
	}
	commitErr := tx.Commit()
	return commitErr
}

func (repo *repository) Update(id int64, user *User) (*User, error) {
	tx, errTx := repo.db.Begin()
	if errTx != nil {
		return nil, errTx
	}
	oldUser := NewDefaultUser()
	row := tx.QueryRow("SELECT * FROM users WHERE id = $1;", id)
	scanErr := row.Scan(&oldUser.Id, &oldUser.UserDocument, &oldUser.CreditCardToken, &oldUser.Value)
	if scanErr == sql.ErrNoRows {
		tx.Rollback()
		return nil, ErrUserNotFound
	}
	if scanErr != nil {
		tx.Rollback()
		return nil, scanErr
	}

	userDocument, encErr := crypto.Encrypt(user.UserDocument)
	if encErr != nil {
		tx.Rollback()
		return nil, encErr
	}
	creditToken, encErr := crypto.Encrypt(user.CreditCardToken)
	if encErr != nil {
		tx.Rollback()
		return nil, encErr
	}
	_, updateErr := tx.Exec("UPDATE users SET user_document = $1, credit_card_token = $2, value = $3 WHERE id = $4;", userDocument, creditToken, user.Value, id)
	if updateErr != nil {
		tx.Rollback()
		return nil, updateErr
	}
	commitErr := tx.Commit()
	if commitErr != nil {
		return nil, commitErr
	}

	nUser := NewUser(userDocument, creditToken, user.Value)
	nUser.Id = id
	return nUser, nil
}
