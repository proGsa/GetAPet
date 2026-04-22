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

func TestPetRepositoryIntegration(t *testing.T) {
	ctx := context.Background()
	db, cleanup := setupTestPostgresPet(t, ctx)
	defer cleanup()

	repo := repository.NewPetRepository(db)

	t.Run("Create and GetByID", func(t *testing.T) {
		clearPetTable(t, db)

		pet := &models.Pet{
			VetPassportID:    uuid.New(),
			SellerID:         uuid.New(),
			PetName:          "Barsik",
			Species:          "cat",
			PetAge:           3,
			Color:            "black",
			PetGender:        "male",
			Breed:            "british",
			Pedigree:         true,
			GoodWithChildren: true,
			GoodWithAnimals:  true,
			PetDescription:   "cute cat",
			IsActive:         true,
			Price:            100.0,
		}

		created, err := repo.Create(pet)
		if err != nil {
			t.Fatalf("create pet: %v", err)
		}
		if created.ID == uuid.Nil {
			t.Fatal("expected non-nil id")
		}

		got, err := repo.GetByID(created.ID)
		if err != nil {
			t.Fatalf("get by id: %v", err)
		}
		if got.PetName != "Barsik" {
			t.Fatalf("expected Barsik, got %s", got.PetName)
		}
	})

	t.Run("GetAll", func(t *testing.T) {
		clearPetTable(t, db)

		for _, p := range []models.Pet{
			{
				VetPassportID: uuid.New(),
				SellerID:      uuid.New(),
				PetName:       "A",
				Species:       "dog",
				PetAge:        1,
				IsActive:      true,
				Price:         10,
			},
			{
				VetPassportID: uuid.New(),
				SellerID:      uuid.New(),
				PetName:       "B",
				Species:       "cat",
				PetAge:        2,
				IsActive:      true,
				Price:         20,
			},
		} {
			if _, err := repo.Create(&p); err != nil {
				t.Fatalf("create pet: %v", err)
			}
		}

		pets, err := repo.GetAll()
		if err != nil {
			t.Fatalf("get all: %v", err)
		}
		if len(pets) != 2 {
			t.Fatalf("expected 2 pets, got %d", len(pets))
		}
	})

	t.Run("GetBySellerID", func(t *testing.T) {
		clearPetTable(t, db)

		sellerID := uuid.New()

		_, _ = repo.Create(&models.Pet{
			VetPassportID: uuid.New(),
			SellerID:      sellerID,
			PetName:       "Owned",
			Species:       "dog",
			PetAge:        2,
			IsActive:      true,
			Price:         10,
		})

		_, _ = repo.Create(&models.Pet{
			VetPassportID: uuid.New(),
			SellerID:      uuid.New(),
			PetName:       "Other",
			Species:       "cat",
			PetAge:        3,
			IsActive:      true,
			Price:         20,
		})

		pets, err := repo.GetBySellerID(sellerID)
		if err != nil {
			t.Fatalf("get by seller: %v", err)
		}
		if len(pets) != 1 {
			t.Fatalf("expected 1 pet, got %d", len(pets))
		}
	})

	t.Run("Update", func(t *testing.T) {
		clearPetTable(t, db)

		created, _ := repo.Create(&models.Pet{
			VetPassportID: uuid.New(),
			SellerID:      uuid.New(),
			PetName:       "Old",
			Species:       "dog",
			PetAge:        1,
			IsActive:      true,
			Price:         10,
		})

		updated, err := repo.Update(created.ID, &models.Pet{
			PetName:  "New",
			Species:  "cat",
			PetAge:   5,
			IsActive: false,
			Price:    999,
		})
		if err != nil {
			t.Fatalf("update: %v", err)
		}

		if updated.PetName != "New" {
			t.Fatalf("expected New, got %s", updated.PetName)
		}
	})

	t.Run("GetByID not found", func(t *testing.T) {
		clearPetTable(t, db)

		_, err := repo.GetByID(uuid.New())
		if !errors.Is(err, repository.ErrPetNotFound) {
			t.Fatalf("expected ErrPetNotFound, got %v", err)
		}
	})

	t.Run("CheckBelonging", func(t *testing.T) {
		clearPetTable(t, db)

		sellerID := uuid.New()

		pet, _ := repo.Create(&models.Pet{
			VetPassportID: uuid.New(),
			SellerID:      sellerID,
			PetName:       "Test",
			Species:       "dog",
			PetAge:        1,
			IsActive:      true,
			Price:         10,
		})

		ok, err := repo.CheckBelonging(pet.ID, sellerID)
		if err != nil {
			t.Fatalf("check belonging: %v", err)
		}
		if !ok {
			t.Fatal("expected true")
		}
	})
	t.Run("Delete", func(t *testing.T) {
		clearPetTable(t, db)

		pet, err := repo.Create(&models.Pet{
			VetPassportID: uuid.New(),
			SellerID:      uuid.New(),
			PetName:       "ToDelete",
			Species:       "dog",
			PetAge:        2,
			IsActive:      true,
			Price:         50,
		})
		if err != nil {
			t.Fatalf("create pet: %v", err)
		}

		// delete
		err = repo.Delete(pet.ID)
		if err != nil {
			t.Fatalf("delete: %v", err)
		}

		// check that it's gone
		_, err = repo.GetByID(pet.ID)
		if !errors.Is(err, repository.ErrPetNotFound) {
			t.Fatalf("expected ErrPetNotFound, got %v", err)
		}
	})
}

func setupTestPostgresPet(t *testing.T, ctx context.Context) (*sql.DB, func()) {
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

		initPetSchema(ctx, db)
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

	initPetSchema(ctx, db)

	return db, func() {
		_ = db.Close()
		_ = pg.Terminate(context.Background())
	}
}

func initPetSchema(ctx context.Context, db *sql.DB) {
	db.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS pet (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		vet_passport_id UUID,
		seller_id UUID,
		pet_name TEXT,
		species TEXT,
		pet_age INT,
		color TEXT,
		pet_gender TEXT,
		breed TEXT,
		pedigree BOOL,
		good_with_children BOOL,
		good_with_animals BOOL,
		pet_description TEXT,
		is_active BOOL,
		price NUMERIC
	);
	`)
}

func clearPetTable(t *testing.T, db *sql.DB) {
	t.Helper()
	if _, err := db.Exec(`TRUNCATE TABLE pet CASCADE`); err != nil {
		t.Fatalf("truncate pet: %v", err)
	}
}
