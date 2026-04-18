package repository

import (
	"database/sql"

	"getapet-backend/internal/models"
	"github.com/google/uuid"
)

type VetPassportRepository struct {
	db *sql.DB
}

func NewVetPassportRepository(db *sql.DB) *VetPassportRepository {
	return &VetPassportRepository{db: db}
}

// func (r *VetPassportRepository) Create(_ *models.VetPassport) (*models.VetPassport, error) {
// 	return nil, errors.New("not implemented")
// }

func (r *VetPassportRepository) Create(p *models.VetPassport) (*models.VetPassport, error) {
	query := `
	INSERT INTO vet_passport (
		chipping, sterilization, health_issues, vaccinations, parasite_treatments
	) VALUES ($1,$2,$3,$4,$5)
	RETURNING id
	`

	err := r.db.QueryRow(
		query,
		p.Chipping,
		p.Sterilization,
		p.HealthIssues,
		p.Vaccinations,
		p.ParasiteTreatments,
	).Scan(&p.ID)

	if err != nil {
		return nil, err
	}

	return p, nil
}

func (r *VetPassportRepository) GetAll() ([]models.VetPassport, error) {
	query := `
	SELECT id, chipping, sterilization, health_issues, vaccinations, parasite_treatments
	FROM vet_passport
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var passports []models.VetPassport

	for rows.Next() {
		var p models.VetPassport

		err := rows.Scan(
			&p.ID,
			&p.Chipping,
			&p.Sterilization,
			&p.HealthIssues,
			&p.Vaccinations,
			&p.ParasiteTreatments,
		)
		if err != nil {
			return nil, err
		}

		passports = append(passports, p)
	}

	return passports, nil
}

func (r *VetPassportRepository) GetByID(id uuid.UUID) (*models.VetPassport, error) {
	query := `
	SELECT id, chipping, sterilization, health_issues, vaccinations, parasite_treatments
	FROM vet_passport
	WHERE id = $1
	`

	var p models.VetPassport

	err := r.db.QueryRow(query, id).Scan(
		&p.ID,
		&p.Chipping,
		&p.Sterilization,
		&p.HealthIssues,
		&p.Vaccinations,
		&p.ParasiteTreatments,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrVetPassportNotFound
		}
		return nil, err
	}

	return &p, nil
}

func (r *VetPassportRepository) Update(id uuid.UUID, p *models.VetPassport) (*models.VetPassport, error) {
	query := `
	UPDATE vet_passport SET
		chipping = $1,
		sterilization = $2,
		health_issues = $3,
		vaccinations = $4,
		parasite_treatments = $5
	WHERE id = $6
	RETURNING id
	`

	var updatedID uuid.UUID

	err := r.db.QueryRow(
		query,
		p.Chipping,
		p.Sterilization,
		p.HealthIssues,
		p.Vaccinations,
		p.ParasiteTreatments,
		id,
	).Scan(&updatedID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrVetPassportNotFound
		}
		return nil, err
	}

	p.ID = updatedID
	return p, nil
}

func (r *VetPassportRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM vet_passport WHERE id = $1`

	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return models.ErrVetPassportNotFound
	}

	return nil
}
