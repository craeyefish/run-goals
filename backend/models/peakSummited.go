package models

type PeakSummited struct {
	Peak
	IsSummited bool `json:"is_summited"`
}
