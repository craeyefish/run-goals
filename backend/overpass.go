package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

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

func fetchPeaks() (*OverpassResponse, error) {
	query := `
[out:json];
area["name"="Western Cape"]["admin_level"="4"]->.searchArea;
(
  node["natural"="peak"](area.searchArea);
);
out;
    `
	resp, err := http.Post("https://overpass-api.de/api/interpreter",
		"application/x-www-form-urlencoded",
		strings.NewReader(query),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query overpass: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("overpass request failed: %d", resp.StatusCode)
	}

	var data OverpassResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to parse overpass json: %w", err)
	}
	return &data, nil
}
