package models

import (
	"errors"
	"time"
)

type PurchaseRequest struct {
	ID          int       `json:"id" db:"id"`
	PetID       int       `json:"pet_id" db:"pet_id"`
	SellerID    int       `json:"seller_id" db:"seller_id"`
	Status      string    `json:"status" db:"status"`
	RequestDate time.Time `json:"request_date" db:"request_date"`
}

type PurchaseRequestRepository interface {
	Create(request *PurchaseRequest) (*PurchaseRequest, error)
	GetAll() ([]PurchaseRequest, error)
	GetByID(id int) (*PurchaseRequest, error)
	GetBySellerID(sellerID int) ([]PurchaseRequest, error)
	GetByPetID(petID int) ([]PurchaseRequest, error)
	UpdateStatus(id int, status string) (*PurchaseRequest, error)
	Delete(id int) error
}

type PurchaseRequestService interface {
	Create(request *PurchaseRequest) (*PurchaseRequest, error)
	GetAll() ([]PurchaseRequest, error)
	GetByID(id int) (*PurchaseRequest, error)
	GetBySellerID(sellerID int) ([]PurchaseRequest, error)
	GetByPetID(petID int) ([]PurchaseRequest, error)
	UpdateStatus(id int, sellerID int, status string) (*PurchaseRequest, error)
	Delete(id int, sellerID int) error
}

var (
	ErrPurchaseRequestNotFound = errors.New("purchase request not found")
)
