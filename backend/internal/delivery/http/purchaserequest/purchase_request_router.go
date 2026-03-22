package purchaserequest

import (
	"getapet-backend/internal/models"

	"github.com/gorilla/mux"
)

type PurchaseRequestRouter struct {
	PurchaseRequestUsecase models.PurchaseRequestService
}

func NewPurchaseRequestRouter(prs models.PurchaseRequestService) *PurchaseRequestRouter {
	return &PurchaseRequestRouter{PurchaseRequestUsecase: prs}
}

func (pr *PurchaseRequestRouter) SetupRoutes(router *mux.Router) {
	purchaseRequestRouter := router.PathPrefix("/purchase-requests").Subrouter()

	purchaseRequestRouter.HandleFunc("", pr.CreatePurchaseRequest).Methods("POST", "OPTIONS")
	purchaseRequestRouter.HandleFunc("", pr.GetPurchaseRequests).Methods("GET")
	purchaseRequestRouter.HandleFunc("/{id}", pr.GetPurchaseRequest).Methods("GET")
	purchaseRequestRouter.HandleFunc("/{id}/status", pr.UpdatePurchaseRequestStatus).Methods("PATCH", "OPTIONS")
	purchaseRequestRouter.HandleFunc("/{id}", pr.DeletePurchaseRequest).Methods("DELETE", "OPTIONS")
}
