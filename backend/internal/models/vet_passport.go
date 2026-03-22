package models

import "errors"

type VetPassport struct {
	ID                 int    `json:"id" db:"id"`
	Chipping           bool   `json:"chipping" db:"chipping"`
	Sterilization      bool   `json:"sterilization" db:"sterilization"`
	HealthIssues       string `json:"health_issues" db:"health_issues"`
	Vaccinations       string `json:"vaccinations" db:"vaccinations"`
	ParasiteTreatments string `json:"parasite_treatments" db:"parasite_treatments"`
}

type VetPassportRepository interface {
	Create(passport *VetPassport) (*VetPassport, error)
	GetAll() ([]VetPassport, error)
	GetByID(id int) (*VetPassport, error)
	Update(id int, passport *VetPassport) (*VetPassport, error)
	Delete(id int) error
}

type VetPassportService interface {
	Create(passport *VetPassport) (*VetPassport, error)
	GetAll() ([]VetPassport, error)
	GetByID(id int) (*VetPassport, error)
	Update(id int, passport *VetPassport) (*VetPassport, error)
	Delete(id int) error
}

var (
	ErrVetPassportNotFound = errors.New("vet passport not found")
)
