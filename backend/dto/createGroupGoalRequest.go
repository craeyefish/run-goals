package dto

import "time"

type CreateGroupGoalRequest struct {
	GroupID       int64     `json:"group_id"`
	Name          string    `json:"name"`
	Description   *string   `json:"description,omitempty"`
	GoalType      string    `json:"goal_type"`
	TargetValue   float64   `json:"target_value"`
	TargetSummits []int64   `json:"target_summits,omitempty"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
}
