package vetpassport

import (
	"getapet-backend/internal/models"

	"github.com/gorilla/mux"
)

type VetPassportRouter struct {
	VetPassportUsecase models.VetPassportService
}

func NewVetPassportRouter(vps models.VetPassportService) *VetPassportRouter {
	return &VetPassportRouter{VetPassportUsecase: vps}
}

func (vr *VetPassportRouter) SetupRoutes(router *mux.Router) {
	vetPassportRouter := router.PathPrefix("/vet-passports").Subrouter()

	vetPassportRouter.HandleFunc("", vr.CreateVetPassport).Methods("POST", "OPTIONS")
	vetPassportRouter.HandleFunc("", vr.GetVetPassports).Methods("GET")
	vetPassportRouter.HandleFunc("/{id}", vr.GetVetPassport).Methods("GET")
	vetPassportRouter.HandleFunc("/{id}", vr.UpdateVetPassport).Methods("PUT", "OPTIONS")
	vetPassportRouter.HandleFunc("/{id}", vr.DeleteVetPassport).Methods("DELETE", "OPTIONS")
}
