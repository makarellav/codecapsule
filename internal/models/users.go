package models

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (um *UserModel) Insert(name, email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	stmt := `INSERT INTO users(name, email, hashed_password, created) 
	VALUES (?, ?, ?, UTC_TIMESTAMP())`

	args := []any{name, email, hash}

	_, err = um.DB.Exec(stmt, args...)

	if err != nil {
		var mysqlErr *mysql.MySQLError

		switch {
		case errors.As(err, &mysqlErr):
			if mysqlErr.Number == 1062 {
				return ErrDuplicateEmail
			}
		default:
			return err
		}
	}

	return nil
}

func (um *UserModel) Authenticate(email, password string) (int, error) {
	stmt := `SELECT id, hashed_password FROM users WHERE email = ?`

	var id int
	var hash []byte

	err := um.DB.QueryRow(stmt, email).Scan(&id, &hash)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return 0, ErrInvalidCredentials
		default:
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hash, []byte(password))

	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return 0, ErrInvalidCredentials
		default:
			return 0, err
		}
	}

	return id, nil
}

func (um *UserModel) Exists(id int) (bool, error) {
	stmt := `SELECT EXISTS(SELECT true FROM users WHERE id = ?)`

	var exists bool

	err := um.DB.QueryRow(stmt, id).Scan(&exists)

	return exists, err
}

func (um *UserModel) Get(id int) (*User, error) {
	stmt := `SELECT id, name, email, created FROM users WHERE id = ?`

	var user User

	err := um.DB.QueryRow(stmt, id).Scan(&user.ID, &user.Name, &user.Email, &user.Created)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRecord
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (um *UserModel) UpdatePassword(id int, currentPassword, newPassword string) error {
	stmt := `SELECT hashed_password FROM users WHERE id = ?`

	var hash []byte

	err := um.DB.QueryRow(stmt, id).Scan(&hash)

	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(hash, []byte(currentPassword))

	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return ErrInvalidCredentials
		default:
			return err
		}
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	stmt = `UPDATE users SET hashed_password = ? WHERE id = ?`

	args := []any{newHash, id}

	_, err = um.DB.Exec(stmt, args...)

	return err
}
