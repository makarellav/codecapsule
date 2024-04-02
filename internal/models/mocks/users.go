package mocks

import (
	"github.com/makarellav/codecapsule/internal/models"
	"time"
)

var mockUser = models.User{
	ID:      1,
	Name:    "Alice Doe",
	Email:   "alice@example.com",
	Created: time.Now(),
}

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

func (um *UserModel) Get(id int) (*models.User, error) {
	switch id {
	case 1:
		return &mockUser, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (um *UserModel) UpdatePassword(id int, currentPassword, newPassword string) error {
	if id == 1 {
		if currentPassword != "pa$$word" {
			return models.ErrInvalidCredentials
		}

		return nil
	}

	return models.ErrNoRecord
}
