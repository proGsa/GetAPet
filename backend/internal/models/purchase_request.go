package models

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type PurchaseRequest struct {
	ID          uuid.UUID `json:"id" db:"id"`
	PetID       uuid.UUID `json:"pet_id" db:"pet_id"`
	BuyerID     uuid.UUID `json:"buyer_id" db:"buyer_id"`
	SellerID    uuid.UUID `json:"seller_id"`
	Status      string    `json:"status" db:"status"`
	RequestDate time.Time `json:"request_date" db:"request_date"`
}

type PurchaseRequestRepository interface {
	Create(request *PurchaseRequest) (*PurchaseRequest, error)
	GetAll() ([]PurchaseRequest, error)
	GetByID(id uuid.UUID) (*PurchaseRequest, error)
	GetByBuyerID(buyerID uuid.UUID) ([]PurchaseRequest, error)
	GetBySellerID(sellerID uuid.UUID) ([]PurchaseRequest, error)
	GetByPetID(petID uuid.UUID) ([]PurchaseRequest, error)
	UpdateStatus(id uuid.UUID, status string) (*PurchaseRequest, error)
	UpdateStatusBySeller(id uuid.UUID, sellerID uuid.UUID, status string) (*PurchaseRequest, error)
	Delete(id uuid.UUID) error
	DeleteByBuyer(id uuid.UUID, buyerID uuid.UUID) error
}

type PurchaseRequestService interface {
	Create(request *PurchaseRequest) (*PurchaseRequest, error)
	GetAll() ([]PurchaseRequest, error)
	GetByID(id uuid.UUID) (*PurchaseRequest, error)
	GetByBuyerID(buyerID uuid.UUID) ([]PurchaseRequest, error)
	GetBySellerID(sellerID uuid.UUID) ([]PurchaseRequest, error)
	GetByPetID(petID uuid.UUID) ([]PurchaseRequest, error)
	UpdateStatus(id uuid.UUID, sellerID uuid.UUID, status string) (*PurchaseRequest, error)
	UpdateStatusBySeller(id uuid.UUID, sellerID uuid.UUID, status string) (*PurchaseRequest, error)
	Delete(id uuid.UUID, buyerID uuid.UUID) error
	DeleteByBuyer(id uuid.UUID, buyerID uuid.UUID) error
}

var (
	ErrPurchaseRequestNotFound              = errors.New("purchase request not found")
	ErrPurchaseRequestForbidden             = errors.New("purchase request access denied")
	ErrPurchaseRequestStatusRequired        = errors.New("purchase request status is required")
	ErrPurchaseRequestPetNotAvailable       = errors.New("pet is not available for purchase")
	ErrPurchaseRequestDuplicatePetBuyer     = errors.New("duplicate purchase request for pet and buyer")
	ErrPurchaseRequestAlreadyApprovedForPet = errors.New("another request for this pet is already approved")
	ErrPurchaseRequestUniqueViolation       = errors.New("purchase request unique constraint violated")
)
