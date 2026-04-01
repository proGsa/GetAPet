package purchaserequest

import (
	"encoding/json"
	"errors"
	"net/http"

	"getapet-backend/internal/dto"
	"getapet-backend/internal/models"
	"github.com/google/uuid"
)

func writePurchaseRequestMappedError(w http.ResponseWriter, err error, defaultMessage string) {
	switch {
	case errors.Is(err, models.ErrPurchaseRequestNotFound):
		writeErrorResponse(w, http.StatusNotFound, err, "Заявка на покупку не найдена")
	case errors.Is(err, models.ErrPurchaseRequestForbidden):
		writeErrorResponse(w, http.StatusForbidden, err, "Заявка принадлежит другому пользователю, поэтому у вас нет прав на это действие")
	case errors.Is(err, models.ErrPetNotFound):
		writeErrorResponse(w, http.StatusNotFound, err, "Питомец не найден")
	case errors.Is(err, models.ErrPurchaseRequestStatusRequired):
		writeErrorResponse(w, http.StatusBadRequest, err, "Необходимо передать значение статуса")
	case errors.Is(err, models.ErrPurchaseRequestPetNotAvailable):
		writeErrorResponse(w, http.StatusConflict, err, "Питомец недоступен для покупки")
	case errors.Is(err, models.ErrPurchaseRequestDuplicatePetBuyer):
		writeErrorResponse(w, http.StatusConflict, err, "Повторная заявка покупателя на одного и того же питомца запрещена")
	case errors.Is(err, models.ErrPurchaseRequestAlreadyApprovedForPet):
		writeErrorResponse(w, http.StatusConflict, err, "Продавец уже одобрил заявку на этого питомца другому покупателю")
	case errors.Is(err, models.ErrPurchaseRequestUniqueViolation):
		writeErrorResponse(w, http.StatusConflict, err, "Операция невозможна: запись с такими данными уже существует")
	default:
		writeErrorResponse(w, http.StatusInternalServerError, err, defaultMessage)
	}
}

// CreatePurchaseRequest godoc
// @Summary Create purchase request
// @Tags purchase-requests
// @Accept json
// @Produce json
// @Param request body dto.CreatePurchaseRequest true "Purchase request payload"
// @Security BearerAuth
// @Success 201 {object} dto.CreatePurchaseRequestResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 503 {object} ErrorResponse
// @Router /api/purchase-requests [post]
func (pr *PurchaseRequestRouter) CreatePurchaseRequest(w http.ResponseWriter, r *http.Request) {
	if pr.PurchaseRequestUsecase == nil {
		writeServiceUnavailable(w)
		return
	}

	var payload dto.CreatePurchaseRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Неверный формат JSON")
		return
	}
	if payload.PetID == uuid.Nil {
		writeErrorResponse(w, http.StatusBadRequest, errors.New("pet_id is required"), "Некорректный pet_id")
		return
	}
	if payload.BuyerID == uuid.Nil {
		writeErrorResponse(w, http.StatusBadRequest, errors.New("buyer_id is required"), "Некорректный buyer_id")
		return
	}

	createPurchaseRequest := dto.CreatePurchaseRequestFromDTO(payload)
	createdPurchaseRequest, err := pr.PurchaseRequestUsecase.Create(&createPurchaseRequest)
	if err != nil {
		writePurchaseRequestMappedError(w, err, "Не удалось создать заявку на покупку")
		return
	}

	writeSuccessResponse(w, http.StatusCreated, dto.PurchaseRequestToCreateDto(*createdPurchaseRequest))
}

// GetPurchaseRequests godoc
// @Summary Get all purchase requests
// @Tags purchase-requests
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.GetPurchaseRequestResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 503 {object} ErrorResponse
// @Router /api/purchase-requests [get]
func (pr *PurchaseRequestRouter) GetPurchaseRequests(w http.ResponseWriter, _ *http.Request) {
	if pr.PurchaseRequestUsecase == nil {
		writeServiceUnavailable(w)
		return
	}

	requests, err := pr.PurchaseRequestUsecase.GetAll()
	if err != nil {
		writePurchaseRequestMappedError(w, err, "Не удалось получить заявки на покупку")
		return
	}

	writeSuccessResponse(w, http.StatusOK, dto.PurchaseRequestsToGetDtos(requests))
}

// GetPurchaseRequestsByBuyer godoc
// @Summary Get purchase requests by buyer ID
// @Tags purchase-requests
// @Produce json
// @Param id path string true "Buyer ID (UUID)"
// @Security BearerAuth
// @Success 200 {array} dto.GetPurchaseRequestResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 503 {object} ErrorResponse
// @Router /api/purchase-requests/buyer/{id} [get]
func (pr *PurchaseRequestRouter) GetPurchaseRequestsByBuyer(w http.ResponseWriter, r *http.Request) {
	if pr.PurchaseRequestUsecase == nil {
		writeServiceUnavailable(w)
		return
	}

	buyerID, err := parseIDFromPath(r)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Некорректный buyer id")
		return
	}

	requests, err := pr.PurchaseRequestUsecase.GetByBuyerID(buyerID)
	if err != nil {
		writePurchaseRequestMappedError(w, err, "Не удалось получить заявки покупателя")
		return
	}

	writeSuccessResponse(w, http.StatusOK, dto.PurchaseRequestsToGetDtos(requests))
}

