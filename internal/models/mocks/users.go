package mocks

import (
	"time"

	"github.com/sparrowsl/snippetbox/internal/models"
)

var mockUser = &models.User{
	ID:       1,
	Name:     "jack",
	Email:    "jack@mail.com",
	Password: []byte("password"),
	Created:  time.Now(),
}

type UserModel struct{}

func (mock *UserModel) Insert(name string, email string, password string) error {
	switch email {
	case "jack@mail.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (mock *UserModel) Authenticate(email string, password string) (int, error) {
	if email == "jack@mail.com" && password == "password" {
		return 1, nil
	}

	return 0, models.ErrInvalidCredentials
}

func (mock *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}
