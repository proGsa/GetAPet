package purchaserequest

import (
	"encoding/json"
	"net/http"

	"getapet-backend/internal/models"
)

func (pr *PurchaseRequestRouter) CreatePurchaseRequest(w http.ResponseWriter, r *http.Request) {
	if pr.PurchaseRequestUsecase == nil {
		writeServiceUnavailable(w)
		return
	}

	var createPurchaseRequest models.PurchaseRequest
	if err := json.NewDecoder(r.Body).Decode(&createPurchaseRequest); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Неверный формат JSON")
		return
	}

	createdPurchaseRequest, err := pr.PurchaseRequestUsecase.Create(&createPurchaseRequest)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось создать заявку")
		return
	}

	writeSuccessResponse(w, http.StatusCreated, createdPurchaseRequest)
}

func (pr *PurchaseRequestRouter) GetPurchaseRequests(w http.ResponseWriter, _ *http.Request) {
	if pr.PurchaseRequestUsecase == nil {
		writeServiceUnavailable(w)
		return
	}

	requests, err := pr.PurchaseRequestUsecase.GetAll()
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось получить заявки")
		return
	}

	writeSuccessResponse(w, http.StatusOK, requests)
}

func (pr *PurchaseRequestRouter) GetPurchaseRequest(w http.ResponseWriter, r *http.Request) {
	if pr.PurchaseRequestUsecase == nil {
		writeServiceUnavailable(w)
		return
	}

	id, err := parseIDFromPath(r)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Некорректный id")
		return
	}

	request, err := pr.PurchaseRequestUsecase.GetByID(id)
	if err != nil {
		if err == models.ErrPurchaseRequestNotFound {
			writeErrorResponse(w, http.StatusNotFound, err, "Заявка не найдена")
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось получить заявку")
		return
	}

	writeSuccessResponse(w, http.StatusOK, request)
}

func (pr *PurchaseRequestRouter) UpdatePurchaseRequestStatus(w http.ResponseWriter, r *http.Request) {
	if pr.PurchaseRequestUsecase == nil {
		writeServiceUnavailable(w)
		return
	}

	id, err := parseIDFromPath(r)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Некорректный id")
		return
	}

	var payload struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Неверный формат JSON")
		return
	}

	// sellerID will be taken from auth middleware later.
	updatedRequest, err := pr.PurchaseRequestUsecase.UpdateStatus(id, 0, payload.Status)
	if err != nil {
		if err == models.ErrPurchaseRequestNotFound {
			writeErrorResponse(w, http.StatusNotFound, err, "Заявка не найдена")
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось обновить статус заявки")
		return
	}

	writeSuccessResponse(w, http.StatusOK, updatedRequest)
}

func (pr *PurchaseRequestRouter) DeletePurchaseRequest(w http.ResponseWriter, r *http.Request) {
	if pr.PurchaseRequestUsecase == nil {
		writeServiceUnavailable(w)
		return
	}

	id, err := parseIDFromPath(r)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Некорректный id")
		return
	}

	// sellerID will be taken from auth middleware later.
	err = pr.PurchaseRequestUsecase.Delete(id, 0)
	if err != nil {
		if err == models.ErrPurchaseRequestNotFound {
			writeErrorResponse(w, http.StatusNotFound, err, "Заявка не найдена")
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, err, "Не удалось удалить заявку")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
