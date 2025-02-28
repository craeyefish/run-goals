package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func handleStravaWebhookVerification(w http.ResponseWriter, r *http.Request) {
	// Strava makes a GET request with `hub.challenge` to verify
	challenge := r.URL.Query().Get("hub.challenge")
	if challenge == "" {
		http.Error(w, "missing hub.challenge", http.StatusBadRequest)
		return
	}
	type responseBody struct {
		HubChallenge string `json:"hub.challenge"`
	}
	json.NewEncoder(w).Encode(responseBody{HubChallenge: challenge})
}

type StravaWebhookPayload struct {
	AspectType     string                 `json:"aspect_type"`
	EventTime      int64                  `json:"event_time"`
	ObjectType     string                 `json:"object_type"`
	ObjectID       int64                  `json:"object_id"`
	OwnerID        int64                  `json:"owner_id"`
	SubscriptionID int64                  `json:"subscription_id"`
	Updates        map[string]interface{} `json:"updates"`
}

func handleStravaWebhookEvents(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// This is just the verification challenge
		handleStravaWebhookVerification(w, r)
		return
	case http.MethodPost:
		// Parse the incoming JSON
		var payload StravaWebhookPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "failed to parse webhook payload", http.StatusBadRequest)
			return
		}

		log.Printf("Received Strava webhook event: %+v\n", payload)

		// We only care about activities; ignore other object_types if needed
		if payload.ObjectType == "activity" {
			// Find the user in DB
			var user User
			// 'strava_athlete_id' is how we store the Strava athlete ID in the User model
			if err := DB.Where("strava_athlete_id = ?", payload.OwnerID).First(&user).Error; err != nil {
				// If no user found or DB error, just log and move on
				log.Printf("No matching user or error: %v\n", err)
			} else {
				// We found the user - now fetch updated stats
				dist, err := fetchUserDistance(user)
				if err != nil {
					log.Printf("Error fetching updated distance for user %d: %v\n", user.ID, err)
				} else {
					// Update and save
					user.LastDistance = dist
					user.LastUpdated = time.Now()
					if err := DB.Save(&user).Error; err != nil {
						log.Printf("Error saving user distance: %v\n", err)
					} else {
						log.Printf("Updated user %d with new distance: %.2f km\n", user.ID, dist)
					}
				}
			}
		}

		// Respond 200 so Strava knows we received it
		w.WriteHeader(http.StatusOK)
	default:
		// Strava only does GET (for verify) or POST (for events).
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
