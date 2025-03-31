package middleware

import (
	"context"
	"net/http"
	"run-goals/meta"
	"run-goals/services"
	"strings"

	"github.com/golang-jwt/jwt"
)

func JWT(jwtService *services.JWTService, next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(rw, "Unauthorized", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(rw, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenStr := parts[1]
		token, err := jwtService.ValidateToken(tokenStr)
		if err != nil || !token.Valid {
			http.Error(rw, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(rw, "Unauthorized claims", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["sub"].(float64) // watch out for type
		if !ok {
			http.Error(rw, "Unauthorized user", http.StatusUnauthorized)
			return
		}

		// Attach userID to context
		ctx := context.WithValue(r.Context(), meta.ContextKeyUserID, int64(userID))
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
