package models

import (
	"errors"
	"time"
	"github.com/google/uuid"

)

type PurchaseRequest struct {
	ID          uuid.UUID       `json:"id" db:"id"`
	PetID       uuid.UUID       `json:"pet_id" db:"pet_id"`
	SellerID    uuid.UUID       `json:"seller_id" db:"seller_id"`
	Status      string    `json:"status" db:"status"`
	RequestDate time.Time `json:"request_date" db:"request_date"`
}

type PurchaseRequestRepository interface {
	Create(request *PurchaseRequest) (*PurchaseRequest, error)
	GetAll() ([]PurchaseRequest, error)
	GetByID(id uuid.UUID) (*PurchaseRequest, error)
	GetBySellerID(sellerID uuid.UUID) ([]PurchaseRequest, error)
	GetByPetID(petID uuid.UUID) ([]PurchaseRequest, error)
	UpdateStatus(id uuid.UUID, status string) (*PurchaseRequest, error)
	Delete(id uuid.UUID) error
}

type PurchaseRequestService interface {
	Create(request *PurchaseRequest) (*PurchaseRequest, error)
	GetAll() ([]PurchaseRequest, error)
	GetByID(id uuid.UUID) (*PurchaseRequest, error)
	GetBySellerID(sellerID uuid.UUID) ([]PurchaseRequest, error)
	GetByPetID(petID uuid.UUID) ([]PurchaseRequest, error)
	UpdateStatus(id uuid.UUID, sellerID uuid.UUID, status string) (*PurchaseRequest, error)
	Delete(id uuid.UUID, sellerID uuid.UUID) error
}

var (
	ErrPurchaseRequestNotFound = errors.New("purchase request not found")
)
