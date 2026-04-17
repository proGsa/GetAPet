//go:build integration

package integration_test

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"
	"time"

	"getapet-backend/internal/models"
	"getapet-backend/internal/repository"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	testcontainers "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestUserRepositoryIntegration(t *testing.T) {
	ctx := context.Background()
	db, cleanup := setupTestPostgres(t, ctx)
	defer cleanup()

	repo := repository.NewUserRepository(db)

	t.Run("Create and GetByID", func(t *testing.T) {
		clearUsersTable(t, db)

		created, err := repo.Create(&models.User{
			FIO:             "Ivan Ivanov",
			TelephoneNumber: "+79991234567",
			City:            "Moscow",
			UserLogin:       "ivan_login",
			UserPassword:    "hash1",
			Status:          "active",
			UserDescription: "desc",
		})
		if err != nil {
			t.Fatalf("create user: %v", err)
		}
		if created.ID == uuid.Nil {
			t.Fatal("expected non-nil user id")
		}

		got, err := repo.GetByID(created.ID)
		if err != nil {
			t.Fatalf("get by id: %v", err)
		}
		if got.UserLogin != "ivan_login" {
			t.Fatalf("expected login ivan_login, got %s", got.UserLogin)
		}
	})

	t.Run("GetByLogin", func(t *testing.T) {
		clearUsersTable(t, db)

		_, err := repo.Create(&models.User{
			FIO:             "Maria Petrova",
			TelephoneNumber: "+79990000000",
			City:            "Kazan",
			UserLogin:       "maria_login",
			UserPassword:    "hash2",
			Status:          "active",
			UserDescription: "desc2",
		})
		if err != nil {
			t.Fatalf("create user: %v", err)
		}

		got, err := repo.GetByLogin("maria_login")
		if err != nil {
			t.Fatalf("get by login: %v", err)
		}
		if got.FIO != "Maria Petrova" {
			t.Fatalf("expected fio Maria Petrova, got %s", got.FIO)
		}
	})

	t.Run("GetAll", func(t *testing.T) {
		clearUsersTable(t, db)

		for _, u := range []models.User{
			{
				FIO:             "A User",
				TelephoneNumber: "+70000000001",
				City:            "City1",
				UserLogin:       "u1",
				UserPassword:    "h1",
				Status:          "active",
			},
			{
				FIO:             "B User",
				TelephoneNumber: "+70000000002",
				City:            "City2",
				UserLogin:       "u2",
				UserPassword:    "h2",
				Status:          "active",
			},
		} {
			if _, err := repo.Create(&u); err != nil {
				t.Fatalf("create user: %v", err)
			}
		}

		users, err := repo.GetAll()
		if err != nil {
			t.Fatalf("get all users: %v", err)
		}
		if len(users) != 2 {
			t.Fatalf("expected 2 users, got %d", len(users))
		}
	})

	t.Run("Update", func(t *testing.T) {
		clearUsersTable(t, db)

		created, err := repo.Create(&models.User{
			FIO:             "Old Name",
			TelephoneNumber: "+71111111111",
			City:            "Old City",
			UserLogin:       "old_login",
			UserPassword:    "old_hash",
			Status:          "active",
			UserDescription: "old_desc",
		})
		if err != nil {
			t.Fatalf("create user: %v", err)
		}

		updated, err := repo.Update(created.ID, &models.User{
			FIO:             "New Name",
			TelephoneNumber: "+72222222222",
			City:            "New City",
			UserLogin:       "new_login",
			UserPassword:    "new_hash",
			Status:          "blocked",
			UserDescription: "new_desc",
		})
		if err != nil {
			t.Fatalf("update user: %v", err)
		}
		if updated.FIO != "New Name" || updated.UserLogin != "new_login" {
			t.Fatalf("update result mismatch: %+v", updated)
		}
	})

	t.Run("Delete and not found behavior", func(t *testing.T) {
		clearUsersTable(t, db)

		created, err := repo.Create(&models.User{
			FIO:             "Delete Me",
			TelephoneNumber: "+73333333333",
			City:            "Perm",
			UserLogin:       "to_delete",
			UserPassword:    "hash",
			Status:          "active",
			UserDescription: "tmp",
		})
		if err != nil {
			t.Fatalf("create user: %v", err)
		}

		if err := repo.Delete(created.ID); err != nil {
			t.Fatalf("delete user: %v", err)
		}

		_, err = repo.GetByID(created.ID)
		if !errors.Is(err, models.ErrUserNotFound) {
			t.Fatalf("expected ErrUserNotFound, got %v", err)
		}

		err = repo.Delete(created.ID)
		if !errors.Is(err, models.ErrUserNotFound) {
			t.Fatalf("expected ErrUserNotFound on second delete, got %v", err)
		}
	})

	t.Run("GetByLogin not found", func(t *testing.T) {
		clearUsersTable(t, db)

		_, err := repo.GetByLogin("unknown_login")
		if !errors.Is(err, models.ErrUserNotFound) {
			t.Fatalf("expected ErrUserNotFound, got %v", err)
		}
	})

	t.Run("Update not found", func(t *testing.T) {
		clearUsersTable(t, db)

		_, err := repo.Update(uuid.New(), &models.User{
			FIO:             "Ghost User",
			TelephoneNumber: "+70000000000",
			City:            "Nowhere",
			UserLogin:       "ghost",
			UserPassword:    "hash",
			Status:          "active",
			UserDescription: "ghost",
		})
		if !errors.Is(err, models.ErrUserNotFound) {
			t.Fatalf("expected ErrUserNotFound, got %v", err)
		}
	})

	t.Run("Create duplicate login", func(t *testing.T) {
		clearUsersTable(t, db)

		_, err := repo.Create(&models.User{
			FIO:             "First User",
			TelephoneNumber: "+74444444444",
			City:            "Moscow",
			UserLogin:       "duplicate_login",
			UserPassword:    "hash",
			Status:          "active",
			UserDescription: "first",
		})
		if err != nil {
			t.Fatalf("create first user: %v", err)
		}

		_, err = repo.Create(&models.User{
			FIO:             "Second User",
			TelephoneNumber: "+75555555555",
			City:            "Moscow",
			UserLogin:       "duplicate_login",
			UserPassword:    "hash",
			Status:          "active",
			UserDescription: "second",
		})
		if err == nil {
			t.Fatal("expected unique constraint error for duplicate user_login")
		}
	})

	t.Run("Create with telephone_number too long", func(t *testing.T) {
		clearUsersTable(t, db)

		_, err := repo.Create(&models.User{
			FIO:             "Bad Phone User",
			TelephoneNumber: "+700000000000000000000000000000",
			City:            "Moscow",
			UserLogin:       "bad_phone",
			UserPassword:    "hash",
			Status:          "active",
		})
		if err == nil {
			t.Fatal("expected length constraint error for telephone_number")
		}
	})

	t.Run("Create with city too long", func(t *testing.T) {
		clearUsersTable(t, db)

		_, err := repo.Create(&models.User{
			FIO:             "Long City User",
			TelephoneNumber: "+76666666666",
			City:            "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz",
			UserLogin:       "long_city",
			UserPassword:    "hash",
			Status:          "active",
		})
		if err == nil {
			t.Fatal("expected length constraint error for city")
		}
	})
}

