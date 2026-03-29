package user

import (
	"getapet-backend/internal/models"

	"github.com/gorilla/mux"
)

type UserRouter struct {
	UserUsecase models.UserService
}

func NewUserRouter(us models.UserService) *UserRouter {
	return &UserRouter{UserUsecase: us}
}

func (ur *UserRouter) SetupRoutes(router *mux.Router) {
	userRouter := router.PathPrefix("/users").Subrouter()

	userRouter.HandleFunc("", ur.CreateUser).Methods("POST", "OPTIONS")
	userRouter.HandleFunc("/login", ur.Login).Methods("POST", "OPTIONS")
	userRouter.HandleFunc("", ur.GetUsers).Methods("GET")
	userRouter.HandleFunc("/{id}", ur.GetUser).Methods("GET")
	userRouter.HandleFunc("/{id}", ur.UpdateUser).Methods("PUT", "OPTIONS")
	userRouter.HandleFunc("/{id}", ur.DeleteUser).Methods("DELETE", "OPTIONS")
}
