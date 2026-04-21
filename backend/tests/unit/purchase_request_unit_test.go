package usecase_test

import (
	"testing"

	"getapet-backend/internal/models"
	"getapet-backend/internal/usecase"

	"github.com/google/uuid"
)

type mockPurchaseRequestRepo struct {
	createFn               func(*models.PurchaseRequest) (*models.PurchaseRequest, error)
	getAllFn               func() ([]models.PurchaseRequest, error)
	getByIDFn              func(uuid.UUID) (*models.PurchaseRequest, error)
	getByBuyerIDFn         func(uuid.UUID) ([]models.PurchaseRequest, error)
	getBySellerIDFn        func(uuid.UUID) ([]models.PurchaseRequest, error)
	getByPetIDFn           func(uuid.UUID) ([]models.PurchaseRequest, error)
	updateStatusFn         func(uuid.UUID, string) (*models.PurchaseRequest, error)
	updateStatusBySellerFn func(uuid.UUID, uuid.UUID, string) (*models.PurchaseRequest, error)
	deleteFn               func(uuid.UUID) error
	deleteByBuyerFn        func(uuid.UUID, uuid.UUID) error
}

func (m *mockPurchaseRequestRepo) Create(request *models.PurchaseRequest) (*models.PurchaseRequest, error) {
	return m.createFn(request)
}

func (m *mockPurchaseRequestRepo) GetAll() ([]models.PurchaseRequest, error) {
	return m.getAllFn()
}

func (m *mockPurchaseRequestRepo) GetByID(id uuid.UUID) (*models.PurchaseRequest, error) {
	return m.getByIDFn(id)
}

func (m *mockPurchaseRequestRepo) GetByBuyerID(buyerID uuid.UUID) ([]models.PurchaseRequest, error) {
	return m.getByBuyerIDFn(buyerID)
}

func (m *mockPurchaseRequestRepo) GetBySellerID(sellerID uuid.UUID) ([]models.PurchaseRequest, error) {
	return m.getBySellerIDFn(sellerID)
}

func (m *mockPurchaseRequestRepo) GetByPetID(petID uuid.UUID) ([]models.PurchaseRequest, error) {
	return m.getByPetIDFn(petID)
}

func (m *mockPurchaseRequestRepo) UpdateStatus(id uuid.UUID, status string) (*models.PurchaseRequest, error) {
	return m.updateStatusFn(id, status)
}

func (m *mockPurchaseRequestRepo) UpdateStatusBySeller(id uuid.UUID, sellerID uuid.UUID, status string) (*models.PurchaseRequest, error) {
	return m.updateStatusBySellerFn(id, sellerID, status)
}

func (m *mockPurchaseRequestRepo) Delete(id uuid.UUID) error {
	return m.deleteFn(id)
}

func (m *mockPurchaseRequestRepo) DeleteByBuyer(id uuid.UUID, buyerID uuid.UUID) error {
	return m.deleteByBuyerFn(id, buyerID)
}

