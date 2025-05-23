package models

import (
	"strings"
	"time"
)

type Activity struct {
	ID               int64     `json:"id"`
	StravaActivityId int64     `json:"strava_activity_id"`
	StravaAthleteId  int64     `json:"strava_athlete_id"`
	UserID           int64     `json:"user_id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Distance         float64   `json:"distance"`
	StartDate        time.Time `json:"start_date"`
	MapPolyline      string    `json:"map_polyline"`
	PhotoURL         string    `json:"photo_url"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Non Strava
	HasSummit bool `json:"has_summit"`
}

func (a *Activity) IsHG() bool {
	return strings.Contains(strings.ToLower(a.Name), "#hg")
}
