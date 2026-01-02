package models

import "time"

type PersonalYearlyGoal struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	Year          int       `json:"year"`
	DistanceGoal  float64   `json:"distance_goal"`  // km
	ElevationGoal float64   `json:"elevation_goal"` // meters
	SummitGoal    int       `json:"summit_goal"`    // count
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
