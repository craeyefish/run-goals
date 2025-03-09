package models

import "time"

type UserPeak struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	PeakID     int64     `json:"peak_id"`
	ActivityID int64     `json:"activity_id"` // the activity that triggered the "bag"
	SummitedAt time.Time `json:"summited_at"` // when we detected the visit
	// optional: distance threshold or actual min distance for reference
}
