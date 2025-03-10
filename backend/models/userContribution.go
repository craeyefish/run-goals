package models

type UserContribution struct {
	ID            int64   `json:"id"`
	TotalDistance float64 `json:"totalDistance"` // in kilometers
}
