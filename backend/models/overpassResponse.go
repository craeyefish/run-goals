package models

type OverpassResponse struct {
	Version   float64           `json:"version"`
	Generator string            `json:"generator"`
	Osm3s     map[string]string `json:"osm3s"` // or a custom struct if needed
	Elements  []Element         `json:"elements"`
}

type Element struct {
	Type string            `json:"type"`
	ID   int64             `json:"id"`
	Lat  float64           `json:"lat"`
	Lon  float64           `json:"lon"`
	Tags map[string]string `json:"tags"`
}
