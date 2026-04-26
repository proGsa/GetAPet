package models

import (
	"errors"
	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID `json:"id" db:"id"`
	FIO             string    `json:"fio" db:"fio"`
	TelephoneNumber string    `json:"telephone_number" db:"telephone_number"`
	City            string    `json:"city" db:"city"`
	UserLogin       string    `json:"user_login" db:"user_login"`
	UserPassword    string    `json:"user_password" db:"user_password"`
	Status          string    `json:"status" db:"status"`
	UserDescription string    `json:"user_description" db:"user_description"`
}

type UserRepository interface {
	Create(user *User) (*User, error)
	GetAll() ([]User, error)
	GetByID(id uuid.UUID) (*User, error)
	GetByLogin(login string) (*User, error)
	Update(id uuid.UUID, user *User) (*User, error)
	Delete(id uuid.UUID) error
}

type UserService interface {
	Create(user *User) (*User, error)
	GetAll() ([]User, error)
	GetByID(id uuid.UUID) (*User, error)
	GetByLogin(login string) (*User, error)
	Update(id uuid.UUID, user *User) (*User, error)
	Delete(id uuid.UUID) error
	Login(login string, password string) (*User, error)
}

var (
	ErrUserNotFound           = errors.New("user not found")
	ErrInvalidCredentials     = errors.New("invalid credentials")
	ErrUserLoginAlreadyExists = errors.New("user login already exists")
)

