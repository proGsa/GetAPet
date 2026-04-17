package usecase_test

import (
	"errors"
	"testing"

	"getapet-backend/internal/models"
	"getapet-backend/internal/usecase"

	"github.com/google/uuid"
)

//
// MOCK REPOSITORY
//
type mockPetRepo struct {
	createFn          func(*models.Pet) (*models.Pet, error)
	getAllFn          func() ([]models.Pet, error)
	getByIDFn         func(uuid.UUID) (*models.Pet, error)
	getBySellerIDFn   func(uuid.UUID) ([]models.Pet, error)
	updateFn          func(uuid.UUID, *models.Pet) (*models.Pet, error)
	deleteFn          func(uuid.UUID) error
	checkBelongingFn  func(uuid.UUID, uuid.UUID) (bool, error)
}

func (m *mockPetRepo) Create(p *models.Pet) (*models.Pet, error) {
	return m.createFn(p)
}

func (m *mockPetRepo) GetAll() ([]models.Pet, error) {
	return m.getAllFn()
}

func (m *mockPetRepo) GetByID(id uuid.UUID) (*models.Pet, error) {
	return m.getByIDFn(id)
}

func (m *mockPetRepo) GetBySellerID(id uuid.UUID) ([]models.Pet, error) {
	return m.getBySellerIDFn(id)
}

func (m *mockPetRepo) Update(id uuid.UUID, p *models.Pet) (*models.Pet, error) {
	return m.updateFn(id, p)
}

func (m *mockPetRepo) Delete(id uuid.UUID) error {
	return m.deleteFn(id)
}

func (m *mockPetRepo) CheckBelonging(id, sellerID uuid.UUID) (bool, error) {
	return m.checkBelongingFn(id, sellerID)
}

//
// TESTS
//

func TestPetUsecase_Update_Success(t *testing.T) {
	repo := &mockPetRepo{
		checkBelongingFn: func(_, _ uuid.UUID) (bool, error) {
			return true, nil
		},
		updateFn: func(id uuid.UUID, p *models.Pet) (*models.Pet, error) {
			p.PetName = "Updated"
			return p, nil
		},
	}

	u := usecase.NewPetUsecase(repo)

	pet := &models.Pet{PetName: "Old"}

	updated, err := u.Update(uuid.New(), uuid.New(), pet)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if updated.PetName != "Updated" {
		t.Fatalf("expected Updated, got %s", updated.PetName)
	}
}

func TestPetUsecase_Update_NotOwner(t *testing.T) {
	repo := &mockPetRepo{
		checkBelongingFn: func(_, _ uuid.UUID) (bool, error) {
			return false, nil
		},
	}

	u := usecase.NewPetUsecase(repo)

	_, err := u.Update(uuid.New(), uuid.New(), &models.Pet{})
	if !errors.Is(err, models.ErrPetNotFound) {
		t.Fatalf("expected ErrPetNotFound, got %v", err)
	}
}

func TestPetUsecase_Delete_Success(t *testing.T) {
	repo := &mockPetRepo{
		checkBelongingFn: func(_, _ uuid.UUID) (bool, error) {
			return true, nil
		},
		deleteFn: func(uuid.UUID) error {
			return nil
		},
	}

	u := usecase.NewPetUsecase(repo)

	err := u.Delete(uuid.New(), uuid.New())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestPetUsecase_Delete_NotOwner(t *testing.T) {
	repo := &mockPetRepo{
		checkBelongingFn: func(_, _ uuid.UUID) (bool, error) {
			return false, nil
		},
	}

	u := usecase.NewPetUsecase(repo)

	err := u.Delete(uuid.New(), uuid.New())
	if !errors.Is(err, models.ErrPetNotFound) {
		t.Fatalf("expected ErrPetNotFound, got %v", err)
	}
}