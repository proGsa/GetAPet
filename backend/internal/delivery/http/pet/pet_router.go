package pet

import (
	"getapet-backend/internal/delivery/middleware"
	"getapet-backend/internal/models"

	"github.com/gorilla/mux"
)

type PetRouter struct {
	PetUsecase models.PetService
	JWTSecret  string
}

func NewPetRouter(ps models.PetService, jwtSecret string) *PetRouter {
	return &PetRouter{PetUsecase: ps, JWTSecret: jwtSecret}
}

func (pr *PetRouter) SetupRoutes(router *mux.Router) {
	petRouter := router.PathPrefix("/pets").Subrouter()

	petRouter.HandleFunc("", pr.GetPets).Methods("GET")
	petRouter.HandleFunc("/{id}", pr.GetPet).Methods("GET")

	protected := petRouter.NewRoute().Subrouter()
	protected.Use(middleware.JWTMiddleware(pr.JWTSecret))
	protected.HandleFunc("", pr.CreatePet).Methods("POST", "OPTIONS")
	protected.HandleFunc("/{id}", pr.UpdatePet).Methods("PUT", "OPTIONS")
	protected.HandleFunc("/{id}", pr.DeletePet).Methods("DELETE", "OPTIONS")
}
