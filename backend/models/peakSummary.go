package models

type PeakSummary struct {
	PeakID   int64              `json:"peak_id"`
	PeakName string             `json:"peak_name"`
	Summits  []SummitedActivity `json:"summits"`
}
