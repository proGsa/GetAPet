package models

import ("errors"
	"github.com/google/uuid"
)

// type UserDomain struct {
// 	ID              uuid.UUID    `json:"id" `
// 	FIO             string `json:"fio"`
// 	TelephoneNumber string `json:"telephone_number"`
// 	City            string `json:"city" `
// 	UserLogin       string `json:"user_login"`
// 	UserPassword    string `json:"user_password" `
// 	Status          string `json:"status"`
// 	UserDescription string `json:"user_description"`
// }

type User struct {
	ID              uuid.UUID    
	FIO             string 
	TelephoneNumber string 
	City            string
	UserLogin       string 
	UserPassword    string
	Status          string
	UserDescription string 
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
	Login(login string, password string) error
}

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrInvalidCredentials   = errors.New("invalid credentials")
)
