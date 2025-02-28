package main

import "time"

// Represents a single user's aggregated data
type UserContribution struct {
	ID            int64   `json:"id"`
	TotalDistance float64 `json:"totalDistance"` // in kilometers
}

// Response for the overall goal
type GoalProgress struct {
	Goal            float64            `json:"goal"` // total km goal
	CurrentProgress float64            `json:"currentProgress"`
	Contributions   []UserContribution `json:"contributions"`
}

// User represents a person who has joined your challenge.
type User struct {
	ID              uint  `gorm:"primaryKey"`
	StravaAthleteID int64 `gorm:"uniqueIndex"` // Strava athlete ID, unique
	AccessToken     string
	RefreshToken    string
	ExpiresAt       int64     // Unix timestamp
	LastDistance    float64   // Cached distance in km
	LastUpdated     time.Time // When we last fetched from Strava

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Activity struct {
	ID               uint      `gorm:"primaryKey"     json:"id"`
	StravaActivityID int64     `gorm:"uniqueIndex"    json:"strava_activity_id"`
	UserID           uint      `json:"user_id"`
	Name             string    `json:"name"`
	Distance         float64   `json:"distance"`
	StartDate        time.Time `json:"start_date"`
	MapPolyline      string    `json:"map_polyline"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Non Strava
	HasSummit bool `json:"has_summit"`
}

type Peak struct {
	ID    uint    `gorm:"primaryKey" json:"id"`
	OsmID int64   `gorm:"uniqueIndex" json:"osm_id"` // OSM Node ID
	Lat   float64 `json:"lat"`
	Lon   float64 `json:"lon"`
	Name  string  `json:"name"`
	ElevM float64 `json:"elev_m"` // Elevation in meters (parse from "ele" tag if present)
}

type UserPeak struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     uint      `json:"user_id"`
	PeakID     uint      `json:"peak_id"`
	ActivityID uint      `json:"activity_id"` // the activity that triggered the "bag"
	SummitedAt time.Time `json:"summited_at"` // when we detected the visit
	// optional: distance threshold or actual min distance for reference
}
