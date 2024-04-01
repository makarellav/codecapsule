package mocks

import "github.com/makarellav/codecapsule/internal/models"

type UserModel struct{}

func (um *UserModel) Insert(name, email, password string) error {
	switch email {
	case "dupe@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (um *UserModel) Authenticate(email, password string) (int, error) {
	if email == "alice@example.com" && password == "pa$$word" {
		return 1, nil
	}

	return 0, models.ErrInvalidCredentials
}

func (um *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}
