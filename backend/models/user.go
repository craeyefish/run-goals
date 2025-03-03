package models

import (
	"time"
)

type User struct {
	StravaAthleteID int64 `json:"strava_athelete_id"` // Strava athlete ID, unique
	AccessToken     string `json:"access_token"`
	RefreshToken    string `json:"refresh_token"`
	ExpiresAt       int64 `json:"expires_at"`     // Unix timestamp
	LastDistance    float64 `json:"last_distance"`   // Cached distance in km
	LastUpdated     time.Time `json:"last_update"` // When we last fetched from Strava

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
