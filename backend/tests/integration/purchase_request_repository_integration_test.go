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

func TestPurchaseRequestRepositoryIntegration(t *testing.T) {
	ctx := context.Background()
	db, cleanup := setupTestPostgresPurchaseRequest(t, ctx)
	defer cleanup()

	repo := repository.NewPurchaseRequestRepository(db)

	t.Run("Create and GetByID", func(t *testing.T) {
		clearPurchaseTables(t, db)

		sellerID := createUserForPurchase(t, db, "seller")
		buyerID := createUserForPurchase(t, db, "buyer")
		petID := createPetForPurchase(t, db, sellerID, true)

		created, err := repo.Create(&models.PurchaseRequest{
			PetID:   petID,
			BuyerID: buyerID,
		})
		if err != nil {
			t.Fatalf("create request: %v", err)
		}
		if created.ID == uuid.Nil {
			t.Fatal("expected non-nil request id")
		}
		if created.Status != "pending" {
			t.Fatalf("expected pending status, got %s", created.Status)
		}

		got, err := repo.GetByID(created.ID)
		if err != nil {
			t.Fatalf("get by id: %v", err)
		}
		if got.BuyerID != buyerID {
			t.Fatalf("expected buyer %s, got %s", buyerID, got.BuyerID)
		}
	})

	t.Run("Create fails for inactive pet", func(t *testing.T) {
		clearPurchaseTables(t, db)

		sellerID := createUserForPurchase(t, db, "seller")
		buyerID := createUserForPurchase(t, db, "buyer")
		petID := createPetForPurchase(t, db, sellerID, false)

		_, err := repo.Create(&models.PurchaseRequest{
			PetID:   petID,
			BuyerID: buyerID,
			Status:  "pending",
		})
		if err == nil {
			t.Fatal("expected create error for inactive pet")
		}
	})

	t.Run("GetAll GetByBuyerID GetByPetID", func(t *testing.T) {
		clearPurchaseTables(t, db)

		sellerID := createUserForPurchase(t, db, "seller")
		buyer1 := createUserForPurchase(t, db, "buyer")
		buyer2 := createUserForPurchase(t, db, "buyer")
		pet1 := createPetForPurchase(t, db, sellerID, true)
		pet2 := createPetForPurchase(t, db, sellerID, true)

		r1, err := repo.Create(&models.PurchaseRequest{PetID: pet1, BuyerID: buyer1})
		if err != nil {
			t.Fatalf("create r1: %v", err)
		}
		_, _ = r1, buyer2

		if _, err := repo.Create(&models.PurchaseRequest{PetID: pet1, BuyerID: buyer2}); err != nil {
			t.Fatalf("create r2: %v", err)
		}

		if _, err := repo.Create(&models.PurchaseRequest{PetID: pet2, BuyerID: buyer2}); err != nil {
			t.Fatalf("create r3: %v", err)
		}

		all, err := repo.GetAll()
		if err != nil {
			t.Fatalf("get all: %v", err)
		}
		if len(all) != 3 {
			t.Fatalf("expected 3 requests, got %d", len(all))
		}

		byBuyer1, err := repo.GetByBuyerID(buyer1)
		if err != nil {
			t.Fatalf("get by buyer: %v", err)
		}
		if len(byBuyer1) != 1 {
			t.Fatalf("expected 1 request for buyer1, got %d", len(byBuyer1))
		}

		byPet1, err := repo.GetByPetID(pet1)
		if err != nil {
			t.Fatalf("get by pet: %v", err)
		}
		if len(byPet1) != 2 {
			t.Fatalf("expected 2 requests for pet1, got %d", len(byPet1))
		}
	})

	t.Run("GetBySellerID", func(t *testing.T) {
		clearPurchaseTables(t, db)

		seller1 := createUserForPurchase(t, db, "seller")
		seller2 := createUserForPurchase(t, db, "seller")
		buyer := createUserForPurchase(t, db, "buyer")

		pet1 := createPetForPurchase(t, db, seller1, true)
		pet2 := createPetForPurchase(t, db, seller2, true)

		if _, err := repo.Create(&models.PurchaseRequest{PetID: pet1, BuyerID: buyer}); err != nil {
			t.Fatalf("create req for seller1: %v", err)
		}
		if _, err := repo.Create(&models.PurchaseRequest{PetID: pet2, BuyerID: buyer}); err != nil {
			t.Fatalf("create req for seller2: %v", err)
		}

		incoming, err := repo.GetBySellerID(seller1)
		if err != nil {
			t.Fatalf("get incoming: %v", err)
		}
		if len(incoming) != 1 {
			t.Fatalf("expected 1 incoming request, got %d", len(incoming))
		}
		if incoming[0].PetID != pet1 {
			t.Fatalf("expected pet %s, got %s", pet1, incoming[0].PetID)
		}
	})

	t.Run("UpdateStatusBySeller approves and closes sale", func(t *testing.T) {
		clearPurchaseTables(t, db)

		sellerID := createUserForPurchase(t, db, "seller")
		buyer1 := createUserForPurchase(t, db, "buyer")
		buyer2 := createUserForPurchase(t, db, "buyer")
		petID := createPetForPurchase(t, db, sellerID, true)

		req1, err := repo.Create(&models.PurchaseRequest{PetID: petID, BuyerID: buyer1})
		if err != nil {
			t.Fatalf("create req1: %v", err)
		}
		req2, err := repo.Create(&models.PurchaseRequest{PetID: petID, BuyerID: buyer2})
		if err != nil {
			t.Fatalf("create req2: %v", err)
		}

		updated, err := repo.UpdateStatusBySeller(req1.ID, sellerID, "approved")
		if err != nil {
			t.Fatalf("approve request: %v", err)
		}
		if updated.Status != "approved" {
			t.Fatalf("expected approved, got %s", updated.Status)
		}

		second, err := repo.GetByID(req2.ID)
		if err != nil {
			t.Fatalf("get second request: %v", err)
		}
		if second.Status != "rejected" {
			t.Fatalf("expected rejected, got %s", second.Status)
		}

		var isActive bool
		if err := db.QueryRow(`SELECT is_active FROM pet WHERE id = $1`, petID).Scan(&isActive); err != nil {
			t.Fatalf("select pet is_active: %v", err)
		}
		if isActive {
			t.Fatal("expected pet to become inactive")
		}
	})

	t.Run("UpdateStatusBySeller forbidden", func(t *testing.T) {
		clearPurchaseTables(t, db)

		ownerSeller := createUserForPurchase(t, db, "seller")
		otherSeller := createUserForPurchase(t, db, "seller")
		buyer := createUserForPurchase(t, db, "buyer")
		petID := createPetForPurchase(t, db, ownerSeller, true)

		req, err := repo.Create(&models.PurchaseRequest{PetID: petID, BuyerID: buyer})
		if err != nil {
			t.Fatalf("create req: %v", err)
		}

		_, err = repo.UpdateStatusBySeller(req.ID, otherSeller, "approved")
		if !errors.Is(err, models.ErrPurchaseRequestForbidden) {
			t.Fatalf("expected ErrPurchaseRequestForbidden, got %v", err)
		}
	})

	t.Run("DeleteByBuyer", func(t *testing.T) {
		clearPurchaseTables(t, db)

		sellerID := createUserForPurchase(t, db, "seller")
		buyerID := createUserForPurchase(t, db, "buyer")
		petID := createPetForPurchase(t, db, sellerID, true)

		req, err := repo.Create(&models.PurchaseRequest{PetID: petID, BuyerID: buyerID})
		if err != nil {
			t.Fatalf("create req: %v", err)
		}

		if err := repo.DeleteByBuyer(req.ID, buyerID); err != nil {
			t.Fatalf("delete by buyer: %v", err)
		}

		_, err = repo.GetByID(req.ID)
		if !errors.Is(err, models.ErrPurchaseRequestNotFound) {
			t.Fatalf("expected ErrPurchaseRequestNotFound, got %v", err)
		}
	})

	t.Run("DeleteByBuyer forbidden", func(t *testing.T) {
		clearPurchaseTables(t, db)

		sellerID := createUserForPurchase(t, db, "seller")
		ownerBuyer := createUserForPurchase(t, db, "buyer")
		otherBuyer := createUserForPurchase(t, db, "buyer")
		petID := createPetForPurchase(t, db, sellerID, true)

		req, err := repo.Create(&models.PurchaseRequest{PetID: petID, BuyerID: ownerBuyer})
		if err != nil {
			t.Fatalf("create req: %v", err)
		}

		err = repo.DeleteByBuyer(req.ID, otherBuyer)
		if !errors.Is(err, models.ErrPurchaseRequestForbidden) {
			t.Fatalf("expected ErrPurchaseRequestForbidden, got %v", err)
		}
	})
}

