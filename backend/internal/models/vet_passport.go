package models

import ("errors"
"github.com/google/uuid"
)

type VetPassport struct {
	ID                 uuid.UUID    `json:"id" db:"id"`
	Chipping           bool   `json:"chipping" db:"chipping"`
	Sterilization      bool   `json:"sterilization" db:"sterilization"`
	HealthIssues       string `json:"health_issues" db:"health_issues"`
	Vaccinations       string `json:"vaccinations" db:"vaccinations"`
	ParasiteTreatments string `json:"parasite_treatments" db:"parasite_treatments"`
}

type VetPassportRepository interface {
	Create(passport *VetPassport) (*VetPassport, error)
	GetAll() ([]VetPassport, error)
	GetByID(id uuid.UUID) (*VetPassport, error)
	Update(id uuid.UUID, passport *VetPassport) (*VetPassport, error)
	Delete(id uuid.UUID) error
}

type VetPassportService interface {
	Create(passport *VetPassport) (*VetPassport, error)
	GetAll() ([]VetPassport, error)
	GetByID(id uuid.UUID) (*VetPassport, error)
	Update(id uuid.UUID, passport *VetPassport) (*VetPassport, error)
	Delete(id uuid.UUID) error
}

var (
	ErrVetPassportNotFound = errors.New("vet passport not found")
)
