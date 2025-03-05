package models

import (
	"time"
)

type Activity struct {
	ID               int64     `json:"id"`
	StravaActivityID int64     `json:"strava_activity_id"`
	UserID           int64     `json:"user_id"`
	Name             string    `json:"name"`
	Distance         float64   `json:"distance"`
	StartDate        time.Time `json:"start_date"`
	MapPolyline      string    `json:"map_polyline"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Non Strava
	HasSummit bool `json:"has_summit"`
}