func setupTestPostgresPurchaseRequest(t *testing.T, ctx context.Context) (*sql.DB, func()) {
	t.Helper()

	if dsn := os.Getenv("INTEGRATION_DB_DSN"); dsn != "" {
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			t.Fatalf("open db: %v", err)
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
			t.Fatalf("ping db: %v", pingErr)
		}

		initPurchaseRequestSchema(ctx, db)
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

	connStr, err := pg.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = pg.Terminate(context.Background())
		t.Fatalf("connection string: %v", err)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		_ = pg.Terminate(context.Background())
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
		_ = pg.Terminate(context.Background())
		t.Fatalf("ping db: %v", pingErr)
	}

	initPurchaseRequestSchema(ctx, db)

	return db, func() {
		_ = db.Close()
		_ = pg.Terminate(context.Background())
	}
}

func initPurchaseRequestSchema(ctx context.Context, db *sql.DB) {
	db.ExecContext(ctx, `
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

		CREATE TABLE IF NOT EXISTS pet (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			vet_passport_id UUID,
			seller_id UUID REFERENCES users(id),
			pet_name VARCHAR(255) NOT NULL,
			species VARCHAR(50) NOT NULL,
			pet_age INT NOT NULL,
			color VARCHAR(50),
			pet_gender VARCHAR(20),
			breed VARCHAR(255),
			pedigree BOOLEAN DEFAULT FALSE,
			good_with_children BOOLEAN DEFAULT TRUE,
			good_with_animals BOOLEAN DEFAULT TRUE,
			pet_description TEXT,
			is_active BOOLEAN DEFAULT TRUE,
			price DECIMAL(10,2)
		);

		CREATE TABLE IF NOT EXISTS purchase_request (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			pet_id UUID NOT NULL REFERENCES pet(id),
			buyer_id UUID NOT NULL REFERENCES users(id),
			status VARCHAR(50) DEFAULT 'pending',
			request_date TIMESTAMP DEFAULT NOW()
		);
	`)
}

