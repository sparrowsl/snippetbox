package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
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

	// Hash password with bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	_, err = model.DB.Exec(statement, name, email, hashedPassword)
	if err != nil {
		var mysqlError *mysql.MySQLError

		if errors.As(err, &mysqlError) {
			if mysqlError.Number == 1062 && strings.Contains(mysqlError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}

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
