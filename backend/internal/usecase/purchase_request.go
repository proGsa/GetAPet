package usecase

import "getapet-backend/internal/models"

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

func (u *PurchaseRequestUsecase) GetByID(id int) (*models.PurchaseRequest, error) {
	return u.purchaseRequestRepo.GetByID(id)
}

func (u *PurchaseRequestUsecase) GetBySellerID(sellerID int) ([]models.PurchaseRequest, error) {
	return u.purchaseRequestRepo.GetBySellerID(sellerID)
}

func (u *PurchaseRequestUsecase) GetByPetID(petID int) ([]models.PurchaseRequest, error) {
	return u.purchaseRequestRepo.GetByPetID(petID)
}

func (u *PurchaseRequestUsecase) UpdateStatus(id int, _ int, status string) (*models.PurchaseRequest, error) {
	return u.purchaseRequestRepo.UpdateStatus(id, status)
}

func (u *PurchaseRequestUsecase) Delete(id int, _ int) error {
	return u.purchaseRequestRepo.Delete(id)
}
