package handlers

import (
	"log"
	"net/http"
	"run-goals/controllers"
)

type StravaHandler struct {
	l                *log.Logger
	stravaController *controllers.StravaController
}

func NewStravaHandler(
	l *log.Logger,
	stravaController *controllers.StravaController,
) *StravaHandler {
	return &StravaHandler{
		l,
		stravaController,
	}
}

// ServeHTTP is the main entry point for the handler and satisfies the handler interface
func (handler *StravaHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// handle request to get activities
	switch r.URL.Path {
	case "/webhook/strava":
		switch r.Method {
		case http.MethodGet:
			// This is just the verification challenge
			handler.stravaController.VerifyWebhookEvent(rw, r)
			return
		case http.MethodPost:
			// Process the webhook event
			handler.stravaController.ProcessWebhookEvent(rw, r)
			return
		default:
			http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
