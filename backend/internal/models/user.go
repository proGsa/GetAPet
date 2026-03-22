package models

import "errors"

type User struct {
	ID              int    `json:"id" db:"id"`
	FIO             string `json:"fio" db:"fio"`
	TelephoneNumber string `json:"telephone_number" db:"telephone_number"`
	City            string `json:"city" db:"city"`
	UserLogin       string `json:"user_login" db:"user_login"`
	UserPassword    string `json:"user_password" db:"user_password"`
	Status          string `json:"status" db:"status"`
	UserDescription string `json:"user_description" db:"user_description"`
}

type UserRepository interface {
	Create(user *User) (*User, error)
	GetAll() ([]User, error)
	GetByID(id int) (*User, error)
	GetByLogin(login string) (*User, error)
	Update(id int, user *User) (*User, error)
	Delete(id int) error
}

type UserService interface {
	Create(user *User) (*User, error)
	GetAll() ([]User, error)
	GetByID(id int) (*User, error)
	GetByLogin(login string) (*User, error)
	Update(id int, user *User) (*User, error)
	Delete(id int) error
}

var (
	ErrUserNotFound = errors.New("user not found")
)
