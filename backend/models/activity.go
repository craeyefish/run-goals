package models

import (
	"io"
	"time"
	"encoding/json"
)

type Activity struct {
	ID               uint      `json:"id"`
	StravaActivityID int64     `json:"strava_activity_id"`
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

// ToJSON serializes the given interface into a string based JSON format
func (activity *Activity) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(activity)
}

// FromJSON deserializes the object from JSON string, in an io.Reader, to the given interface
func (activity *Activity) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(activity)
}
