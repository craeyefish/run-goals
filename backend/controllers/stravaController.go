package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"run-goals/models"
	"run-goals/services"
)

type StravaController struct {
	l             *log.Logger
	stravaService *services.StravaService
}

func NewStravaController(
	l *log.Logger,
	stravaService *services.StravaService,
) *StravaController {
	return &StravaController{
		l:             l,
		stravaService: stravaService,
	}
}

func (c *StravaController) VerifyWebhookEvent(rw http.ResponseWriter, r *http.Request) {
	// Strava makes a GET request with `hub.challenge` to verify
	challenge := r.URL.Query().Get("hub.challenge")
	if challenge == "" {
		http.Error(rw, "missing hub.challenge", http.StatusBadRequest)
		return
	}
	type responseBody struct {
		HubChallenge string `json:"hub.challenge"`
	}
	json.NewEncoder(rw).Encode(responseBody{HubChallenge: challenge})
}

func (c *StravaController) ProcessWebhookEvent(rw http.ResponseWriter, r *http.Request) {
	var payload models.StravaWebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(rw, "failed to parse webhook payload", http.StatusBadRequest)
		return
	}
	c.l.Printf("Received Strava webhook event: %+v\n", payload)

	// We only care about activities; ignore other object_types if needed
	if payload.ObjectType == "activity" {
		c.stravaService.ProcessWebhookEvent(payload)
	}
	rw.WriteHeader(http.StatusOK)
}

func (c *StravaController) ProcessCallback(rw http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(rw, "missing code", http.StatusBadRequest)
		return
	}

	err := c.stravaService.ProcessCallback(code)
	if err != nil {
		http.Error(rw, "Failed to process callback", http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, r, "<add domain here>?joined=1", http.StatusFound)
}
