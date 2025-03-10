package models

type Peak struct {
	ID              int64   `json:"id"`
	OsmID           int64   `json:"osm_id"` // OSM Node ID
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	Name            string  `json:"name"`
	ElevationMeters float64 `json:"elevation_meters"` // Elevation in meters (parse from "ele" tag if present)
}
