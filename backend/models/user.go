package models

import (
	"time"
)

type User struct {
	ID              int64     `json:"id"`
	StravaAthleteID int64     `json:"strava_athelete_id"` // Strava athlete ID, unique
	AccessToken     string    `json:"access_token"`
	RefreshToken    string    `json:"refresh_token"`
	ExpiresAt       time.Time `json:"expires_at"`
	LastDistance    float64   `json:"last_distance"` // Cached distance in km
	LastUpdated     time.Time `json:"last_updated"`  // When we last fetched from Strava
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
