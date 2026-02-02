package middleware

import (
	"context"
	"net/http"
	"strings"

	"vietio/internal/auth"
)

type contextKey string

const UserIdKey contextKey = "user_id"

func AuthJWT(authService *auth.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			token := parts[1]

			claims, err := authService.ParseAndValidateJWT(token)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIdKey, claims.UserId)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
