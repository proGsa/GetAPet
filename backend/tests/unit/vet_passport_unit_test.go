package usecase_test

import (
	"testing"

	"getapet-backend/internal/models"
	"getapet-backend/internal/usecase"

	"github.com/google/uuid"
)

//
// MOCK REPOSITORY
//

type mockVetPassportRepo struct {
	createFn  func(*models.VetPassport) (*models.VetPassport, error)
	getAllFn  func() ([]models.VetPassport, error)
	getByIDFn func(uuid.UUID) (*models.VetPassport, error)
	updateFn  func(uuid.UUID, *models.VetPassport) (*models.VetPassport, error)
	deleteFn  func(uuid.UUID) error
}

func (m *mockVetPassportRepo) Create(p *models.VetPassport) (*models.VetPassport, error) {
	return m.createFn(p)
}

func (m *mockVetPassportRepo) GetAll() ([]models.VetPassport, error) {
	return m.getAllFn()
}

func (m *mockVetPassportRepo) GetByID(id uuid.UUID) (*models.VetPassport, error) {
	return m.getByIDFn(id)
}

func (m *mockVetPassportRepo) Update(id uuid.UUID, p *models.VetPassport) (*models.VetPassport, error) {
	return m.updateFn(id, p)
}

func (m *mockVetPassportRepo) Delete(id uuid.UUID) error {
	return m.deleteFn(id)
}

//
// TESTS
//

func TestVetPassportUsecase_Create(t *testing.T) {
	repo := &mockVetPassportRepo{
		createFn: func(p *models.VetPassport) (*models.VetPassport, error) {
			p.ID = uuid.New()
			return p, nil
		},
	}

	u := usecase.NewVetPassportUsecase(repo)

	passport := &models.VetPassport{}

	res, err := u.Create(passport)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.ID == uuid.Nil {
		t.Fatal("expected ID to be set")
	}
}

func TestVetPassportUsecase_GetAll(t *testing.T) {
	repo := &mockVetPassportRepo{
		getAllFn: func() ([]models.VetPassport, error) {
			return []models.VetPassport{
				{ID: uuid.New()},
				{ID: uuid.New()},
			}, nil
		},
	}

	u := usecase.NewVetPassportUsecase(repo)

	res, err := u.GetAll()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(res) != 2 {
		t.Fatalf("expected 2 passports, got %d", len(res))
	}
}

func TestVetPassportUsecase_GetByID(t *testing.T) {
	id := uuid.New()

	repo := &mockVetPassportRepo{
		getByIDFn: func(in uuid.UUID) (*models.VetPassport, error) {
			if in != id {
				t.Fatalf("wrong id passed")
			}
			return &models.VetPassport{ID: id}, nil
		},
	}

	u := usecase.NewVetPassportUsecase(repo)

	res, err := u.GetByID(id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.ID != id {
		t.Fatalf("expected %v, got %v", id, res.ID)
	}
}

func TestVetPassportUsecase_Update(t *testing.T) {
	id := uuid.New()

	repo := &mockVetPassportRepo{
		updateFn: func(in uuid.UUID, p *models.VetPassport) (*models.VetPassport, error) {
			if in != id {
				t.Fatalf("wrong id passed")
			}
			p.ID = id
			return p, nil
		},
	}

	u := usecase.NewVetPassportUsecase(repo)

	res, err := u.Update(id, &models.VetPassport{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.ID != id {
		t.Fatalf("expected %v, got %v", id, res.ID)
	}
}

func TestVetPassportUsecase_Delete(t *testing.T) {
	id := uuid.New()

	repo := &mockVetPassportRepo{
		deleteFn: func(in uuid.UUID) error {
			if in != id {
				t.Fatalf("wrong id passed")
			}
			return nil
		},
	}

	u := usecase.NewVetPassportUsecase(repo)

	err := u.Delete(id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
