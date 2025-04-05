package services

import (
	"log"
	"run-goals/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTServiceInterface interface {
	GenerateToken(userID int64) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
}

type JWTService struct {
	l         *log.Logger
	secretKey []byte
}

func NewJWTService(
	l *log.Logger,
	config *config.Config,
) *JWTService {
	return &JWTService{
		l:         l,
		secretKey: []byte(config.JWT.Secret),
	}
}

func (j *JWTService) GenerateToken(userID int64) (string, error) {
	// Create the claims
	claims := jwt.MapClaims{
		"sub": float64(userID),
		"exp": time.Now().Add(time.Hour * 24).Unix(), // 24 hour expiry
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

func (j *JWTService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return j.secretKey, nil
	})
}