func clearPurchaseTables(t *testing.T, db *sql.DB) {
	t.Helper()
	if _, err := db.Exec(`TRUNCATE TABLE purchase_request, pet, users`); err != nil {
		t.Fatalf("truncate tables: %v", err)
	}
}

func createUserForPurchase(t *testing.T, db *sql.DB, status string) uuid.UUID {
	t.Helper()

	var id uuid.UUID
	err := db.QueryRow(
		`INSERT INTO users (fio, telephone_number, city, user_login, user_password, status, user_description)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id`,
		"Test User "+status,
		"+79990000000",
		"Moscow",
		"user_"+uuid.NewString(),
		"hash",
		status,
		"desc",
	).Scan(&id)
	if err != nil {
		t.Fatalf("insert user: %v", err)
	}

	return id
}

func createPetForPurchase(t *testing.T, db *sql.DB, sellerID uuid.UUID, isActive bool) uuid.UUID {
	t.Helper()

	var id uuid.UUID
	err := db.QueryRow(
		`INSERT INTO pet (
			vet_passport_id, seller_id, pet_name, species, pet_age,
			color, pet_gender, breed, pedigree, good_with_children,
			good_with_animals, pet_description, is_active, price
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
		RETURNING id`,
		uuid.New(),
		sellerID,
		"PetName",
		"Кошка",
		2,
		"Черный",
		"Девочка",
		"Британская",
		false,
		true,
		true,
		"desc",
		isActive,
		1000.0,
	).Scan(&id)
	if err != nil {
		t.Fatalf("insert pet: %v", err)
	}

	return id
}
