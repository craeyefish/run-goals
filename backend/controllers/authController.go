package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"run-goals/services"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type AuthControllerInterface interface {
	RefreshToken(rw http.ResponseWriter, r *http.Request)
}

type AuthController struct {
	l          *log.Logger
	jwtService *services.JWTService
}

func NewAuthController(
	l *log.Logger,
	jwtService *services.JWTService,
) *AuthController {
	return &AuthController{
		l:          l,
		jwtService: jwtService,
	}
}

// POST /auth/refresh
func (h *AuthController) RefreshToken(rw http.ResponseWriter, r *http.Request) {
	refreshToken := r.Header.Get("Authorization")
	refreshToken = strings.TrimPrefix(refreshToken, "Bearer ")

	token, err := h.jwtService.ValidateToken(refreshToken)
	if err != nil {
		http.Error(rw, "invalid refresh token", http.StatusBadRequest)
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := int64(claims["sub"].(float64))

	newAccessToken, _ := h.jwtService.GenerateAccessToken(userID)

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(map[string]string{"accessToken": newAccessToken})
}
