package pet

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func writeSuccessResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeErrorResponse(w http.ResponseWriter, status int, err error, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error":   err.Error(),
		"message": message,
	})
}

func parseIDFromPath(r *http.Request) (int, error) {
	return strconv.Atoi(mux.Vars(r)["id"])
}

func writeServiceUnavailable(w http.ResponseWriter) {
	writeErrorResponse(w, http.StatusServiceUnavailable, errors.New("service unavailable"), "Сервис временно недоступен")
}
