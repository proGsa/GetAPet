package usecase

import ("getapet-backend/internal/models"
"github.com/google/uuid"
)

type PurchaseRequestUsecase struct {
	purchaseRequestRepo models.PurchaseRequestRepository
}

func NewPurchaseRequestUsecase(purchaseRequestRepo models.PurchaseRequestRepository) *PurchaseRequestUsecase {
	return &PurchaseRequestUsecase{purchaseRequestRepo: purchaseRequestRepo}
}

func (u *PurchaseRequestUsecase) Create(request *models.PurchaseRequest) (*models.PurchaseRequest, error) {
	return u.purchaseRequestRepo.Create(request)
}

func (u *PurchaseRequestUsecase) GetAll() ([]models.PurchaseRequest, error) {
	return u.purchaseRequestRepo.GetAll()
}

func (u *PurchaseRequestUsecase) GetByID(id uuid.UUID) (*models.PurchaseRequest, error) {
	return u.purchaseRequestRepo.GetByID(id)
}

func (u *PurchaseRequestUsecase) GetBySellerID(sellerID uuid.UUID) ([]models.PurchaseRequest, error) {
	return u.purchaseRequestRepo.GetBySellerID(sellerID)
}

func (u *PurchaseRequestUsecase) GetByPetID(petID uuid.UUID) ([]models.PurchaseRequest, error) {
	return u.purchaseRequestRepo.GetByPetID(petID)
}

func (u *PurchaseRequestUsecase) UpdateStatus(id uuid.UUID, _ uuid.UUID, status string) (*models.PurchaseRequest, error) {
	return u.purchaseRequestRepo.UpdateStatus(id, status)
}

func (u *PurchaseRequestUsecase) Delete(id uuid.UUID, _ uuid.UUID) error {
	return u.purchaseRequestRepo.Delete(id)
}
