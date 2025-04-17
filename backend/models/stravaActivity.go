package models

// StravaActivity is a partial struct to parse activity JSON from Strava
type StravaActivity struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Distance  float64 `json:"distance"` // in meters
	StartDate string  `json:"start_date_local"`
	Map       struct {
		SummaryPolyline string `json:"summary_polyline"`
	} `json:"map"`
	Photos struct {
		Count   int `json:"count"`
		Primary struct {
			ID       string            `json:"id"`
			Source   int               `json:"source"`
			UniqueID string            `json:"unique_id"`
			Urls     map[string]string `json:"urls"`
		} `json:"primary"`
	} `json:"photos"`
	// add other fields if needed
}
