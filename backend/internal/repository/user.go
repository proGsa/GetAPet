package repository

import (
	"database/sql"
	"errors"

	"getapet-backend/internal/models"
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

func (r *UserRepository) GetByID(_ int) (*models.User, error) {
	return nil, errors.New("not implemented")
}

func (r *UserRepository) GetByLogin(_ string) (*models.User, error) {
	return nil, errors.New("not implemented")
}

func (r *UserRepository) Update(_ int, _ *models.User) (*models.User, error) {
	return nil, errors.New("not implemented")
}

func (r *UserRepository) Delete(_ int) error {
	return errors.New("not implemented")
}
