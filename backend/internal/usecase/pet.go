package usecase

import "getapet-backend/internal/models"

type PetUsecase struct {
	petRepo models.PetRepository
}

func NewPetUsecase(petRepo models.PetRepository) *PetUsecase {
	return &PetUsecase{petRepo: petRepo}
}

func (u *PetUsecase) Create(pet *models.Pet) (*models.Pet, error) {
	return u.petRepo.Create(pet)
}

func (u *PetUsecase) GetAll() ([]models.Pet, error) {
	return u.petRepo.GetAll()
}

func (u *PetUsecase) GetByID(id int) (*models.Pet, error) {
	return u.petRepo.GetByID(id)
}

func (u *PetUsecase) GetBySellerID(sellerID int) ([]models.Pet, error) {
	return u.petRepo.GetBySellerID(sellerID)
}

func (u *PetUsecase) Update(id, sellerID int, pet *models.Pet) (*models.Pet, error) {
	ok, err := u.petRepo.CheckBelonging(id, sellerID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, models.ErrPetNotFound
	}
	return u.petRepo.Update(id, pet)
}

func (u *PetUsecase) Delete(id, sellerID int) error {
	ok, err := u.petRepo.CheckBelonging(id, sellerID)
	if err != nil {
		return err
	}
	if !ok {
		return models.ErrPetNotFound
	}
	return u.petRepo.Delete(id)
}
