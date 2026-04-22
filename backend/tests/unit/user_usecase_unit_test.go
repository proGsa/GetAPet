package usecase_test

import (
	"errors"
	"strings"
	"testing"

	"getapet-backend/internal/models"
	"getapet-backend/internal/usecase"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepo struct {
	createFn     func(*models.User) (*models.User, error)
	getAllFn     func() ([]models.User, error)
	getByIDFn    func(uuid.UUID) (*models.User, error)
	getByLoginFn func(string) (*models.User, error)
	updateFn     func(uuid.UUID, *models.User) (*models.User, error)
	deleteFn     func(uuid.UUID) error
}

func (m *mockUserRepo) Create(user *models.User) (*models.User, error) {
	return m.createFn(user)
}

func (m *mockUserRepo) GetAll() ([]models.User, error) {
	return m.getAllFn()
}

func (m *mockUserRepo) GetByID(id uuid.UUID) (*models.User, error) {
	return m.getByIDFn(id)
}

func (m *mockUserRepo) GetByLogin(login string) (*models.User, error) {
	return m.getByLoginFn(login)
}

func (m *mockUserRepo) Update(id uuid.UUID, user *models.User) (*models.User, error) {
	return m.updateFn(id, user)
}

func (m *mockUserRepo) Delete(id uuid.UUID) error {
	return m.deleteFn(id)
}

func TestUserUsecase_Create_Success(t *testing.T) {
	rawPassword := "my-plain-password"

	repo := &mockUserRepo{
		createFn: func(user *models.User) (*models.User, error) {
			if user.UserPassword == rawPassword {
				t.Fatal("password must be hashed before repository call")
			}

			if err := bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(rawPassword)); err != nil {
				t.Fatalf("password hash does not match original: %v", err)
			}

			user.ID = uuid.New()
			return user, nil
		},
	}

	u := usecase.NewUserUsecase(repo)

	user := &models.User{
		UserLogin:    "john",
		UserPassword: rawPassword,
	}

	res, err := u.Create(user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.ID == uuid.Nil {
		t.Fatal("expected ID to be set")
	}
}

func TestUserUsecase_Create_HashError(t *testing.T) {
	repo := &mockUserRepo{
		createFn: func(*models.User) (*models.User, error) {
			t.Fatal("repository Create must not be called on hash error")
			return nil, nil
		},
	}

	u := usecase.NewUserUsecase(repo)

	_, err := u.Create(&models.User{
		UserPassword: strings.Repeat("a", 73),
	})
	if !errors.Is(err, bcrypt.ErrPasswordTooLong) {
		t.Fatalf("expected ErrPasswordTooLong, got %v", err)
	}
}

func TestUserUsecase_GetAll(t *testing.T) {
	repo := &mockUserRepo{
		getAllFn: func() ([]models.User, error) {
			return []models.User{
				{ID: uuid.New()},
				{ID: uuid.New()},
			}, nil
		},
	}

	u := usecase.NewUserUsecase(repo)

	res, err := u.GetAll()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(res) != 2 {
		t.Fatalf("expected 2 users, got %d", len(res))
	}
}

func TestUserUsecase_GetByID(t *testing.T) {
	id := uuid.New()

	repo := &mockUserRepo{
		getByIDFn: func(in uuid.UUID) (*models.User, error) {
			if in != id {
				t.Fatalf("wrong id passed")
			}
			return &models.User{ID: id}, nil
		},
	}

	u := usecase.NewUserUsecase(repo)

	res, err := u.GetByID(id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.ID != id {
		t.Fatalf("expected %v, got %v", id, res.ID)
	}
}

func TestUserUsecase_GetByLogin(t *testing.T) {
	login := "john"

	repo := &mockUserRepo{
		getByLoginFn: func(in string) (*models.User, error) {
			if in != login {
				t.Fatalf("wrong login passed")
			}
			return &models.User{UserLogin: login}, nil
		},
	}

	u := usecase.NewUserUsecase(repo)

	res, err := u.GetByLogin(login)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.UserLogin != login {
		t.Fatalf("expected %s, got %s", login, res.UserLogin)
	}
}

func TestUserUsecase_Update_Success(t *testing.T) {
	id := uuid.New()
	rawPassword := "my-new-password"

	repo := &mockUserRepo{
		updateFn: func(in uuid.UUID, user *models.User) (*models.User, error) {
			if in != id {
				t.Fatalf("wrong id passed")
			}

			if user.UserPassword == rawPassword {
				t.Fatal("password must be hashed before repository call")
			}

			if err := bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(rawPassword)); err != nil {
				t.Fatalf("password hash does not match original: %v", err)
			}

			user.ID = id
			return user, nil
		},
	}

	u := usecase.NewUserUsecase(repo)

	res, err := u.Update(id, &models.User{UserPassword: rawPassword})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.ID != id {
		t.Fatalf("expected %v, got %v", id, res.ID)
	}
}

func TestUserUsecase_Update_HashError(t *testing.T) {
	repo := &mockUserRepo{
		updateFn: func(uuid.UUID, *models.User) (*models.User, error) {
			t.Fatal("repository Update must not be called on hash error")
			return nil, nil
		},
	}

	u := usecase.NewUserUsecase(repo)

	_, err := u.Update(uuid.New(), &models.User{
		UserPassword: strings.Repeat("a", 73),
	})
	if !errors.Is(err, bcrypt.ErrPasswordTooLong) {
		t.Fatalf("expected ErrPasswordTooLong, got %v", err)
	}
}

func TestUserUsecase_Delete(t *testing.T) {
	id := uuid.New()

	repo := &mockUserRepo{
		deleteFn: func(in uuid.UUID) error {
			if in != id {
				t.Fatalf("wrong id passed")
			}
			return nil
		},
	}

	u := usecase.NewUserUsecase(repo)

	err := u.Delete(id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestUserUsecase_Login_Success(t *testing.T) {
	login := "john"
	password := "valid-password"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password for test setup: %v", err)
	}

	repo := &mockUserRepo{
		getByLoginFn: func(in string) (*models.User, error) {
			if in != login {
				t.Fatalf("wrong login passed")
			}
			return &models.User{
				ID:           uuid.New(),
				UserLogin:    login,
				UserPassword: string(hashedPassword),
			}, nil
		},
	}

	u := usecase.NewUserUsecase(repo)

	res, err := u.Login(login, password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.UserLogin != login {
		t.Fatalf("expected login %s, got %s", login, res.UserLogin)
	}
}

func TestUserUsecase_Login_UserNotFound(t *testing.T) {
	repo := &mockUserRepo{
		getByLoginFn: func(string) (*models.User, error) {
			return nil, models.ErrUserNotFound
		},
	}

	u := usecase.NewUserUsecase(repo)

	_, err := u.Login("missing-user", "any-password")
	if !errors.Is(err, models.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestUserUsecase_Login_InvalidPassword(t *testing.T) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password for test setup: %v", err)
	}

	repo := &mockUserRepo{
		getByLoginFn: func(string) (*models.User, error) {
			return &models.User{UserPassword: string(hashedPassword)}, nil
		},
	}

	u := usecase.NewUserUsecase(repo)

	_, err = u.Login("john", "wrong-password")
	if !errors.Is(err, models.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestUserUsecase_Login_RepoError(t *testing.T) {
	expectedErr := errors.New("database unavailable")

	repo := &mockUserRepo{
		getByLoginFn: func(string) (*models.User, error) {
			return nil, expectedErr
		},
	}

	u := usecase.NewUserUsecase(repo)

	_, err := u.Login("john", "password")
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}
