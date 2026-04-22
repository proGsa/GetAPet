package usecase

import (
	"getapet-backend/internal/models"
	"github.com/google/uuid"
)

type VetPassportUsecase struct {
	vetPassportRepo models.VetPassportRepository
}

func NewVetPassportUsecase(vetPassportRepo models.VetPassportRepository) *VetPassportUsecase {
	return &VetPassportUsecase{vetPassportRepo: vetPassportRepo}
}

func (u *VetPassportUsecase) Create(passport *models.VetPassport) (*models.VetPassport, error) {
	return u.vetPassportRepo.Create(passport)
}

func (u *VetPassportUsecase) GetAll() ([]models.VetPassport, error) {
	return u.vetPassportRepo.GetAll()
}

func (u *VetPassportUsecase) GetByID(id uuid.UUID) (*models.VetPassport, error) {
	return u.vetPassportRepo.GetByID(id)
}

func (u *VetPassportUsecase) Update(id uuid.UUID, passport *models.VetPassport) (*models.VetPassport, error) {
	return u.vetPassportRepo.Update(id, passport)
}

func (u *VetPassportUsecase) Delete(id uuid.UUID) error {
	return u.vetPassportRepo.Delete(id)
}
