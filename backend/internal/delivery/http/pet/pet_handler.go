package pet

import (
	"encoding/json"
	"errors"
	"net/http"

	"getapet-backend/internal/delivery/middleware"
	"getapet-backend/internal/dto"
	"getapet-backend/internal/models"

	"github.com/google/uuid"
)

func getUserIDFromContext(r *http.Request) uuid.UUID {
	val := r.Context().Value(middleware.UserIDKey)
	if val == nil {
		return uuid.Nil
	}

	idStr, ok := val.(string)
	if !ok {
		return uuid.Nil
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil
	}

	return id
}

// CreatePet godoc
// @Summary Создать питомца
// @Tags pets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param pet body dto.CreatePetRequest true "Данные питомца"
// @Success 201 {object} dto.CreatePetResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 503 {object} ErrorResponse
// @Router /api/pets [post]
func (pr *PetRouter) CreatePet(w http.ResponseWriter, r *http.Request) {
	if pr.PetUsecase == nil {
		writeServiceUnavailable(w)
		return
	}

	sellerID := getUserIDFromContext(r)
	if sellerID == uuid.Nil {
		writeErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"), "Требуется авторизация")
		return
	}

	var req dto.CreatePetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Неверный формат JSON")
		return
	}

	createPet, err := dto.CreatePetRequestFromDTO(req, sellerID)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Некорректный vet_passport_id")
		return
	}

	createdPet, err := pr.PetUsecase.Create(&createPet)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось создать питомца")
		return
	}

	writeSuccessResponse(w, http.StatusCreated, dto.CreatePetResponseFromDomain(*createdPet))
}

// GetPets godoc
// @Summary Получить список питомцев
// @Tags pets
// @Produce json
// @Success 200 {array} dto.PetResponse
// @Failure 500 {object} ErrorResponse
// @Failure 503 {object} ErrorResponse
// @Router /api/pets [get]
func (pr *PetRouter) GetPets(w http.ResponseWriter, _ *http.Request) {
	if pr.PetUsecase == nil {
		writeServiceUnavailable(w)
		return
	}

	pets, err := pr.PetUsecase.GetAll()
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось получить питомцев")
		return
	}

	writeSuccessResponse(w, http.StatusOK, dto.PetsToDto(pets))
}

// GetPet godoc
// @Summary Получить питомца по ID
// @Tags pets
// @Produce json
// @Param id path string true "ID питомца"
// @Success 200 {object} dto.PetResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 503 {object} ErrorResponse
// @Router /api/pets/{id} [get]
func (pr *PetRouter) GetPet(w http.ResponseWriter, r *http.Request) {
	if pr.PetUsecase == nil {
		writeServiceUnavailable(w)
		return
	}

	id, err := parseIDFromPath(r)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Некорректный id")
		return
	}

	pet, err := pr.PetUsecase.GetByID(id)
	if err != nil {
		if err == models.ErrPetNotFound {
			writeErrorResponse(w, http.StatusNotFound, err, "Питомец не найден")
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось получить питомца")
		return
	}

	writeSuccessResponse(w, http.StatusOK, dto.PetToDto(*pet))
}

// UpdatePet godoc
// @Summary Обновить питомца
// @Tags pets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID питомца"
// @Param pet body dto.UpdatePetRequest true "Обновленные данные питомца"
// @Success 200 {object} dto.UpdatePetResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 503 {object} ErrorResponse
// @Router /api/pets/{id} [put]
func (pr *PetRouter) UpdatePet(w http.ResponseWriter, r *http.Request) {
	if pr.PetUsecase == nil {
		writeServiceUnavailable(w)
		return
	}

	id, err := parseIDFromPath(r)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Некорректный id")
		return
	}

	var req dto.UpdatePetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Неверный формат JSON")
		return
	}

	updatePet := dto.UpdatePetRequestToModel(req)

	// sellerID will be taken from auth middleware later.
	sellerID := getUserIDFromContext(r)
	updatedPet, err := pr.PetUsecase.Update(id, sellerID, &updatePet)
	if err != nil {
		if err == models.ErrPetNotFound {
			writeErrorResponse(w, http.StatusNotFound, err, "Питомец не найден")
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось обновить питомца")
		return
	}

	writeSuccessResponse(w, http.StatusOK, dto.UpdatePetResponseFromDomain(*updatedPet))
}

// DeletePet godoc
// @Summary Удалить питомца
// @Tags pets
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID питомца"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 503 {object} ErrorResponse
// @Router /api/pets/{id} [delete]
func (pr *PetRouter) DeletePet(w http.ResponseWriter, r *http.Request) {
	if pr.PetUsecase == nil {
		writeServiceUnavailable(w)
		return
	}

	id, err := parseIDFromPath(r)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Некорректный id")
		return
	}

	sellerID := getUserIDFromContext(r)
	if sellerID == uuid.Nil {
		writeErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"), "Требуется авторизация")
		return
	}

	err = pr.PetUsecase.Delete(id, sellerID)
	if err != nil {
		if err == models.ErrPetNotFound {
			writeErrorResponse(w, http.StatusNotFound, err, "Питомец не найден")
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось удалить питомца")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
