package main

import "time"

type SummitedActivity struct {
	UserName   string    `json:"user_name"`
	UserID     uint      `json:"user_id"`
	ActivityID uint      `json:"activity_id"`
	SummitedAt time.Time `json:"summited_at"`
}

type PeakSummary struct {
	PeakID   uint               `json:"peak_id"`
	PeakName string             `json:"peak_name"`
	Summits  []SummitedActivity `json:"summits"`
}
