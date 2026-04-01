package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"getapet-backend/internal/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type PurchaseRequestRepository struct {
	db *sql.DB
}

const (
	uqPurchaseRequestPetBuyer  = "uq_purchase_request_pet_buyer"
	uxOneApprovedRequestPerPet = "ux_purchase_request_one_approved_per_pet"
)

func mapPurchaseRequestConstraintError(err error) error {
	var pqErr *pq.Error
	if !errors.As(err, &pqErr) {
		return err
	}

	switch pqErr.Constraint {
	case uqPurchaseRequestPetBuyer:
		return models.ErrPurchaseRequestDuplicatePetBuyer
	case uxOneApprovedRequestPerPet:
		return models.ErrPurchaseRequestAlreadyApprovedForPet
	}

	if pqErr.Code == "23505" {
		return models.ErrPurchaseRequestUniqueViolation
	}

	return err
}

func NewPurchaseRequestRepository(db *sql.DB) *PurchaseRequestRepository {
	return &PurchaseRequestRepository{db: db}
}

func (r *PurchaseRequestRepository) Create(request *models.PurchaseRequest) (*models.PurchaseRequest, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	var isActive bool
	//FOR UPDATE ставит блокировку на выбранные строки до конца текущей транзакции
	err = tx.QueryRow(`SELECT is_active FROM pet WHERE id = $1 FOR UPDATE`, request.PetID).Scan(&isActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrPetNotFound
		}
		return nil, err
	}

	if !isActive {
		return nil, models.ErrPurchaseRequestPetNotAvailable
	}

	if strings.TrimSpace(request.Status) == "" {
		request.Status = "pending"
	}

	const query = `
		INSERT INTO purchase_request (pet_id, buyer_id, status)
		VALUES ($1, $2, $3)
		RETURNING id, pet_id, buyer_id, status, request_date
	`

	err = tx.QueryRow(query, request.PetID, request.BuyerID, request.Status).Scan(
		&request.ID,
		&request.PetID,
		&request.BuyerID,
		&request.Status,
		&request.RequestDate,
	)
	if err != nil {
		return nil, mapPurchaseRequestConstraintError(err)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return request, nil
}

