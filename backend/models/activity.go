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
	Type             string    `json:"type"`       // Run, Hike, Walk, etc.
	SportType        string    `json:"sport_type"` // More specific: Trail Run, etc.
	Description      string    `json:"description"`
	Distance         float64   `json:"distance"`
	Elevation        float64   `json:"total_elevation_gain"`
	MovingTime       float64   `json:"moving_time"`
	StartDate        time.Time `json:"start_date"`
	MapPolyline      string    `json:"map_polyline"`
	PhotoURL         string    `json:"photo_url"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Non Strava
	HasSummit          bool `json:"has_summit"`
	SummitsCalculated  bool `json:"summits_calculated"` // Whether summit detection has been run for this activity
}

func (a *Activity) IsHG() bool {
	return strings.Contains(strings.ToLower(a.Name), "#hg")
}

// ActivityWithUser includes user information for display
type ActivityWithUser struct {
	Activity
	UserName        string  `json:"userName"`
	StravaAthleteID int64   `json:"stravaAthleteId"`
	PeakNames       *string `json:"peakNames,omitempty"` // Comma-separated peak names for summit activities
}
