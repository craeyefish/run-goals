package dto

import "time"

type CreateGroupGoalRequest struct {
	GroupID     int64     `json:"group_id"`
	Name        string    `json:"name"`
	TargetValue string    `json:"target_value"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}
