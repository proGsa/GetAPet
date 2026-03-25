package usecase

import "getapet-backend/internal/models"

type UserUsecase struct {
	userRepo models.UserRepository
}

func NewUserUsecase(userRepo models.UserRepository) *UserUsecase {
	return &UserUsecase{userRepo: userRepo}
}

func (u *UserUsecase) Create(user *models.User) (*models.User, error) {
	return u.userRepo.Create(user)
}

func (u *UserUsecase) GetAll() ([]models.User, error) {
	return u.userRepo.GetAll()
}

func (u *UserUsecase) GetByID(id int) (*models.User, error) {
	return u.userRepo.GetByID(id)
}

func (u *UserUsecase) GetByLogin(login string) (*models.User, error) {
	return u.userRepo.GetByLogin(login)
}

func (u *UserUsecase) Update(id int, user *models.User) (*models.User, error) {
	return u.userRepo.Update(id, user)
}

func (u *UserUsecase) Delete(id int) error {
	return u.userRepo.Delete(id)
}
