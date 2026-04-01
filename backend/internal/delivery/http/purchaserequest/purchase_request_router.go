package purchaserequest

import (
	"getapet-backend/internal/delivery/middleware"
	"getapet-backend/internal/models"

	"github.com/gorilla/mux"
)

type PurchaseRequestRouter struct {
	PurchaseRequestUsecase models.PurchaseRequestService
	JWTSecret              string
}

func NewPurchaseRequestRouter(prs models.PurchaseRequestService, jwtSecret string) *PurchaseRequestRouter {
	return &PurchaseRequestRouter{PurchaseRequestUsecase: prs, JWTSecret: jwtSecret}
}

func (pr *PurchaseRequestRouter) SetupRoutes(router *mux.Router) {
	purchaseRequestRouter := router.PathPrefix("/purchase-requests").Subrouter()

	protected := purchaseRequestRouter.NewRoute().Subrouter()
	protected.Use(middleware.JWTMiddleware(pr.JWTSecret))
	protected.HandleFunc("", pr.CreatePurchaseRequest).Methods("POST", "OPTIONS")
	protected.HandleFunc("", pr.GetPurchaseRequests).Methods("GET")
	protected.HandleFunc("/buyer/{id}", pr.GetPurchaseRequestsByBuyer).Methods("GET")
	protected.HandleFunc("/seller/{id}", pr.GetPurchaseRequestsBySeller).Methods("GET")
	protected.HandleFunc("/pet/{pet_id}", pr.GetPurchaseRequestsByPet).Methods("GET")
	protected.HandleFunc("/{id}", pr.GetPurchaseRequest).Methods("GET")
	protected.HandleFunc("/{id}/status", pr.UpdatePurchaseRequestStatus).Methods("PATCH", "OPTIONS")
	protected.HandleFunc("/{id}", pr.DeletePurchaseRequest).Methods("DELETE", "OPTIONS")
}