func (r *PurchaseRequestRepository) GetAll() ([]models.PurchaseRequest, error) {
	const query = `
		SELECT id, pet_id, buyer_id, status, request_date
		FROM purchase_request
		ORDER BY request_date DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requests := make([]models.PurchaseRequest, 0)
	for rows.Next() {
		var req models.PurchaseRequest
		if err := rows.Scan(
			&req.ID,
			&req.PetID,
			&req.BuyerID,
			&req.Status,
			&req.RequestDate,
		); err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return requests, nil
}

func (r *PurchaseRequestRepository) GetByID(id uuid.UUID) (*models.PurchaseRequest, error) {
	const query = `
		SELECT id, pet_id, buyer_id, status, request_date
		FROM purchase_request
		WHERE id = $1
	`

	var req models.PurchaseRequest
	err := r.db.QueryRow(query, id).Scan(
		&req.ID,
		&req.PetID,
		&req.BuyerID,
		&req.Status,
		&req.RequestDate,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrPurchaseRequestNotFound
		}
		return nil, err
	}

	return &req, nil
}

func (r *PurchaseRequestRepository) GetByBuyerID(buyerID uuid.UUID) ([]models.PurchaseRequest, error) {
	const query = `
		SELECT id, pet_id, buyer_id, status, request_date
		FROM purchase_request
		WHERE buyer_id = $1
		ORDER BY request_date DESC
	`

	rows, err := r.db.Query(query, buyerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requests := make([]models.PurchaseRequest, 0)
	for rows.Next() {
		var req models.PurchaseRequest
		if err := rows.Scan(
			&req.ID,
			&req.PetID,
			&req.BuyerID,
			&req.Status,
			&req.RequestDate,
		); err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return requests, nil
}

func (r *PurchaseRequestRepository) GetBySellerID(sellerID uuid.UUID) ([]models.PurchaseRequest, error) {
	const query = `
		SELECT pr.id, pr.pet_id, pr.buyer_id, pr.status, pr.request_date
		FROM purchase_request pr
		JOIN pet p ON p.id = pr.pet_id
		WHERE p.seller_id = $1
		ORDER BY pr.request_date DESC
	`

	rows, err := r.db.Query(query, sellerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requests := make([]models.PurchaseRequest, 0)
	for rows.Next() {
		var req models.PurchaseRequest
		if err := rows.Scan(
			&req.ID,
			&req.PetID,
			&req.BuyerID,
			&req.Status,
			&req.RequestDate,
		); err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return requests, nil
}

func (r *PurchaseRequestRepository) GetByPetID(petID uuid.UUID) ([]models.PurchaseRequest, error) {
	const query = `
		SELECT id, pet_id, buyer_id, status, request_date
		FROM purchase_request
		WHERE pet_id = $1
		ORDER BY request_date DESC
	`

	rows, err := r.db.Query(query, petID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requests := make([]models.PurchaseRequest, 0)
	for rows.Next() {
		var req models.PurchaseRequest
		if err := rows.Scan(
			&req.ID,
			&req.PetID,
			&req.BuyerID,
			&req.Status,
			&req.RequestDate,
		); err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return requests, nil
}

func (r *PurchaseRequestRepository) requestExistsTx(tx *sql.Tx, id uuid.UUID) (bool, error) {
	var exists bool
	if err := tx.QueryRow(`SELECT EXISTS (SELECT 1 FROM purchase_request WHERE id = $1)`, id).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (r *PurchaseRequestRepository) UpdateStatus(id uuid.UUID, status string) (*models.PurchaseRequest, error) {
	normalizedStatus := strings.ToLower(strings.TrimSpace(status))
	if normalizedStatus == "" {
		return nil, models.ErrPurchaseRequestStatusRequired
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	var req models.PurchaseRequest
	err = tx.QueryRow(
		`SELECT id, pet_id, buyer_id, status, request_date FROM purchase_request WHERE id = $1 FOR UPDATE`,
		id,
	).Scan(
		&req.ID,
		&req.PetID,
		&req.BuyerID,
		&req.Status,
		&req.RequestDate,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrPurchaseRequestNotFound
		}
		return nil, err
	}

	previousStatus := req.Status

	if normalizedStatus == "approved" {
		var isActive bool
		err = tx.QueryRow(`SELECT is_active FROM pet WHERE id = $1 FOR UPDATE`, req.PetID).Scan(&isActive)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, models.ErrPetNotFound
			}
			return nil, err
		}
		if !isActive && req.Status != "approved" {
			return nil, models.ErrPurchaseRequestPetNotAvailable
		}
	}

	// обновление статуса требуемой заявки
	err = tx.QueryRow(
		`UPDATE purchase_request
		 SET status = $1
		 WHERE id = $2
		 RETURNING id, pet_id, buyer_id, status, request_date`,
		normalizedStatus,
		id,
	).Scan(
		&req.ID,
		&req.PetID,
		&req.BuyerID,
		&req.Status,
		&req.RequestDate,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrPurchaseRequestNotFound
		}
		return nil, mapPurchaseRequestConstraintError(err)
	}

	switch normalizedStatus {
	/* если заявка одобрена, то оставшиеся заявки со статусом pending надо перевести в статус rejected
	и для сделать объявление питомца неактивным*/
	case "approved":
		_, err = tx.Exec(
			`UPDATE purchase_request
			 SET status = 'rejected'
			 WHERE pet_id = $1 AND id <> $2 AND status = 'pending'`,
			req.PetID,
			req.ID,
		)
		if err != nil {
			return nil, err
		}

		if _, err = tx.Exec(`UPDATE pet SET is_active = false WHERE id = $1`, req.PetID); err != nil {
			return nil, err
		}
	default:
		/* если заявка НЕ одобрена, а меняется на pending/regected */
		if previousStatus != "approved" {
			//было заявка одобрена, так что ничего в таблице pet менять не надо, то есть peisactive уже и так равен false
			break
		}
		/* заявка была до этого одобрена, но стала rejected/penging => надо проверить есть ли
		approved заявки: если есть, то оставить is_active у питомца равным false, иначе - равным true*/
		var hasApproved bool
		//проверка на наличие активных других заявок
		err = tx.QueryRow(
			`SELECT EXISTS (SELECT 1 FROM purchase_request WHERE pet_id = $1 AND status = 'approved')`,
			req.PetID,
		).Scan(&hasApproved)
		if err != nil {
			return nil, err
		}
		//делаем активные объявления неактивными
		if _, err = tx.Exec(`UPDATE pet SET is_active = $1 WHERE id = $2`, !hasApproved, req.PetID); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &req, nil
}

func (r *PurchaseRequestRepository) UpdateStatusBySeller(id uuid.UUID, sellerID uuid.UUID, status string) (*models.PurchaseRequest, error) {
	normalizedStatus := strings.ToLower(strings.TrimSpace(status))
	if normalizedStatus == "" {
		return nil, models.ErrPurchaseRequestStatusRequired
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	var req models.PurchaseRequest
	var isActive bool
	err = tx.QueryRow(
		`SELECT pr.id, pr.pet_id, pr.buyer_id, pr.status, pr.request_date, p.is_active
		 FROM purchase_request pr
		 JOIN pet p ON p.id = pr.pet_id
		 WHERE pr.id = $1 AND p.seller_id = $2
		 FOR UPDATE OF pr, p`,
		id,
		sellerID,
	).Scan(
		&req.ID,
		&req.PetID,
		&req.BuyerID,
		&req.Status,
		&req.RequestDate,
		&isActive,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exists, existsErr := r.requestExistsTx(tx, id)
			if existsErr != nil {
				return nil, existsErr
			}
			if exists {
				return nil, models.ErrPurchaseRequestForbidden
			}
			return nil, models.ErrPurchaseRequestNotFound
		}
		return nil, err
	}

	previousStatus := req.Status
	if normalizedStatus == "approved" && !isActive && req.Status != "approved" {
		return nil, models.ErrPurchaseRequestPetNotAvailable
	}

	err = tx.QueryRow(
		`UPDATE purchase_request
		 SET status = $1
		 WHERE id = $2
		 RETURNING id, pet_id, buyer_id, status, request_date`,
		normalizedStatus,
		id,
	).Scan(
		&req.ID,
		&req.PetID,
		&req.BuyerID,
		&req.Status,
		&req.RequestDate,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrPurchaseRequestNotFound
		}
		return nil, mapPurchaseRequestConstraintError(err)
	}

	switch normalizedStatus {
	case "approved":
		_, err = tx.Exec(
			`UPDATE purchase_request
			 SET status = 'rejected'
			 WHERE pet_id = $1 AND id <> $2 AND status = 'pending'`,
			req.PetID,
			req.ID,
		)
		if err != nil {
			return nil, err
		}

		if _, err = tx.Exec(`UPDATE pet SET is_active = false WHERE id = $1`, req.PetID); err != nil {
			return nil, err
		}
	default:
		if previousStatus != "approved" {
			break
		}
		var hasApproved bool
		err = tx.QueryRow(
			`SELECT EXISTS (SELECT 1 FROM purchase_request WHERE pet_id = $1 AND status = 'approved')`,
			req.PetID,
		).Scan(&hasApproved)
		if err != nil {
			return nil, err
		}

		if _, err = tx.Exec(`UPDATE pet SET is_active = $1 WHERE id = $2`, !hasApproved, req.PetID); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &req, nil
}

func (r *PurchaseRequestRepository) Delete(id uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	var petID uuid.UUID
	var status string
	err = tx.QueryRow(
		`SELECT pet_id, status FROM purchase_request WHERE id = $1 FOR UPDATE`,
		id,
	).Scan(&petID, &status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrPurchaseRequestNotFound
		}
		return err
	}

	result, err := tx.Exec(`DELETE FROM purchase_request WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return models.ErrPurchaseRequestNotFound
	}

	if status == "approved" {
		var hasApproved bool
		err = tx.QueryRow(
			`SELECT EXISTS (SELECT 1 FROM purchase_request WHERE pet_id = $1 AND status = 'approved')`,
			petID,
		).Scan(&hasApproved)
		if err != nil {
			return err
		}

		if _, err = tx.Exec(`UPDATE pet SET is_active = $1 WHERE id = $2`, !hasApproved, petID); err != nil {
			return fmt.Errorf("не удалось обновить доступность питомца: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *PurchaseRequestRepository) DeleteByBuyer(id uuid.UUID, buyerID uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	var petID uuid.UUID
	var status string
	err = tx.QueryRow(
		`SELECT pet_id, status FROM purchase_request WHERE id = $1 AND buyer_id = $2 FOR UPDATE`,
		id,
		buyerID,
	).Scan(&petID, &status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exists, existsErr := r.requestExistsTx(tx, id)
			if existsErr != nil {
				return existsErr
			}
			if exists {
				return models.ErrPurchaseRequestForbidden
			}
			return models.ErrPurchaseRequestNotFound
		}
		return err
	}

	result, err := tx.Exec(`DELETE FROM purchase_request WHERE id = $1 AND buyer_id = $2`, id, buyerID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return models.ErrPurchaseRequestNotFound
	}

	if status == "approved" {
		var hasApproved bool
		err = tx.QueryRow(
			`SELECT EXISTS (SELECT 1 FROM purchase_request WHERE pet_id = $1 AND status = 'approved')`,
			petID,
		).Scan(&hasApproved)
		if err != nil {
			return err
		}

		if _, err = tx.Exec(`UPDATE pet SET is_active = $1 WHERE id = $2`, !hasApproved, petID); err != nil {
			return fmt.Errorf("не удалось обновить доступность питомца: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
