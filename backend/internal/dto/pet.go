package dto

import (
	"getapet-backend/internal/models"
	"github.com/google/uuid"
)

type CreatePetRequest struct {
	VetPassportID    string  `json:"vet_passport_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	PetName          string  `json:"pet_name" validate:"required,min=1,max=100" example:"Барсик"`
	Species          string  `json:"species" validate:"required,oneof=cat dog bird other" example:"cat"`
	PetAge           int     `json:"pet_age" validate:"required,gte=0,lte=50" example:"2"`
	Color            string  `json:"color" validate:"required,max=50" example:"black"`
	PetGender        string  `json:"pet_gender" validate:"required,oneof=male female" example:"male"`
	Breed            string  `json:"breed" validate:"omitempty,max=100" example:"british"`
	Pedigree         bool    `json:"pedigree" example:"true"`
	GoodWithChildren bool    `json:"good_with_children" example:"true"`
	GoodWithAnimals  bool    `json:"good_with_animals" example:"false"`
	PetDescription   string  `json:"pet_description" validate:"omitempty,max=1000" example:"Очень дружелюбный и активный питомец"`
	Price            float64 `json:"price" validate:"required,gte=0" example:"15000"`
}
type CreatePetResponse struct {
	ID string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type UpdatePetRequest struct {
	PetName          string  `json:"pet_name" validate:"required,min=1,max=100" example:"Барсик"`
	Species          string  `json:"species" validate:"required,oneof=cat dog bird other" example:"cat"`
	PetAge           int     `json:"pet_age" validate:"required,gte=0,lte=50" example:"3"`
	Color            string  `json:"color" validate:"required,max=50" example:"black"`
	PetGender        string  `json:"pet_gender" validate:"required,oneof=male female" example:"male"`
	Breed            string  `json:"breed" validate:"omitempty,max=100" example:"british"`
	Pedigree         bool    `json:"pedigree" example:"true"`
	GoodWithChildren bool    `json:"good_with_children" example:"true"`
	GoodWithAnimals  bool    `json:"good_with_animals" example:"false"`
	PetDescription   string  `json:"pet_description" validate:"omitempty,max=1000" example:"Спокойный и ласковый"`
	IsActive         bool    `json:"is_active" example:"true"`
	Price            float64 `json:"price" validate:"required,gte=0" example:"18000"`
}
type UpdatePetResponse struct {
	ID               string  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	PetName          string  `json:"pet_name" validate:"required,min=1,max=100" example:"Барсик"`
	Species          string  `json:"species" validate:"required,oneof=cat dog bird other" example:"cat"`
	PetAge           int     `json:"pet_age" validate:"required,gte=0,lte=50" example:"3"`
	Color            string  `json:"color" validate:"required,max=50" example:"black"`
	PetGender        string  `json:"pet_gender" validate:"required,oneof=male female" example:"male"`
	Breed            string  `json:"breed" validate:"omitempty,max=100" example:"british"`
	Pedigree         bool    `json:"pedigree" example:"true"`
	GoodWithChildren bool    `json:"good_with_children" example:"true"`
	GoodWithAnimals  bool    `json:"good_with_animals" example:"false"`
	PetDescription   string  `json:"pet_description" validate:"omitempty,max=1000" example:"Спокойный и ласковый"`
	IsActive         bool    `json:"is_active" example:"true"`
	Price            float64 `json:"price" validate:"required,gte=0" example:"18000"`
}

// type UpdateUserResponse struct {
// 	ID string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
// }

type PetResponse struct {
	ID               string  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	VetPassportID    string  `json:"vet_passport_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	SellerID         string  `json:"seller_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	PetName          string  `json:"pet_name" example:"Барсик"`
	Species          string  `json:"species" example:"cat"`
	PetAge           int     `json:"pet_age" example:"2"`
	Color            string  `json:"color" example:"black"`
	PetGender        string  `json:"pet_gender" example:"male"`
	Breed            string  `json:"breed" example:"british"`
	Pedigree         bool    `json:"pedigree" example:"true"`
	GoodWithChildren bool    `json:"good_with_children" example:"true"`
	GoodWithAnimals  bool    `json:"good_with_animals" example:"false"`
	PetDescription   string  `json:"pet_description" example:"Очень дружелюбный"`
	IsActive         bool    `json:"is_active" example:"true"`
	Price            float64 `json:"price" example:"15000"`
}

func CreatePetRequestFromDTO(req CreatePetRequest, sellerID uuid.UUID) (models.Pet, error) {
	vetPassportID, err := uuid.Parse(req.VetPassportID)
	if err != nil {
		return models.Pet{}, err
	}
	return models.Pet{
		VetPassportID:    vetPassportID,
		SellerID:         sellerID,
		PetName:          req.PetName,
		Species:          req.Species,
		PetAge:           req.PetAge,
		Color:            req.Color,
		PetGender:        req.PetGender,
		Breed:            req.Breed,
		Pedigree:         req.Pedigree,
		GoodWithChildren: req.GoodWithChildren,
		GoodWithAnimals:  req.GoodWithAnimals,
		PetDescription:   req.PetDescription,
		IsActive:         true,
		Price:            req.Price,
	}, nil
}

func CreatePetResponseFromDomain(p models.Pet) CreatePetResponse {
	return CreatePetResponse{
		ID: p.ID.String(),
	}
}

func UpdatePetResponseFromDomain(p models.Pet) UpdatePetResponse {
	return UpdatePetResponse{
		ID:               p.ID.String(),
		PetName:          p.PetName,
		Species:          p.Species,
		PetAge:           p.PetAge,
		Color:            p.Color,
		PetGender:        p.PetGender,
		Breed:            p.Breed,
		Pedigree:         p.Pedigree,
		GoodWithChildren: p.GoodWithChildren,
		GoodWithAnimals:  p.GoodWithAnimals,
		PetDescription:   p.PetDescription,
		IsActive:         p.IsActive,
		Price:            p.Price,
	}
}

func PetToDto(p models.Pet) PetResponse {
	return PetResponse{
		ID:               p.ID.String(),
		VetPassportID:    p.VetPassportID.String(),
		SellerID:    p.SellerID.String(),
		PetName:          p.PetName,
		Species:          p.Species,
		PetAge:           p.PetAge,
		Color:            p.Color,
		PetGender:        p.PetGender,
		Breed:            p.Breed,
		Pedigree:         p.Pedigree,
		GoodWithChildren: p.GoodWithChildren,
		GoodWithAnimals:  p.GoodWithAnimals,
		PetDescription:   p.PetDescription,
		IsActive:         p.IsActive,
		Price:            p.Price,
	}
}

func PetsToDto(pets []models.Pet) []PetResponse {
	res := make([]PetResponse, len(pets))
	for i, p := range pets {
		res[i] = PetToDto(p)
	}
	return res
}

func UpdatePetRequestToModel(req UpdatePetRequest) models.Pet {
	return models.Pet{
		PetName:          req.PetName,
		Species:          req.Species,
		PetAge:           req.PetAge,
		Color:            req.Color,
		PetGender:        req.PetGender,
		Breed:            req.Breed,
		Pedigree:         req.Pedigree,
		GoodWithChildren: req.GoodWithChildren,
		GoodWithAnimals:  req.GoodWithAnimals,
		PetDescription:   req.PetDescription,
		IsActive:         req.IsActive,
		Price:            req.Price,
	}
}
