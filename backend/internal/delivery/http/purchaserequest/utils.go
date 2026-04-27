package purchaserequest

import (
	"encoding/json"
	"errors"
	"net/http"

	// "getapet-backend/internal/delivery/middleware"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func writeSuccessResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeErrorResponse(w http.ResponseWriter, status int, err error, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error(), Message: message})
}

func parseIDFromPath(r *http.Request) (uuid.UUID, error) {
	return parseIDFromPathParam(r, "id")
}

func parseIDFromPathParam(r *http.Request, name string) (uuid.UUID, error) {
	return uuid.Parse(mux.Vars(r)[name])
}

// не используется - удалить ?
// func userIDFromContext(r *http.Request) (uuid.UUID, error) {
// 	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
// 	if !ok || userID == "" {
// 		return uuid.Nil, errors.New("user_id is missing in context")
// 	}
// 	return uuid.Parse(userID)
// }

func writeServiceUnavailable(w http.ResponseWriter) {
	writeErrorResponse(w, http.StatusServiceUnavailable, errors.New("service unavailable"), "Service is temporarily unavailable")
}
