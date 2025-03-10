package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"run-goals/daos"
	"run-goals/models"
	"strings"
)

type OverpassServiceInterface interface {
	FetchPeaks() error
}

type OverpassService struct {
	l        *log.Logger
	peaksDao *daos.PeaksDao
}

func NewOverpassService(
	l *log.Logger,
	peaksDao *daos.PeaksDao,
) *OverpassService {
	return &OverpassService{
		l:        l,
		peaksDao: peaksDao,
	}
}

func (s *OverpassService) FetchPeaks() error {
	peaks, err := s.peaksDao.GetPeaks()
	if err != nil {
		s.l.Printf("Error calling PeaksDao: %v", err)
		return err
	}

	if len(peaks) > 0 {
		return nil
	}

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
		return fmt.Errorf("failed to query overpass: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("overpass request failed: %d", resp.StatusCode)
	}

	var data models.OverpassResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fmt.Errorf("failed to parse overpass json: %w", err)
	}

	return nil
}
