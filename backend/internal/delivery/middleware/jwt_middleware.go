package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type userIDCtxKey string

const UserIDKey userIDCtxKey = "userID"

type authErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func JWTMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if jwtSecret == "" {
				writeUnauthorized(w, errors.New("jwt secret is empty"), "Сервис авторизации не настроен")
				return
			}

			authHeader := r.Header.Get("Authorization")
			tokenString, err := ExtractRawToken(authHeader)
			if err != nil {
				writeUnauthorized(w, err, err.Error())
				return
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return []byte(jwtSecret), nil
			})
			if err != nil || !token.Valid {
				writeUnauthorized(w, errors.New("invalid token"), "Невалидный JWT токен")
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				writeUnauthorized(w, errors.New("invalid token claims"), "Некорректные claims токена")
				return
			}

			userID, ok := claims["user_id"].(string)
			if !ok || userID == "" {
				writeUnauthorized(w, errors.New("user_id claim is missing"), "Токен не содержит user_id")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ExtractRawToken(authHeader string) (string, error) {
	authHeader = strings.TrimSpace(authHeader)
	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}
	if len(authHeader) > 7 && strings.EqualFold(authHeader[:7], "Bearer ") {
		authHeader = strings.TrimSpace(authHeader[7:])
	}
	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}
	return authHeader, nil
}

func writeUnauthorized(w http.ResponseWriter, err error, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(authErrorResponse{
		Error:   err.Error(),
		Message: message,
	})
}
