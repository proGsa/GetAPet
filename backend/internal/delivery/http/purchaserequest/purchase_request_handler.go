package purchaserequest

import (
	"encoding/json"
	"net/http"

	"getapet-backend/internal/dto"
	"getapet-backend/internal/models"
)

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
		writeErrorResponse(w, http.StatusBadRequest, err, "Invalid JSON format")
		return
	}

	buyerID, err := userIDFromContext(r)
	if err != nil {
		writeErrorResponse(w, http.StatusUnauthorized, err, "Invalid user_id in token")
		return
	}

	createPurchaseRequest := dto.CreatePurchaseRequestFromDTO(payload, buyerID)
	createdPurchaseRequest, err := pr.PurchaseRequestUsecase.Create(&createPurchaseRequest)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err, "Failed to create purchase request")
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
		writeErrorResponse(w, http.StatusInternalServerError, err, "Failed to get purchase requests")
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
		writeErrorResponse(w, http.StatusBadRequest, err, "Invalid buyer id")
		return
	}

	requests, err := pr.PurchaseRequestUsecase.GetByBuyerID(buyerID)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err, "Failed to get buyer requests")
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
		writeErrorResponse(w, http.StatusBadRequest, err, "Invalid seller id")
		return
	}

	requests, err := pr.PurchaseRequestUsecase.GetBySellerID(sellerID)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err, "Failed to get seller requests")
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
		writeErrorResponse(w, http.StatusBadRequest, err, "Invalid pet_id")
		return
	}

	requests, err := pr.PurchaseRequestUsecase.GetByPetID(petID)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err, "Failed to get pet requests")
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
		writeErrorResponse(w, http.StatusBadRequest, err, "Invalid id")
		return
	}

	request, err := pr.PurchaseRequestUsecase.GetByID(id)
	if err != nil {
		if err == models.ErrPurchaseRequestNotFound {
			writeErrorResponse(w, http.StatusNotFound, err, "Purchase request not found")
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, err, "Failed to get purchase request")
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
		writeErrorResponse(w, http.StatusBadRequest, err, "Invalid id")
		return
	}

	var payload dto.UpdatePurchaseRequestStatus
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err, "Invalid JSON format")
		return
	}

	sellerID, err := userIDFromContext(r)
	if err != nil {
		writeErrorResponse(w, http.StatusUnauthorized, err, "Invalid user_id in token")
		return
	}

	updatedRequest, err := pr.PurchaseRequestUsecase.UpdateStatus(id, sellerID, payload.Status)
	if err != nil {
		if err == models.ErrPurchaseRequestNotFound {
			writeErrorResponse(w, http.StatusNotFound, err, "Purchase request not found")
			return
		}
		if err == models.ErrPurchaseRequestForbidden {
			writeErrorResponse(w, http.StatusForbidden, err, "Not enough permissions to update request")
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, err, "Failed to update purchase request status")
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
		writeErrorResponse(w, http.StatusBadRequest, err, "Invalid id")
		return
	}

	buyerID, err := userIDFromContext(r)
	if err != nil {
		writeErrorResponse(w, http.StatusUnauthorized, err, "Invalid user_id in token")
		return
	}

	err = pr.PurchaseRequestUsecase.Delete(id, buyerID)
	if err != nil {
		if err == models.ErrPurchaseRequestNotFound {
			writeErrorResponse(w, http.StatusNotFound, err, "Purchase request not found")
			return
		}
		if err == models.ErrPurchaseRequestForbidden {
			writeErrorResponse(w, http.StatusForbidden, err, "Not enough permissions to delete request")
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, err, "Failed to delete purchase request")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