func TestPurchaseRequestUsecase_Create(t *testing.T) {
	repo := &mockPurchaseRequestRepo{
		createFn: func(request *models.PurchaseRequest) (*models.PurchaseRequest, error) {
			request.ID = uuid.New()
			return request, nil
		},
	}

	u := usecase.NewPurchaseRequestUsecase(repo)

	res, err := u.Create(&models.PurchaseRequest{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.ID == uuid.Nil {
		t.Fatal("expected ID to be set")
	}
}

func TestPurchaseRequestUsecase_GetAll(t *testing.T) {
	repo := &mockPurchaseRequestRepo{
		getAllFn: func() ([]models.PurchaseRequest, error) {
			return []models.PurchaseRequest{
				{ID: uuid.New()},
				{ID: uuid.New()},
			}, nil
		},
	}

	u := usecase.NewPurchaseRequestUsecase(repo)

	res, err := u.GetAll()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(res) != 2 {
		t.Fatalf("expected 2 requests, got %d", len(res))
	}
}

func TestPurchaseRequestUsecase_GetByID(t *testing.T) {
	id := uuid.New()

	repo := &mockPurchaseRequestRepo{
		getByIDFn: func(in uuid.UUID) (*models.PurchaseRequest, error) {
			if in != id {
				t.Fatal("wrong id passed")
			}
			return &models.PurchaseRequest{ID: id}, nil
		},
	}

	u := usecase.NewPurchaseRequestUsecase(repo)

	res, err := u.GetByID(id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.ID != id {
		t.Fatalf("expected %v, got %v", id, res.ID)
	}
}

func TestPurchaseRequestUsecase_GetByBuyerID(t *testing.T) {
	buyerID := uuid.New()

	repo := &mockPurchaseRequestRepo{
		getByBuyerIDFn: func(in uuid.UUID) ([]models.PurchaseRequest, error) {
			if in != buyerID {
				t.Fatal("wrong buyerID passed")
			}
			return []models.PurchaseRequest{
				{ID: uuid.New(), BuyerID: buyerID},
			}, nil
		},
	}

	u := usecase.NewPurchaseRequestUsecase(repo)

	res, err := u.GetByBuyerID(buyerID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(res) != 1 {
		t.Fatalf("expected 1 request, got %d", len(res))
	}
}

func TestPurchaseRequestUsecase_GetBySellerID(t *testing.T) {
	sellerID := uuid.New()

	repo := &mockPurchaseRequestRepo{
		getBySellerIDFn: func(in uuid.UUID) ([]models.PurchaseRequest, error) {
			if in != sellerID {
				t.Fatal("wrong sellerID passed")
			}
			return []models.PurchaseRequest{
				{ID: uuid.New(), SellerID: sellerID},
			}, nil
		},
	}

	u := usecase.NewPurchaseRequestUsecase(repo)

	res, err := u.GetBySellerID(sellerID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(res) != 1 {
		t.Fatalf("expected 1 request, got %d", len(res))
	}
}

func TestPurchaseRequestUsecase_GetByPetID(t *testing.T) {
	petID := uuid.New()

	repo := &mockPurchaseRequestRepo{
		getByPetIDFn: func(in uuid.UUID) ([]models.PurchaseRequest, error) {
			if in != petID {
				t.Fatal("wrong petID passed")
			}
			return []models.PurchaseRequest{
				{ID: uuid.New(), PetID: petID},
			}, nil
		},
	}

	u := usecase.NewPurchaseRequestUsecase(repo)

	res, err := u.GetByPetID(petID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(res) != 1 {
		t.Fatalf("expected 1 request, got %d", len(res))
	}
}

func TestPurchaseRequestUsecase_UpdateStatus_WithoutSeller(t *testing.T) {
	id := uuid.New()
	status := "approved"

	repo := &mockPurchaseRequestRepo{
		updateStatusFn: func(inID uuid.UUID, inStatus string) (*models.PurchaseRequest, error) {
			if inID != id {
				t.Fatal("wrong id passed")
			}
			if inStatus != status {
				t.Fatal("wrong status passed")
			}
			return &models.PurchaseRequest{ID: id, Status: status}, nil
		},
		updateStatusBySellerFn: func(uuid.UUID, uuid.UUID, string) (*models.PurchaseRequest, error) {
			t.Fatal("UpdateStatusBySeller must not be called when sellerID = uuid.Nil")
			return nil, nil
		},
	}

	u := usecase.NewPurchaseRequestUsecase(repo)

	res, err := u.UpdateStatus(id, uuid.Nil, status)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.Status != status {
		t.Fatalf("expected status %s, got %s", status, res.Status)
	}
}

func TestPurchaseRequestUsecase_UpdateStatus_WithSeller(t *testing.T) {
	id := uuid.New()
	sellerID := uuid.New()
	status := "declined"

	repo := &mockPurchaseRequestRepo{
		updateStatusFn: func(uuid.UUID, string) (*models.PurchaseRequest, error) {
			t.Fatal("UpdateStatus must not be called when sellerID != uuid.Nil")
			return nil, nil
		},
		updateStatusBySellerFn: func(inID uuid.UUID, inSellerID uuid.UUID, inStatus string) (*models.PurchaseRequest, error) {
			if inID != id {
				t.Fatal("wrong id passed")
			}
			if inSellerID != sellerID {
				t.Fatal("wrong sellerID passed")
			}
			if inStatus != status {
				t.Fatal("wrong status passed")
			}
			return &models.PurchaseRequest{ID: id, SellerID: sellerID, Status: status}, nil
		},
	}

	u := usecase.NewPurchaseRequestUsecase(repo)

	res, err := u.UpdateStatus(id, sellerID, status)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.Status != status {
		t.Fatalf("expected status %s, got %s", status, res.Status)
	}
}

func TestPurchaseRequestUsecase_UpdateStatusBySeller(t *testing.T) {
	id := uuid.New()
	sellerID := uuid.New()
	status := "approved"

	repo := &mockPurchaseRequestRepo{
		updateStatusBySellerFn: func(inID uuid.UUID, inSellerID uuid.UUID, inStatus string) (*models.PurchaseRequest, error) {
			if inID != id {
				t.Fatal("wrong id passed")
			}
			if inSellerID != sellerID {
				t.Fatal("wrong sellerID passed")
			}
			if inStatus != status {
				t.Fatal("wrong status passed")
			}
			return &models.PurchaseRequest{ID: id, SellerID: sellerID, Status: status}, nil
		},
	}

	u := usecase.NewPurchaseRequestUsecase(repo)

	res, err := u.UpdateStatusBySeller(id, sellerID, status)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.Status != status {
		t.Fatalf("expected status %s, got %s", status, res.Status)
	}
}

func TestPurchaseRequestUsecase_Delete_WithoutBuyer(t *testing.T) {
	id := uuid.New()

	repo := &mockPurchaseRequestRepo{
		deleteFn: func(inID uuid.UUID) error {
			if inID != id {
				t.Fatal("wrong id passed")
			}
			return nil
		},
		deleteByBuyerFn: func(uuid.UUID, uuid.UUID) error {
			t.Fatal("DeleteByBuyer must not be called when buyerID = uuid.Nil")
			return nil
		},
	}

	u := usecase.NewPurchaseRequestUsecase(repo)

	err := u.Delete(id, uuid.Nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestPurchaseRequestUsecase_Delete_WithBuyer(t *testing.T) {
	id := uuid.New()
	buyerID := uuid.New()

	repo := &mockPurchaseRequestRepo{
		deleteFn: func(uuid.UUID) error {
			t.Fatal("Delete must not be called when buyerID != uuid.Nil")
			return nil
		},
		deleteByBuyerFn: func(inID uuid.UUID, inBuyerID uuid.UUID) error {
			if inID != id {
				t.Fatal("wrong id passed")
			}
			if inBuyerID != buyerID {
				t.Fatal("wrong buyerID passed")
			}
			return nil
		},
	}

	u := usecase.NewPurchaseRequestUsecase(repo)

	err := u.Delete(id, buyerID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestPurchaseRequestUsecase_DeleteByBuyer(t *testing.T) {
	id := uuid.New()
	buyerID := uuid.New()

	repo := &mockPurchaseRequestRepo{
		deleteByBuyerFn: func(inID uuid.UUID, inBuyerID uuid.UUID) error {
			if inID != id {
				t.Fatal("wrong id passed")
			}
			if inBuyerID != buyerID {
				t.Fatal("wrong buyerID passed")
			}
			return nil
		},
	}

	u := usecase.NewPurchaseRequestUsecase(repo)

	err := u.DeleteByBuyer(id, buyerID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
