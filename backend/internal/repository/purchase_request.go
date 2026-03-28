package repository

import (
	"database/sql"
	"errors"

	"getapet-backend/internal/models"
	"github.com/google/uuid"
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

func (r *PurchaseRequestRepository) GetByID(_ uuid.UUID) (*models.PurchaseRequest, error) {
	return nil, errors.New("not implemented")
}

func (r *PurchaseRequestRepository) GetBySellerID(_ uuid.UUID) ([]models.PurchaseRequest, error) {
	return nil, errors.New("not implemented")
}

func (r *PurchaseRequestRepository) GetByPetID(_ uuid.UUID) ([]models.PurchaseRequest, error) {
	return nil, errors.New("not implemented")
}

func (r *PurchaseRequestRepository) UpdateStatus(_ uuid.UUID, _ string) (*models.PurchaseRequest, error) {
	return nil, errors.New("not implemented")
}

func (r *PurchaseRequestRepository) Delete(_ uuid.UUID) error {
	return errors.New("not implemented")
}
