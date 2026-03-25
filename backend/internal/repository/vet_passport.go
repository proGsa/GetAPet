package repository

import (
	"database/sql"
	"errors"

	"getapet-backend/internal/models"
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

func (r *VetPassportRepository) GetByID(_ int) (*models.VetPassport, error) {
	return nil, errors.New("not implemented")
}

func (r *VetPassportRepository) Update(_ int, _ *models.VetPassport) (*models.VetPassport, error) {
	return nil, errors.New("not implemented")
}

func (r *VetPassportRepository) Delete(_ int) error {
	return errors.New("not implemented")
}
