package repository

import (
	"database/sql"
	"errors"

	"getapet-backend/internal/models"
)

type PurchaseRequestRepository struct {
	db *sql.DB
}

func NewPurchaseRequestRepository(db *sql.DB) *PurchaseRequestRepository {
	return &PurchaseRequestRepository{db: db}
}

func (r *PurchaseRequestRepository) Create(_ *models.PurchaseRequest) (*models.PurchaseRequest, error) {
	return nil, errors.New("not implemented")
}

func (r *PurchaseRequestRepository) GetAll() ([]models.PurchaseRequest, error) {
	return nil, errors.New("not implemented")
}

func (r *PurchaseRequestRepository) GetByID(_ int) (*models.PurchaseRequest, error) {
	return nil, errors.New("not implemented")
}

func (r *PurchaseRequestRepository) GetBySellerID(_ int) ([]models.PurchaseRequest, error) {
	return nil, errors.New("not implemented")
}

func (r *PurchaseRequestRepository) GetByPetID(_ int) ([]models.PurchaseRequest, error) {
	return nil, errors.New("not implemented")
}

func (r *PurchaseRequestRepository) UpdateStatus(_ int, _ string) (*models.PurchaseRequest, error) {
	return nil, errors.New("not implemented")
}

func (r *PurchaseRequestRepository) Delete(_ int) error {
	return errors.New("not implemented")
}
