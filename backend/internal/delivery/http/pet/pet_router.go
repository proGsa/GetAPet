package pet

import (
	"getapet-backend/internal/models"

	"github.com/gorilla/mux"
)

type PetRouter struct {
	PetUsecase models.PetService
}

func NewPetRouter(ps models.PetService) *PetRouter {
	return &PetRouter{PetUsecase: ps}
}

func (pr *PetRouter) SetupRoutes(router *mux.Router) {
	petRouter := router.PathPrefix("/pets").Subrouter()

	petRouter.HandleFunc("", pr.CreatePet).Methods("POST", "OPTIONS")
	petRouter.HandleFunc("", pr.GetPets).Methods("GET")
	petRouter.HandleFunc("/{id}", pr.GetPet).Methods("GET")
	petRouter.HandleFunc("/{id}", pr.UpdatePet).Methods("PUT", "OPTIONS")
	petRouter.HandleFunc("/{id}", pr.DeletePet).Methods("DELETE", "OPTIONS")
}
