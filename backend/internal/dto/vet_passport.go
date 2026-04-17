package dto

import (
	"getapet-backend/internal/models"
	// "github.com/google/uuid"
)

type CreateVetPassportRequest struct {
	Chipping           bool   `json:"chipping" example:"true"`
	Sterilization      bool   `json:"sterilization" example:"false"`
	HealthIssues       string `json:"health_issues" validate:"omitempty,max=1000" example:"Аллергия на корм"`
	Vaccinations       string `json:"vaccinations" validate:"omitempty,max=1000" example:"Привит от бешенства"`
	ParasiteTreatments string `json:"parasite_treatments" validate:"omitempty,max=1000" example:"Обработка от блох"`
}

type UpdateVetPassportRequest struct {
	Chipping           bool   `json:"chipping" example:"true"`
	Sterilization      bool   `json:"sterilization" example:"true"`
	HealthIssues       string `json:"health_issues" validate:"omitempty,max=1000" example:"Нет"`
	Vaccinations       string `json:"vaccinations" validate:"omitempty,max=1000" example:"Все прививки сделаны"`
	ParasiteTreatments string `json:"parasite_treatments" validate:"omitempty,max=1000" example:"Регулярная обработка"`
}

type VetPassportResponse struct {
	ID                 string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Chipping           bool   `json:"chipping" example:"true"`
	Sterilization      bool   `json:"sterilization" example:"false"`
	HealthIssues       string `json:"health_issues" example:"Нет"`
	Vaccinations       string `json:"vaccinations" example:"Привит"`
	ParasiteTreatments string `json:"parasite_treatments" example:"Обработан"`
}

func CreateVetPassportRequestFromDTO(req CreateVetPassportRequest) models.VetPassport {
	return models.VetPassport{
		Chipping:           req.Chipping,
		Sterilization:      req.Sterilization,
		HealthIssues:       req.HealthIssues,
		Vaccinations:       req.Vaccinations,
		ParasiteTreatments: req.ParasiteTreatments,
	}
}

func UpdateVetPassportRequestFromDTO(req UpdateVetPassportRequest) models.VetPassport {
	return models.VetPassport{
		Chipping:           req.Chipping,
		Sterilization:      req.Sterilization,
		HealthIssues:       req.HealthIssues,
		Vaccinations:       req.Vaccinations,
		ParasiteTreatments: req.ParasiteTreatments,
	}
}

func VetPassportToDTO(v models.VetPassport) VetPassportResponse {
	return VetPassportResponse{
		ID:                 v.ID.String(),
		Chipping:           v.Chipping,
		Sterilization:      v.Sterilization,
		HealthIssues:       v.HealthIssues,
		Vaccinations:       v.Vaccinations,
		ParasiteTreatments: v.ParasiteTreatments,
	}
}

func VetPassportsToDTO(passports []models.VetPassport) []VetPassportResponse {
	res := make([]VetPassportResponse, len(passports))
	for i, p := range passports {
		res[i] = VetPassportToDTO(p)
	}
	return res
}
