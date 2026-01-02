package models

import "time"

// SummitFavourite represents a user's favourite/wishlist peak
type SummitFavourite struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	PeakID    int64     `json:"peak_id" db:"peak_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
