package usecase

import 
("getapet-backend/internal/models"
"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	userRepo models.UserRepository
}

func NewUserUsecase(userRepo models.UserRepository) *UserUsecase {
	return &UserUsecase{userRepo: userRepo}
}

func (u *UserUsecase) Create(user *models.User) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.UserPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.UserPassword = string(hashedPassword)

	return u.userRepo.Create(user)
}

func (u *UserUsecase) GetAll() ([]models.User, error) {
	return u.userRepo.GetAll()
}

func (u *UserUsecase) GetByID(id uuid.UUID) (*models.User, error) {
	return u.userRepo.GetByID(id)
}

func (u *UserUsecase) GetByLogin(login string) (*models.User, error) {
	return u.userRepo.GetByLogin(login)
}

func (u *UserUsecase) Update(id uuid.UUID, user *models.User) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.UserPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.UserPassword = string(hashedPassword)
	return u.userRepo.Update(id, user)
}

func (u *UserUsecase) Delete(id uuid.UUID) error {
	return u.userRepo.Delete(id)
}

func (u *UserUsecase) Login(login string, password string) (*models.User, error)  {
	user, err := u.userRepo.GetByLogin(login)
	if err != nil {
		if err == models.ErrUserNotFound {
			return nil,models.ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(password)); err != nil {
		return nil, models.ErrInvalidCredentials
	}

	return user, nil
}
