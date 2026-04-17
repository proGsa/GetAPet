package repository

import (
	"database/sql"
	"errors"

	"getapet-backend/internal/models"
	"github.com/google/uuid"
)

var ErrPetNotFound = errors.New("pet not found")

type PetRepository struct {
	db *sql.DB
}

func NewPetRepository(db *sql.DB) *PetRepository {
	return &PetRepository{db: db}
}

func (r *PetRepository) Create(p *models.Pet) (*models.Pet, error) {
	query := `
	INSERT INTO pet (
		vet_passport_id, seller_id, pet_name, species, pet_age,
		color, pet_gender, breed, pedigree, good_with_children,
		good_with_animals, pet_description, is_active, price
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
	RETURNING id
	`

	err := r.db.QueryRow(
		query,
		p.VetPassportID,
		p.SellerID,
		p.PetName,
		p.Species,
		p.PetAge,
		p.Color,
		p.PetGender,
		p.Breed,
		p.Pedigree,
		p.GoodWithChildren,
		p.GoodWithAnimals,
		p.PetDescription,
		p.IsActive,
		p.Price,
	).Scan(&p.ID)

	if err != nil {
		return nil, err
	}

	return p, nil
}

func (r *PetRepository) GetAll() ([]models.Pet, error) {
	query := `SELECT id, vet_passport_id, seller_id, pet_name, species, pet_age,
				color, pet_gender, breed, pedigree, good_with_children,
				good_with_animals, pet_description, is_active, price
				FROM pet`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pets []models.Pet

	for rows.Next() {
		var p models.Pet

		err := rows.Scan(
			&p.ID,
			&p.VetPassportID,
			&p.SellerID,
			&p.PetName,
			&p.Species,
			&p.PetAge,
			&p.Color,
			&p.PetGender,
			&p.Breed,
			&p.Pedigree,
			&p.GoodWithChildren,
			&p.GoodWithAnimals,
			&p.PetDescription,
			&p.IsActive,
			&p.Price,
		)
		if err != nil {
			return nil, err
		}

		pets = append(pets, p)
	}

	return pets, nil
}

//?
// func (r *PetRepository) GetByID(_ uuid.UUID) (*models.Pet, error) {
// 	return nil, errors.New("not implemented")
// }
//?
// func (r *PetRepository) GetBySellerID(_ uuid.UUID) ([]models.Pet, error) {
// 	return nil, errors.New("not implemented")
// }

func (r *PetRepository) GetBySellerID(sellerID uuid.UUID) ([]models.Pet, error) {
	query := `SELECT id, vet_passport_id, seller_id, pet_name, species, pet_age,
       		color, pet_gender, breed, pedigree, good_with_children,
       		good_with_animals, pet_description, is_active, price
			FROM pet 
			WHERE seller_id = $1`

	rows, err := r.db.Query(query, sellerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pets []models.Pet

	for rows.Next() {
		var p models.Pet

		err := rows.Scan(
			&p.ID,
			&p.VetPassportID,
			&p.SellerID,
			&p.PetName,
			&p.Species,
			&p.PetAge,
			&p.Color,
			&p.PetGender,
			&p.Breed,
			&p.Pedigree,
			&p.GoodWithChildren,
			&p.GoodWithAnimals,
			&p.PetDescription,
			&p.IsActive,
			&p.Price,
		)
		if err != nil {
			return nil, err
		}

		pets = append(pets, p)
	}

	return pets, nil
}

//?
// func (r *PetRepository) Update(_ uuid.UUID, _ *models.Pet) (*models.Pet, error) {
// 	return nil, errors.New("not implemented")
// }

func (r *PetRepository) Update(id uuid.UUID, p *models.Pet) (*models.Pet, error) {
	query := `
	UPDATE pet SET
		pet_name = $1,
		species = $2,
		pet_age = $3,
		color = $4,
		pet_gender = $5,
		breed = $6,
		pedigree = $7,
		good_with_children = $8,
		good_with_animals = $9,
		pet_description = $10,
		is_active = $11,
		price = $12
	WHERE id = $13
	RETURNING id
	`

	var updatedID uuid.UUID

	err := r.db.QueryRow(
		query,
		p.PetName,
		p.Species,
		p.PetAge,
		p.Color,
		p.PetGender,
		p.Breed,
		p.Pedigree,
		p.GoodWithChildren,
		p.GoodWithAnimals,
		p.PetDescription,
		p.IsActive,
		p.Price,
		id,
	).Scan(&updatedID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrPetNotFound
		}
		return nil, err
	}

	p.ID = updatedID
	return p, nil
}

func (r *PetRepository) GetByID(id uuid.UUID) (*models.Pet, error) {
	query := `
			SELECT id, vet_passport_id, seller_id, pet_name, species, pet_age,
				color, pet_gender, breed, pedigree, good_with_children,
				good_with_animals, pet_description, is_active, price
			FROM pet
			WHERE id = $1
			`

	var p models.Pet

	err := r.db.QueryRow(query, id).Scan(
		&p.ID,
		&p.VetPassportID,
		&p.SellerID,
		&p.PetName,
		&p.Species,
		&p.PetAge,
		&p.Color,
		&p.PetGender,
		&p.Breed,
		&p.Pedigree,
		&p.GoodWithChildren,
		&p.GoodWithAnimals,
		&p.PetDescription,
		&p.IsActive,
		&p.Price,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrPetNotFound
		}
		return nil, err
	}

	return &p, nil
}

//?
// func (r *PetRepository) CheckBelonging(_, _ uuid.UUID) (bool, error) {
// 	return false, errors.New("not implemented")
// }

func (r *PetRepository) CheckBelonging(petID, sellerID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS (
		SELECT 1 FROM pet WHERE id = $1 AND seller_id = $2
	)`

	var exists bool

	err := r.db.QueryRow(query, petID, sellerID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *PetRepository) Delete(id uuid.UUID) error {
	query := `
		DELETE FROM pet
		WHERE id = $1
	`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrPetNotFound
	}

	return nil
}
