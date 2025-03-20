package models

import "time"

type GroupGoal struct {
	ID          int64     `json:"id"`
	GroupId     int64     `json:"group_id"`
	Name        string    `json:"name"`
	TargetValue string    `json:"target_value"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	CreatedAt   time.Time `json:"created_at"`
}
