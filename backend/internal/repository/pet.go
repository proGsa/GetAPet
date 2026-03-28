package repository

import (
	"database/sql"
	"errors"

	"getapet-backend/internal/models"
	"github.com/google/uuid"
)

type PetRepository struct {
	db *sql.DB
}

func NewPetRepository(db *sql.DB) *PetRepository {
	return &PetRepository{db: db}
}

func (r *PetRepository) Create(_ *models.Pet) (*models.Pet, error) {
	return nil, errors.New("not implemented")
}

func (r *PetRepository) GetAll() ([]models.Pet, error) {
	return nil, errors.New("not implemented")
}

func (r *PetRepository) GetByID(_ uuid.UUID) (*models.Pet, error) {
	return nil, errors.New("not implemented")
}

func (r *PetRepository) GetBySellerID(_ uuid.UUID) ([]models.Pet, error) {
	return nil, errors.New("not implemented")
}

func (r *PetRepository) Update(_ uuid.UUID, _ *models.Pet) (*models.Pet, error) {
	return nil, errors.New("not implemented")
}

func (r *PetRepository) Delete(_ uuid.UUID) error {
	return errors.New("not implemented")
}

func (r *PetRepository) CheckBelonging(_, _ uuid.UUID) (bool, error) {
	return false, errors.New("not implemented")
}
