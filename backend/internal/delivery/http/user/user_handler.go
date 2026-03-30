package user

import (
	"encoding/json"
	"net/http"
	"time"

	"getapet-backend/internal/dto"
	"getapet-backend/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

// CreateUser godoc
// @Summary Create user
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "User payload"
// @Success 201 {object} dto.RegisterResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 503 {object} ErrorResponse
// @Router /api/users [post]
func (ur *UserRouter) CreateUser(w http.ResponseWriter, r *http.Request) {
	if ur.UserUsecase == nil {
		writeServiceUnavailable(w)
		return
	}

	var createUser dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&createUser); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Неверный формат JSON")
		return
	}
	userDomain := dto.CreateUserRequestFromDTO(createUser)

	createdUser, err := ur.UserUsecase.Create(&userDomain)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось создать пользователя")
		return
	}

	writeSuccessResponse(w, http.StatusCreated, dto.RegisterResponseFromDomain(*createdUser))
}

// Login godoc
// @Summary Login user
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 503 {object} ErrorResponse
// @Router /api/users/login [post]
func (ur *UserRouter) Login(w http.ResponseWriter, r *http.Request) {
	if ur.UserUsecase == nil {
		writeServiceUnavailable(w)
		return
	}

	var creds dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Неверный формат JSON")
		return
	}

	if creds.UserLogin == "" || creds.UserPassword == "" {
		writeErrorResponse(w, http.StatusBadRequest, models.ErrInvalidCredentials, "Логин и пароль обязательны")
		return
	}
	user, err := ur.UserUsecase.Login(creds.UserLogin, creds.UserPassword)
	if err != nil {
		if err == models.ErrInvalidCredentials {
			writeErrorResponse(w, http.StatusUnauthorized, err, "Неверный логин или пароль")
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось выполнить вход")
		return
	}

	token, err := ur.generateJWT(*user)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось сгенерировать токен")
		return
	}

	writeSuccessResponse(w, http.StatusOK, dto.LoginResponseFromDomain(*user, token))
}

// // Logout godoc
// // @Summary Logout user
// // @Tags users
// // @Produce json
// // @Security BearerAuth
// // @Success 200 {object} dto.LogoutResponse
// // @Failure 401 {object} ErrorResponse
// // @Failure 503 {object} ErrorResponse
// // @Router /api/users/logout [post]
// func (ur *UserRouter) Logout(w http.ResponseWriter, _ *http.Request) {
// 	if ur.UserUsecase == nil {
// 		writeServiceUnavailable(w)
// 		return
// 	}
// 
// 	writeSuccessResponse(w, http.StatusOK, dto.LogoutResponse{Message: "Успешный выход из системы"})
// }

func (ur *UserRouter) generateJWT(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"login":   user.UserLogin,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(ur.JWTSecret))
}

// GetUsers godoc
// @Summary Get all users
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.UserResponse
// @Failure 500 {object} ErrorResponse
// @Failure 503 {object} ErrorResponse
// @Router /api/users [get]
func (ur *UserRouter) GetUsers(w http.ResponseWriter, _ *http.Request) {
	if ur.UserUsecase == nil {
		writeServiceUnavailable(w)
		return
	}

	users, err := ur.UserUsecase.GetAll()
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось получить пользователей")
		return
	}

	writeSuccessResponse(w, http.StatusOK, dto.UsersToDto(users))
}

// GetUser godoc
// @Summary Get user by ID
// @Tags users
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Security BearerAuth
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 503 {object} ErrorResponse
// @Router /api/users/{id} [get]
func (ur *UserRouter) GetUser(w http.ResponseWriter, r *http.Request) {
	if ur.UserUsecase == nil {
		writeServiceUnavailable(w)
		return
	}

	id, err := parseIDFromPath(r)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Некорректный id")
		return
	}

	user, err := ur.UserUsecase.GetByID(id)
	if err != nil {
		if err == models.ErrUserNotFound {
			writeErrorResponse(w, http.StatusNotFound, err, "Пользователь не найден")
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось получить пользователя")
		return
	}

	writeSuccessResponse(w, http.StatusOK, dto.UserToDto(*user))
}

// UpdateUser godoc
// @Summary Update user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Param user body dto.UpdateUserRequest true "User data"
// @Security BearerAuth
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 503 {object} ErrorResponse
// @Router /api/users/{id} [put]
func (ur *UserRouter) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if ur.UserUsecase == nil {
		writeServiceUnavailable(w)
		return
	}

	id, err := parseIDFromPath(r)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Некорректный id")
		return
	}

	var updateUser dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&updateUser); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Неверный формат JSON")
		return
	}
	domainUser := models.User{
	ID: id,
	FIO:             updateUser.FIO,
	TelephoneNumber: updateUser.TelephoneNumber,
	City:            updateUser.City,
	UserLogin:       updateUser.UserLogin,
	UserPassword:    updateUser.UserPassword,
	Status:          updateUser.Status,
	UserDescription: updateUser.UserDescription,
}
	updatedUser, err := ur.UserUsecase.Update(id, &domainUser)
	if err != nil {
		if err == models.ErrUserNotFound {
			writeErrorResponse(w, http.StatusNotFound, err, "Пользователь не найден")
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось обновить пользователя")
		return
	}

	writeSuccessResponse(w, http.StatusOK, dto.UserToDto(*updatedUser))
}

// DeleteUser godoc
// @Summary Delete user
// @Tags users
// @Param id path string true "User ID (UUID)"
// @Security BearerAuth
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 503 {object} ErrorResponse
// @Router /api/users/{id} [delete]
func (ur *UserRouter) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if ur.UserUsecase == nil {
		writeServiceUnavailable(w)
		return
	}

	id, err := parseIDFromPath(r)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Некорректный id")
		return
	}

	err = ur.UserUsecase.Delete(id)
	if err != nil {
		if err == models.ErrUserNotFound {
			writeErrorResponse(w, http.StatusNotFound, err, "Пользователь не найден")
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось удалить пользователя")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
