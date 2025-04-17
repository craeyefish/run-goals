package models

import (
	"time"
)

type GroupMemberGoalContribution struct {
	GroupMemberID      int64     `json:"group_member_id"`
	GroupID            int64     `json:"group_id"`
	UserID             int64     `json:"user_id"`
	Role               string    `json:"role"`
	JoinedAt           time.Time `json:"joined_at"`
	TotalActivities    int64     `json:"total_activities"`
	TotalDistance      float64   `json:"total_distance"`
	TotalUniqueSummits int64     `json:"total_unique_summits"`
	TotalSummits       int64     `json:"total_summits"`
}
