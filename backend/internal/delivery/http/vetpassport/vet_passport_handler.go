package vetpassport

import (
	"encoding/json"
	"net/http"

	"getapet-backend/internal/dto"
	"getapet-backend/internal/models"
	// "github.com/google/uuid"
)

// CreateVetPassport godoc
// @Summary Создать ветпаспорт
// @Tags vet-passports
// @Accept json
// @Produce json
// @Param passport body dto.CreateVetPassportRequest true "Данные ветпаспорта"
// @Success 201 {object} dto.VetPassportResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 503 {object} ErrorResponse
// @Router /api/vet-passports [post]
func (vr *VetPassportRouter) CreateVetPassport(w http.ResponseWriter, r *http.Request) {
	if vr.VetPassportUsecase == nil {
		writeServiceUnavailable(w)
		return
	}

	var req dto.CreateVetPassportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Неверный формат JSON")
		return
	}

	passport := dto.CreateVetPassportRequestFromDTO(req)

	createdVetPassport, err := vr.VetPassportUsecase.Create(&passport)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось создать ветпаспорт")
		return
	}

	writeSuccessResponse(w, http.StatusCreated, createdVetPassport)
}

// GetVetPassports godoc
// @Summary Получить список ветпаспортов
// @Tags vet-passports
// @Produce json
// @Success 200 {array} dto.VetPassportResponse
// @Failure 500 {object} ErrorResponse
// @Failure 503 {object} ErrorResponse
// @Router /api/vet-passports [get]
func (vr *VetPassportRouter) GetVetPassports(w http.ResponseWriter, _ *http.Request) {
	if vr.VetPassportUsecase == nil {
		writeServiceUnavailable(w)
		return
	}

	vetPassports, err := vr.VetPassportUsecase.GetAll()
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось получить ветпаспорта")
		return
	}

	writeSuccessResponse(w, http.StatusOK, vetPassports)
}

// GetVetPassport godoc
// @Summary Получить ветпаспорт по ID
// @Tags vet-passports
// @Produce json
// @Param id path string true "ID ветпаспорта"
// @Success 200 {object} dto.VetPassportResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 503 {object} ErrorResponse
// @Router /api/vet-passports/{id} [get]
func (vr *VetPassportRouter) GetVetPassport(w http.ResponseWriter, r *http.Request) {
	if vr.VetPassportUsecase == nil {
		writeServiceUnavailable(w)
		return
	}

	id, err := parseIDFromPath(r)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Некорректный id")
		return
	}

	vetPassport, err := vr.VetPassportUsecase.GetByID(id)
	if err != nil {
		if err == models.ErrVetPassportNotFound {
			writeErrorResponse(w, http.StatusNotFound, err, "Ветпаспорт не найден")
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось получить ветпаспорт")
		return
	}

	writeSuccessResponse(w, http.StatusOK, vetPassport)
}

// UpdateVetPassport godoc
// @Summary Обновить ветпаспорт
// @Tags vet-passports
// @Accept json
// @Produce json
// @Param id path string true "ID ветпаспорта"
// @Param passport body dto.UpdateVetPassportRequest true "Обновленные данные"
// @Success 200 {object} dto.VetPassportResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 503 {object} ErrorResponse
// @Router /api/vet-passports/{id} [put]
func (vr *VetPassportRouter) UpdateVetPassport(w http.ResponseWriter, r *http.Request) {
	if vr.VetPassportUsecase == nil {
		writeServiceUnavailable(w)
		return
	}

	id, err := parseIDFromPath(r)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Некорректный id")
		return
	}

	var req dto.UpdateVetPassportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Неверный формат JSON")
		return
	}

	passport := dto.UpdateVetPassportRequestFromDTO(req)

	updatedVetPassport, err := vr.VetPassportUsecase.Update(id, &passport)
	if err != nil {
		if err == models.ErrVetPassportNotFound {
			writeErrorResponse(w, http.StatusNotFound, err, "Ветпаспорт не найден")
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось обновить ветпаспорт")
		return
	}

	writeSuccessResponse(w, http.StatusOK, updatedVetPassport)
}

// DeleteVetPassport godoc
// @Summary Удалить ветпаспорт
// @Tags vet-passports
// @Param id path string true "ID ветпаспорта"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 503 {object} ErrorResponse
// @Router /api/vet-passports/{id} [delete]
func (vr *VetPassportRouter) DeleteVetPassport(w http.ResponseWriter, r *http.Request) {
	if vr.VetPassportUsecase == nil {
		writeServiceUnavailable(w)
		return
	}

	id, err := parseIDFromPath(r)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Некорректный id")
		return
	}

	err = vr.VetPassportUsecase.Delete(id)
	if err != nil {
		if err == models.ErrVetPassportNotFound {
			writeErrorResponse(w, http.StatusNotFound, err, "Ветпаспорт не найден")
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось удалить ветпаспорт")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
