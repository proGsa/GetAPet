package vetpassport

import (
	"encoding/json"
	"net/http"

	"getapet-backend/internal/models"
)

func (vr *VetPassportRouter) CreateVetPassport(w http.ResponseWriter, r *http.Request) {
	if vr.VetPassportUsecase == nil {
		writeServiceUnavailable(w)
		return
	}

	var createVetPassport models.VetPassport
	if err := json.NewDecoder(r.Body).Decode(&createVetPassport); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Неверный формат JSON")
		return
	}

	createdVetPassport, err := vr.VetPassportUsecase.Create(&createVetPassport)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось создать ветпаспорт")
		return
	}

	writeSuccessResponse(w, http.StatusCreated, createdVetPassport)
}

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

	var updateVetPassport models.VetPassport
	if err := json.NewDecoder(r.Body).Decode(&updateVetPassport); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Неверный формат JSON")
		return
	}

	updatedVetPassport, err := vr.VetPassportUsecase.Update(id, &updateVetPassport)
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
