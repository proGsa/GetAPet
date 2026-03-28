package repository

import (
	"database/sql"
	"errors"

	"getapet-backend/internal/models"
	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(_ *models.User) (*models.User, error) {
	return nil, errors.New("not implemented")
}

func (r *UserRepository) GetAll() ([]models.User, error) {
	return nil, errors.New("not implemented")
}

func (r *UserRepository) GetByID(_ uuid.UUID) (*models.User, error) {
	return nil, errors.New("not implemented")
}

func (r *UserRepository) GetByLogin(_ string) (*models.User, error) {
	return nil, errors.New("not implemented")
}

func (r *UserRepository) Update(_ uuid.UUID, _ *models.User) (*models.User, error) {
	return nil, errors.New("not implemented")
}

func (r *UserRepository) Delete(_ uuid.UUID) error {
	return errors.New("not implemented")
}
