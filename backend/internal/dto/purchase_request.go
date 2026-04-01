package dto

import (
	"getapet-backend/internal/models"
	"time"

	"github.com/google/uuid"
)

type CreatePurchaseRequest struct {
	PetID   uuid.UUID `json:"pet_id"`
	BuyerID uuid.UUID `json:"buyer_id"`
}

type CreatePurchaseRequestResponse struct {
	ID          uuid.UUID `json:"id"`
	PetID       uuid.UUID `json:"pet_id"`
	BuyerID     uuid.UUID `json:"buyer_id"`
	Status      string    `json:"status"`
	RequestDate time.Time `json:"request_date"`
}

type GetPurchaseRequestResponse struct {
	ID          uuid.UUID `json:"id"`
	PetID       uuid.UUID `json:"pet_id"`
	BuyerID     uuid.UUID `json:"buyer_id"`
	Status      string    `json:"status"`
	RequestDate time.Time `json:"request_date"`
}

type UpdatePurchaseRequestStatus struct {
	Status string `json:"status"`
}

type UpdatePurchaseRequestStatusResponse struct {
	ID          uuid.UUID `json:"id"`
	PetID       uuid.UUID `json:"pet_id"`
	BuyerID     uuid.UUID `json:"buyer_id"`
	Status      string    `json:"status"`
	RequestDate time.Time `json:"request_date"`
}

type DeletePurchaseRequestResponse struct {
	Message string `json:"message"`
}

func CreatePurchaseRequestFromDTO(req CreatePurchaseRequest) models.PurchaseRequest {
	return models.PurchaseRequest{
		PetID:   req.PetID,
		BuyerID: req.BuyerID,
	}
}

func PurchaseRequestToCreateDto(req models.PurchaseRequest) CreatePurchaseRequestResponse {
	return CreatePurchaseRequestResponse{
		ID:          req.ID,
		PetID:       req.PetID,
		BuyerID:     req.BuyerID,
		Status:      req.Status,
		RequestDate: req.RequestDate,
	}
}

func PurchaseRequestToGetDto(req models.PurchaseRequest) GetPurchaseRequestResponse {
	return GetPurchaseRequestResponse{
		ID:          req.ID,
		PetID:       req.PetID,
		BuyerID:     req.BuyerID,
		Status:      req.Status,
		RequestDate: req.RequestDate,
	}
}

func PurchaseRequestsToGetDtos(requests []models.PurchaseRequest) []GetPurchaseRequestResponse {
	response := make([]GetPurchaseRequestResponse, 0, len(requests))
	for _, request := range requests {
		response = append(response, GetPurchaseRequestResponse{
			ID:          request.ID,
			PetID:       request.PetID,
			BuyerID:     request.BuyerID,
			Status:      request.Status,
			RequestDate: request.RequestDate,
		})
	}
	return response
}

func PurchaseRequestToUpdateStatusDto(req models.PurchaseRequest) UpdatePurchaseRequestStatusResponse {
	return UpdatePurchaseRequestStatusResponse{
		ID:          req.ID,
		PetID:       req.PetID,
		BuyerID:     req.BuyerID,
		Status:      req.Status,
		RequestDate: req.RequestDate,
	}
}
