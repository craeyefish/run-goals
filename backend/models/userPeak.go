package models

import "time"

type UserPeak struct {
	ID         uint      `json: "id"`
	UserID     uint      `json:"user_id"`
	PeakID     uint      `json:"peak_id"`
	ActivityID uint      `json:"activity_id"` // the activity that triggered the "bag"
	SummitedAt time.Time `json:"summited_at"` // when we detected the visit
	// optional: distance threshold or actual min distance for reference
}
