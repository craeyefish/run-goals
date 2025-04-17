package handlers

import (
	"log"
	"net/http"
	"run-goals/controllers"
)

type AuthHandler struct {
	l                *log.Logger
	authController   *controllers.AuthController
	stravaController *controllers.StravaController
}

func NewAuthHandler(
	l *log.Logger,
	authController *controllers.AuthController,
	stravaController *controllers.StravaController,
) *AuthHandler {
	return &AuthHandler{
		l,
		authController,
		stravaController,
	}
}

// ServeHTTP is the main entry point for the handler and satisfies the handler interface
func (handler *AuthHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// handle request to get activities
	switch r.URL.Path {
	case "/auth/refresh":
		handler.authController.RefreshToken(rw, r)
		return
	case "/auth/strava/callback":
		handler.stravaController.ProcessCallback(rw, r)
		return
	}
}
