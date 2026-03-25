package user

import (
	"encoding/json"
	"net/http"

	"getapet-backend/internal/models"
)

func (ur *UserRouter) CreateUser(w http.ResponseWriter, r *http.Request) {
	if ur.UserUsecase == nil {
		writeServiceUnavailable(w)
		return
	}

	var createUser models.User
	if err := json.NewDecoder(r.Body).Decode(&createUser); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Неверный формат JSON")
		return
	}

	createdUser, err := ur.UserUsecase.Create(&createUser)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось создать пользователя")
		return
	}

	writeSuccessResponse(w, http.StatusCreated, createdUser)
}

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

	writeSuccessResponse(w, http.StatusOK, users)
}

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

	writeSuccessResponse(w, http.StatusOK, user)
}

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

	var updateUser models.User
	if err := json.NewDecoder(r.Body).Decode(&updateUser); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Неверный формат JSON")
		return
	}

	updatedUser, err := ur.UserUsecase.Update(id, &updateUser)
	if err != nil {
		if err == models.ErrUserNotFound {
			writeErrorResponse(w, http.StatusNotFound, err, "Пользователь не найден")
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось обновить пользователя")
		return
	}

	writeSuccessResponse(w, http.StatusOK, updatedUser)
}

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
