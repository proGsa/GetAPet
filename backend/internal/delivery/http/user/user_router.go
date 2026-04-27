package user

import (
	"getapet-backend/internal/delivery/middleware"
	"getapet-backend/internal/models"

	"github.com/gorilla/mux"
)

type UserRouter struct {
	UserUsecase models.UserService
	JWTSecret   string
}

func NewUserRouter(us models.UserService, jwtSecret string) *UserRouter {
	return &UserRouter{UserUsecase: us, JWTSecret: jwtSecret}
}

func (ur *UserRouter) SetupRoutes(router *mux.Router) {
	userRouter := router.PathPrefix("/users").Subrouter()

	userRouter.HandleFunc("", ur.CreateUser).Methods("POST", "OPTIONS")
	userRouter.HandleFunc("/login", ur.Login).Methods("POST", "OPTIONS")

	protected := userRouter.NewRoute().Subrouter()
	protected.Use(middleware.JWTMiddleware(ur.JWTSecret))
	// protected.HandleFunc("/logout", ur.Logout).Methods("POST", "OPTIONS")
	protected.HandleFunc("", ur.GetUsers).Methods("GET")
	protected.HandleFunc("/{id}", ur.GetUser).Methods("GET")
	protected.HandleFunc("/{id}", ur.UpdateUser).Methods("PUT", "OPTIONS")
	protected.HandleFunc("/{id}", ur.DeleteUser).Methods("DELETE", "OPTIONS")
}
