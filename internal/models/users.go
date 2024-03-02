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

func (model *UserModel) Authenticate(email string, pass string) (int, error) {
	statement := `SELECT id, password FROM users WHERE email = ?`

	var id int
	var password []byte

	row := model.DB.QueryRow(statement, email)
	if err := row.Scan(&id, &password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	if err := bcrypt.CompareHashAndPassword(password, []byte(pass)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
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
