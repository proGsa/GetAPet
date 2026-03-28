package repository

import (
	"database/sql"
	"errors"

	"getapet-backend/internal/models"
	"github.com/google/uuid"
)

type VetPassportRepository struct {
	db *sql.DB
}

func NewVetPassportRepository(db *sql.DB) *VetPassportRepository {
	return &VetPassportRepository{db: db}
}

func (r *VetPassportRepository) Create(_ *models.VetPassport) (*models.VetPassport, error) {
	return nil, errors.New("not implemented")
}

func (r *VetPassportRepository) GetAll() ([]models.VetPassport, error) {
	return nil, errors.New("not implemented")
}

func (r *VetPassportRepository) GetByID(_ uuid.UUID) (*models.VetPassport, error) {
	return nil, errors.New("not implemented")
}

func (r *VetPassportRepository) Update(_ uuid.UUID, _ *models.VetPassport) (*models.VetPassport, error) {
	return nil, errors.New("not implemented")
}

func (r *VetPassportRepository) Delete(_ uuid.UUID) error {
	return errors.New("not implemented")
}
