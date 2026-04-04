package httpserver

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "getapet-backend/docs"
	"getapet-backend/internal/delivery/http/pet"
	"getapet-backend/internal/delivery/http/purchaserequest"
	"getapet-backend/internal/delivery/http/user"
	"getapet-backend/internal/delivery/http/vetpassport"
	"getapet-backend/internal/repository"
	"getapet-backend/internal/usecase"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type HTTPServer struct {
	server *http.Server
}

type HealthResponse struct {
	Status string `json:"status"`
}

func NewHTTPServer(addr string) *HTTPServer {
	return &HTTPServer{
		server: &http.Server{
			Addr: addr,
		},
	}
}

func (s *HTTPServer) Start(db *sql.DB) error {
	s.server.Handler = setupRoutes(db)
	log.Println("Server is running on", s.server.Addr)
	return s.server.ListenAndServe()
}

func setupRoutes(db *sql.DB) *mux.Router {
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)

	petRepo := repository.NewPetRepository(db)
	petUsecase := usecase.NewPetUsecase(petRepo)

	vetPassportRepo := repository.NewVetPassportRepository(db)
	vetPassportUsecase := usecase.NewVetPassportUsecase(vetPassportRepo)

	purchaseRequestRepo := repository.NewPurchaseRequestRepository(db)
	purchaseRequestUsecase := usecase.NewPurchaseRequestUsecase(purchaseRequestRepo)

	router := mux.NewRouter()
	router.HandleFunc("/health", healthHandler).Methods(http.MethodGet)
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	api := router.PathPrefix("/api").Subrouter()

	userRouter := user.NewUserRouter(userUsecase, os.Getenv("JWT_SECRET"))
	userRouter.SetupRoutes(api)

	petRouter := pet.NewPetRouter(petUsecase, os.Getenv("JWT_SECRET"))
	petRouter.SetupRoutes(api)

	vetPassportRouter := vetpassport.NewVetPassportRouter(vetPassportUsecase)
	vetPassportRouter.SetupRoutes(api)

	purchaseRequestRouter := purchaserequest.NewPurchaseRequestRouter(purchaseRequestUsecase)
	purchaseRequestRouter.SetupRoutes(api)

	return router
}

// healthHandler godoc
// @Summary Check service health
// @Tags system
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(HealthResponse{Status: "ok"})
}
