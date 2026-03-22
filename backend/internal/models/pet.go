package models

import "errors"

type Pet struct {
	ID               int     `json:"id" db:"id"`
	VetPassportID    int     `json:"vet_passport_id" db:"vet_passport_id"`
	SellerID         int     `json:"seller_id" db:"seller_id"`
	PetName          string  `json:"pet_name" db:"pet_name"`
	Species          string  `json:"species" db:"species"`
	PetAge           int     `json:"pet_age" db:"pet_age"`
	Color            string  `json:"color" db:"color"`
	PetGender        string  `json:"pet_gender" db:"pet_gender"`
	Breed            string  `json:"breed" db:"breed"`
	Pedigree         bool    `json:"pedigree" db:"pedigree"`
	GoodWithChildren bool    `json:"good_with_children" db:"good_with_children"`
	GoodWithAnimals  bool    `json:"good_with_animals" db:"good_with_animals"`
	PetDescription   string  `json:"pet_description" db:"pet_description"`
	IsActive         bool    `json:"is_active" db:"is_active"`
	Price            float64 `json:"price" db:"price"`
}

type PetRepository interface {
	Create(pet *Pet) (*Pet, error)
	GetAll() ([]Pet, error)
	GetByID(id int) (*Pet, error)
	GetBySellerID(sellerID int) ([]Pet, error)
	Update(id int, pet *Pet) (*Pet, error)
	Delete(id int) error
	CheckBelonging(id, sellerID int) (bool, error)
}

type PetService interface {
	Create(pet *Pet) (*Pet, error)
	GetAll() ([]Pet, error)
	GetByID(id int) (*Pet, error)
	GetBySellerID(sellerID int) ([]Pet, error)
	Update(id, sellerID int, pet *Pet) (*Pet, error)
	Delete(id, sellerID int) error
}

var (
	ErrPetNotFound = errors.New("pet not found")
)