func setupTestPostgres(t *testing.T, ctx context.Context) (*sql.DB, func()) {
	t.Helper()

	if dsn := os.Getenv("INTEGRATION_DB_DSN"); dsn != "" {
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			t.Fatalf("open external test db: %v", err)
		}

		var pingErr error
		for i := 0; i < 40; i++ {
			pingErr = db.PingContext(ctx)
			if pingErr == nil {
				break
			}
			time.Sleep(500 * time.Millisecond)
		}
		if pingErr != nil {
			_ = db.Close()
			t.Fatalf("ping external test db: %v", pingErr)
		}

		if err := initUserSchema(ctx, db); err != nil {
			_ = db.Close()
			t.Fatalf("init schema: %v", err)
		}

		return db, func() { _ = db.Close() }
	}

	testcontainers.SkipIfProviderIsNotHealthy(t)

	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("getapet_test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
	)
	if err != nil {
		t.Fatalf("start postgres container: %v", err)
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = pgContainer.Terminate(ctx)
		t.Fatalf("build connection string: %v", err)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		_ = pgContainer.Terminate(ctx)
		t.Fatalf("open db: %v", err)
	}

	var pingErr error
	for i := 0; i < 20; i++ {
		pingErr = db.PingContext(ctx)
		if pingErr == nil {
			break
		}
		time.Sleep(300 * time.Millisecond)
	}
	if pingErr != nil {
		_ = db.Close()
		_ = pgContainer.Terminate(ctx)
		t.Fatalf("ping db: %v", pingErr)
	}

	if err := initUserSchema(ctx, db); err != nil {
		_ = db.Close()
		_ = pgContainer.Terminate(ctx)
		t.Fatalf("init schema: %v", err)
	}

	cleanup := func() {
		_ = db.Close()
		_ = pgContainer.Terminate(context.Background())
	}

	return db, cleanup
}

func initUserSchema(ctx context.Context, db *sql.DB) error {
	const createSchema = `
		CREATE EXTENSION IF NOT EXISTS pgcrypto;
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			fio VARCHAR(255) NOT NULL,
			telephone_number VARCHAR(20) NOT NULL,
			city VARCHAR(50),
			user_login VARCHAR(50) UNIQUE,
			user_password VARCHAR(255),
			status VARCHAR(20),
			user_description TEXT
		);
	`
	_, err := db.ExecContext(ctx, createSchema)
	return err
}

func clearUsersTable(t *testing.T, db *sql.DB) {
	t.Helper()
	if _, err := db.Exec(`TRUNCATE TABLE users CASCADE`); err != nil {
		t.Fatalf("truncate users: %v", err)
	}
}