// GetPurchaseRequestsBySeller godoc
// @Summary Get purchase requests by seller ID
// @Tags purchase-requests
// @Produce json
// @Param id path string true "Seller ID (UUID)"
// @Security BearerAuth
// @Success 200 {array} dto.GetPurchaseRequestResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 503 {object} ErrorResponse
// @Router /api/purchase-requests/seller/{id} [get]
func (pr *PurchaseRequestRouter) GetPurchaseRequestsBySeller(w http.ResponseWriter, r *http.Request) {
	if pr.PurchaseRequestUsecase == nil {
		writeServiceUnavailable(w)
		return
	}

	sellerID, err := parseIDFromPath(r)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Некорректный seller id")
		return
	}

	requests, err := pr.PurchaseRequestUsecase.GetBySellerID(sellerID)
	if err != nil {
		writePurchaseRequestMappedError(w, err, "Не удалось получить заявки продавца")
		return
	}

	writeSuccessResponse(w, http.StatusOK, dto.PurchaseRequestsToGetDtos(requests))
}

// GetPurchaseRequestsByPet godoc
// @Summary Get purchase requests by pet ID
// @Tags purchase-requests
// @Produce json
// @Param pet_id path string true "Pet ID (UUID)"
// @Security BearerAuth
// @Success 200 {array} dto.GetPurchaseRequestResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 503 {object} ErrorResponse
// @Router /api/purchase-requests/pet/{pet_id} [get]
func (pr *PurchaseRequestRouter) GetPurchaseRequestsByPet(w http.ResponseWriter, r *http.Request) {
	if pr.PurchaseRequestUsecase == nil {
		writeServiceUnavailable(w)
		return
	}

	petID, err := parseIDFromPathParam(r, "pet_id")
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Некорректный pet_id")
		return
	}

	requests, err := pr.PurchaseRequestUsecase.GetByPetID(petID)
	if err != nil {
		writePurchaseRequestMappedError(w, err, "Не удалось получить заявки по питомцу")
		return
	}

	writeSuccessResponse(w, http.StatusOK, dto.PurchaseRequestsToGetDtos(requests))
}

// GetPurchaseRequest godoc
// @Summary Get purchase request by ID
// @Tags purchase-requests
// @Produce json
// @Param id path string true "Purchase request ID (UUID)"
// @Security BearerAuth
// @Success 200 {object} dto.GetPurchaseRequestResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 503 {object} ErrorResponse
// @Router /api/purchase-requests/{id} [get]
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
		writePurchaseRequestMappedError(w, err, "Не удалось получить заявку на покупку")
		return
	}

	writeSuccessResponse(w, http.StatusOK, dto.PurchaseRequestToGetDto(*request))
}

// UpdatePurchaseRequestStatus godoc
// @Summary Update purchase request status
// @Tags purchase-requests
// @Accept json
// @Produce json
// @Param id path string true "Purchase request ID (UUID)"
// @Param request body dto.UpdatePurchaseRequestStatus true "Status payload"
// @Security BearerAuth
// @Success 200 {object} dto.UpdatePurchaseRequestStatusResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 503 {object} ErrorResponse
// @Router /api/purchase-requests/{id}/status [patch]
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

	var payload dto.UpdatePurchaseRequestStatus
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Неверный формат JSON")
		return
	}

	// sellerID, err := userIDFromContext(r)
	// if err != nil {
	// 	writeErrorResponse(w, http.StatusUnauthorized, err, "Некорректный user_id в токене")
	// 	return
	// }

	//updatedRequest, err := pr.PurchaseRequestUsecase.UpdateStatus(id, sellerID, payload.Status)
	updatedRequest, err := pr.PurchaseRequestUsecase.UpdateStatus(id, uuid.Nil, payload.Status)
	if err != nil {
		writePurchaseRequestMappedError(w, err, "Не удалось обновить статус заявки")
		return
	}

	writeSuccessResponse(w, http.StatusOK, dto.PurchaseRequestToUpdateStatusDto(*updatedRequest))
}

// DeletePurchaseRequest godoc
// @Summary Delete purchase request
// @Tags purchase-requests
// @Produce json
// @Param id path string true "Purchase request ID (UUID)"
// @Security BearerAuth
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 503 {object} ErrorResponse
// @Router /api/purchase-requests/{id} [delete]
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

	// buyerID, err := userIDFromContext(r)
	// if err != nil {
	// 	writeErrorResponse(w, http.StatusUnauthorized, err, "Некорректный user_id в токене")
	// 	return
	// }

	//err = pr.PurchaseRequestUsecase.Delete(id, buyerID)
	err = pr.PurchaseRequestUsecase.Delete(id, uuid.Nil)
	if err != nil {
		writePurchaseRequestMappedError(w, err, "Не удалось удалить заявку на покупку")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
