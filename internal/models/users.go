package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID       int
	Name     string
	Email    string
	Password []byte
	Created  time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (model *UserModel) Insert(name string, email string, password string) error {
	statement := `INSERT INTO users (name, email, password, created) VALUES (?, ?, ?, UTC_TIMESTAMP())`

	result, err := model.DB.Exec(statement, name, email, password)
	if err != nil {
		return err
	}

	result.LastInsertId()

	return nil
}

func (model *UserModel) Authenticate(email string, password string) (int, error) {
	return 0, nil
}

func (model *UserModel) Exists(id int) bool {
	statement := `SELECT id FROM users WHERE id = ?`
	result, err := model.DB.Exec(statement, id)
	if err != nil {
		return false
	}

	result.LastInsertId()

	return false
}
