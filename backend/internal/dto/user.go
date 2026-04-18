package dto

import (
	"getapet-backend/internal/models"
)

type LoginRequest struct {
	UserLogin    string `json:"user_login" validate:"required,min=3,max=50" example:"maria_petrova"`
	UserPassword string `json:"user_password" validate:"required,min=6,max=255" example:"securepassword123"`
}

type LoginResponse struct {
	Token       string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNTUwZTg0MDAtZTI5Yi00MWQ0LWE3MTYtNDQ2NjU1NDQwMDAwIiwiaXNfYWRtaW4iOmZhbHNlLCJleHAiOjE3MDAwMDAwMDB9"`
	ID      string `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type LogoutResponse struct {
	Message string `json:"message" example:"Успешный выход из системы"`
}

type CreateUserRequest struct {
	FIO             string `json:"fio" validate:"required,min=1,max=255" example:"Петрова Мария Сергеевна"`
	TelephoneNumber string `json:"telephone_number" validate:"required,min=5,max=20" example:"+79998887766"`
	City            string `json:"city" validate:"omitempty,max=50" example:"Saint Petersburg"`
	UserLogin       string `json:"user_login" validate:"required,min=3,max=50" example:"maria_petrova"`
	UserPassword    string `json:"user_password" validate:"required,min=6,max=255" example:"securepassword123"`
	Status          string `json:"status" validate:"omitempty,oneof=active blocked" example:"active"`
	UserDescription string `json:"user_description" validate:"omitempty,max=1000" example:"Волонтёр приюта"`
}

type RegisterResponse struct {
	ID string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type UserResponse struct {
	ID      string `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	FIO             string `json:"fio" example:"Иванов Иван Иванович"`
	TelephoneNumber string `json:"telephone_number" example:"+79991234567"`
	City            string `json:"city" example:"Moscow"`
	UserLogin       string `json:"user_login" example:"ivan_ivanov"`
	Status          string `json:"status" example:"active"`
	UserDescription string `json:"user_description" example:"Люблю животных"`
}

type UpdateUserRequest struct {
	FIO             string `json:"fio" validate:"required,min=1,max=255" example:"Петрова Мария Сергеевна"`
	TelephoneNumber string `json:"telephone_number" validate:"required,min=5,max=20" example:"+79998887766"`
	City            string `json:"city" validate:"omitempty,max=50" example:"Kazan"`
	UserLogin       string `json:"user_login" validate:"required,min=3,max=50" example:"maria_new_login"`
	UserPassword    string `json:"user_password" validate:"required,min=6,max=255" example:"newsecurepassword123"`
	Status          string `json:"status" validate:"required,oneof=active blocked" example:"active"`
	UserDescription string `json:"user_description" validate:"omitempty,max=1000" example:"Обновленное описание"`
}

type UpdateUserResponse struct {
	ID string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required,min=6,max=255" example:"oldpassword123"`
	NewPassword     string `json:"new_password" validate:"required,min=6,max=255" example:"newpassword123"`
}

// конверторы

func CreateUserRequestFromDTO(req CreateUserRequest) models.User {
	return models.User{
		FIO:             req.FIO,
		TelephoneNumber: req.TelephoneNumber,
		City:            req.City,
		UserLogin:       req.UserLogin,
		UserPassword:    req.UserPassword,
		Status:          req.Status,
		UserDescription: req.UserDescription,
	}
}

func UpdateUserRequestFromDTO(req UpdateUserRequest) models.User {
	return models.User{
		FIO:             req.FIO,
		TelephoneNumber: req.TelephoneNumber,
		City:            req.City,
		UserLogin:       req.UserLogin,
		UserPassword:    req.UserPassword,
		Status:          req.Status,
		UserDescription: req.UserDescription,
	}
}


func RegisterResponseFromDomain(u models.User) RegisterResponse {
	return RegisterResponse{
		ID: u.ID.String(),
	}
}

func LoginResponseFromDomain(u models.User, token string) LoginResponse {
	return LoginResponse{
		ID: u.ID.String(),
		Token: token,
	}
}

func UserToDto(domainUser models.User) UserResponse {
	return UserResponse{
		ID:             domainUser.ID.String(),
		FIO:             domainUser.FIO,
		TelephoneNumber: domainUser.TelephoneNumber,
		City:            domainUser.City,
		UserLogin:       domainUser.UserLogin,
		Status:          domainUser.Status,
		UserDescription: domainUser.UserDescription,
	}
}


func UsersToDto(domainUsers []models.User) []UserResponse {
	dtos := make([]UserResponse, len(domainUsers))
	for i, user := range domainUsers {
		dtos[i] = UserToDto(user)
	}
	return dtos
}
