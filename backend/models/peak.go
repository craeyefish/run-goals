package models

type Peak struct {
	ID              int64   `json:"id"`
	OsmID           int64   `json:"osm_id"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	Name            string  `json:"name"`
	ElevationMeters float64 `json:"elevation_meters"`
	// New fields for better differentiation
	AltName     string `json:"alt_name"`      // Alternative name (from alt_name tag)
	NameEN      string `json:"name_en"`       // English name (from name:en tag)
	Region      string `json:"region"`        // Region/area (derived or from is_in tag)
	Wikipedia   string `json:"wikipedia"`     // Wikipedia article link
	Wikidata    string `json:"wikidata"`      // Wikidata ID for more info
	Description string `json:"description"`   // From description tag
	Prominence  float64 `json:"prominence"`   // From prominence tag if available
}
