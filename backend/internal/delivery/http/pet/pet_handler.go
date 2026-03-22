package pet

import (
	"encoding/json"
	"net/http"

	"getapet-backend/internal/models"
)

func (pr *PetRouter) CreatePet(w http.ResponseWriter, r *http.Request) {
	if pr.PetUsecase == nil {
		writeServiceUnavailable(w)
		return
	}

	var createPet models.Pet
	if err := json.NewDecoder(r.Body).Decode(&createPet); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Неверный формат JSON")
		return
	}

	createdPet, err := pr.PetUsecase.Create(&createPet)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось создать питомца")
		return
	}

	writeSuccessResponse(w, http.StatusCreated, createdPet)
}

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

	writeSuccessResponse(w, http.StatusOK, pets)
}

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

	writeSuccessResponse(w, http.StatusOK, pet)
}

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

	var updatePet models.Pet
	if err := json.NewDecoder(r.Body).Decode(&updatePet); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Неверный формат JSON")
		return
	}

	// sellerID will be taken from auth middleware later.
	updatedPet, err := pr.PetUsecase.Update(id, 0, &updatePet)
	if err != nil {
		if err == models.ErrPetNotFound {
			writeErrorResponse(w, http.StatusNotFound, err, "Питомец не найден")
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось обновить питомца")
		return
	}

	writeSuccessResponse(w, http.StatusOK, updatedPet)
}

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

	// sellerID will be taken from auth middleware later.
	err = pr.PetUsecase.Delete(id, 0)
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
