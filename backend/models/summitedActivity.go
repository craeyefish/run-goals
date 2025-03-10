package models

import "time"

type SummitedActivity struct {
	UserName   string    `json:"user_name"`
	UserID     int64     `json:"user_id"`
	ActivityID int64     `json:"activity_id"`
	SummitedAt time.Time `json:"summited_at"`
}
