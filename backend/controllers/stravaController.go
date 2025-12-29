package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"run-goals/daos"
	"run-goals/models"
	"run-goals/services"
)

type StravaController struct {
	l             *log.Logger
	jwtService    *services.JWTService
	stravaService *services.StravaService
	summitService *services.SummitService
	activityDao   *daos.ActivityDao
}

func NewStravaController(
	l *log.Logger,
	jwtService *services.JWTService,
	stravaService *services.StravaService,
	summitService *services.SummitService,
	activityDao *daos.ActivityDao,
) *StravaController {
	return &StravaController{
		l:             l,
		jwtService:    jwtService,
		stravaService: stravaService,
		summitService: summitService,
		activityDao:   activityDao,
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

		// Trigger summit calculation asynchronously for this activity
		go func() {
			activity, err := c.activityDao.GetActivityByStravaID(payload.ObjectID)
			if err != nil {
				c.l.Printf("Error fetching activity for summit calc: %v", err)
				return
			}
			if activity == nil {
				c.l.Printf("Activity %d not found for summit calc", payload.ObjectID)
				return
			}
			if activity.SummitsCalculated {
				c.l.Printf("Summits already calculated for activity %d", activity.ID)
				return
			}
			c.l.Printf("Calculating summits for activity %d from webhook", activity.ID)
			if err := c.summitService.CalculateSummitsForActivity(activity); err != nil {
				c.l.Printf("Error calculating summits: %v", err)
			}
		}()
	}
	rw.WriteHeader(http.StatusOK)
}

func (c *StravaController) ProcessCallback(rw http.ResponseWriter, r *http.Request) {
	var payload struct {
		Code string `json:"code"`
	}
	c.l.Printf("Received Strava callback payload: %+v\n", payload)
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(rw, "failed to parse callback payload", http.StatusBadRequest)
		return
	}

	c.l.Printf("processing strava callback")
	user, err := c.stravaService.ProcessCallback(payload.Code)
	if err != nil {
		http.Error(rw, "Failed to process callback", http.StatusInternalServerError)
		return
	}

	accessTokenString, err := c.jwtService.GenerateAccessToken(user.ID)
	if err != nil {
		http.Error(rw, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	refreshTokenString, err := c.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		http.Error(rw, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(map[string]string{"accessToken": accessTokenString, "refreshToken": refreshTokenString})
}
