package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// NullableString is a wrapper around sql.NullString that marshals to/from JSON as a regular string
type NullableString struct {
	sql.NullString
}

// MarshalJSON converts NullableString to JSON string (empty string if null)
func (ns NullableString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return json.Marshal("")
}

// UnmarshalJSON converts JSON string to NullableString
func (ns *NullableString) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == "" {
		ns.Valid = false
		ns.String = ""
	} else {
		ns.Valid = true
		ns.String = s
	}
	return nil
}

type User struct {
	ID              int64          `json:"id"`
	StravaAthleteID int64          `json:"strava_athelete_id"` // Strava athlete ID, unique
	Username        NullableString `json:"username"`            // User-chosen display name
	IsAdmin         bool           `json:"is_admin"`            // Whether user has admin privileges
	AccessToken     string         `json:"access_token"`
	RefreshToken    string         `json:"refresh_token"`
	ExpiresAt       time.Time      `json:"expires_at"`
	LastDistance    float64        `json:"last_distance"` // Cached distance in km
	LastUpdated     time.Time      `json:"last_updated"`  // When we last fetched from Strava
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}
