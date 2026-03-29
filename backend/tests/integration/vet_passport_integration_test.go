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

func TestVetPassportRepositoryIntegration(t *testing.T) {
	ctx := context.Background()
	db, cleanup := setupTestPostgresVet(t, ctx)
	defer cleanup()

	repo := repository.NewVetPassportRepository(db)

	t.Run("Create and GetByID", func(t *testing.T) {
		clearVetTable(t, db)

		passport := &models.VetPassport{
			Chipping:           true,
			Sterilization:      false,
			HealthIssues:       "none",
			Vaccinations:       "rabies",
			ParasiteTreatments: "yes",
		}

		created, err := repo.Create(passport)
		if err != nil {
			t.Fatalf("create: %v", err)
		}
		if created.ID == uuid.Nil {
			t.Fatal("expected non-nil id")
		}

		got, err := repo.GetByID(created.ID)
		if err != nil {
			t.Fatalf("get by id: %v", err)
		}
		if got.HealthIssues != "none" {
			t.Fatalf("expected none, got %s", got.HealthIssues)
		}
	})

	t.Run("GetAll", func(t *testing.T) {
		clearVetTable(t, db)

		for _, p := range []models.VetPassport{
			{
				Chipping:      true,
				HealthIssues:  "a",
				Vaccinations:  "v1",
			},
			{
				Chipping:      false,
				HealthIssues:  "b",
				Vaccinations:  "v2",
			},
		} {
			if _, err := repo.Create(&p); err != nil {
				t.Fatalf("create: %v", err)
			}
		}

		all, err := repo.GetAll()
		if err != nil {
			t.Fatalf("get all: %v", err)
		}
		if len(all) != 2 {
			t.Fatalf("expected 2, got %d", len(all))
		}
	})

	t.Run("Update", func(t *testing.T) {
		clearVetTable(t, db)

		created, _ := repo.Create(&models.VetPassport{
			Chipping:     true,
			HealthIssues: "old",
		})

		updated, err := repo.Update(created.ID, &models.VetPassport{
			Chipping:     false,
			HealthIssues: "new",
		})
		if err != nil {
			t.Fatalf("update: %v", err)
		}

		if updated.HealthIssues != "new" {
			t.Fatalf("expected new, got %s", updated.HealthIssues)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		clearVetTable(t, db)

		created, _ := repo.Create(&models.VetPassport{
			Chipping: true,
		})

		err := repo.Delete(created.ID)
		if err != nil {
			t.Fatalf("delete: %v", err)
		}

		_, err = repo.GetByID(created.ID)
		if !errors.Is(err, models.ErrVetPassportNotFound) {
			t.Fatalf("expected ErrVetPassportNotFound, got %v", err)
		}
	})

	t.Run("GetByID not found", func(t *testing.T) {
		clearVetTable(t, db)

		_, err := repo.GetByID(uuid.New())
		if !errors.Is(err, models.ErrVetPassportNotFound) {
			t.Fatalf("expected ErrVetPassportNotFound, got %v", err)
		}
	})

	t.Run("Update not found", func(t *testing.T) {
		clearVetTable(t, db)

		_, err := repo.Update(uuid.New(), &models.VetPassport{
			HealthIssues: "ghost",
		})
		if !errors.Is(err, models.ErrVetPassportNotFound) {
			t.Fatalf("expected ErrVetPassportNotFound, got %v", err)
		}
	})

	t.Run("Delete not found", func(t *testing.T) {
		clearVetTable(t, db)

		err := repo.Delete(uuid.New())
		if !errors.Is(err, models.ErrVetPassportNotFound) {
			t.Fatalf("expected ErrVetPassportNotFound, got %v", err)
		}
	})
}

func setupTestPostgresVet(t *testing.T, ctx context.Context) (*sql.DB, func()) {
	t.Helper()

	if dsn := os.Getenv("INTEGRATION_DB_DSN"); dsn != "" {
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			t.Fatalf("open db: %v", err)
		}

		for i := 0; i < 40; i++ {
			if db.PingContext(ctx) == nil {
				break
			}
			time.Sleep(500 * time.Millisecond)
		}

		initVetSchema(ctx, db)
		return db, func() { _ = db.Close() }
	}

	testcontainers.SkipIfProviderIsNotHealthy(t)

	pg, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("getapet_test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
	)
	if err != nil {
		t.Fatalf("container: %v", err)
	}

	conn, _ := pg.ConnectionString(ctx, "sslmode=disable")
	db, _ := sql.Open("postgres", conn)

	for i := 0; i < 20; i++ {
		if db.PingContext(ctx) == nil {
			break
		}
		time.Sleep(300 * time.Millisecond)
	}

	initVetSchema(ctx, db)

	return db, func() {
		_ = db.Close()
		_ = pg.Terminate(context.Background())
	}
}

func initVetSchema(ctx context.Context, db *sql.DB) {
	db.ExecContext(ctx, `
	CREATE EXTENSION IF NOT EXISTS pgcrypto;

	CREATE TABLE IF NOT EXISTS vet_passport (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		chipping BOOLEAN,
		sterilization BOOLEAN,
		health_issues TEXT,
		vaccinations TEXT,
		parasite_treatments TEXT
	);
	`)
}

func clearVetTable(t *testing.T, db *sql.DB) {
	t.Helper()
	if _, err := db.Exec(`TRUNCATE TABLE vet_passport`); err != nil {
		t.Fatalf("truncate vet_passport: %v", err)
	}
}