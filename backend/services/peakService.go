package services

import (
	"log"
	"run-goals/daos"
	"run-goals/models"
	"strconv"
)

type PeakServiceInterface interface {
	ListPeaks(userID int64) ([]models.PeakSummited, error)
	StorePeaks() (resp *models.OverpassResponse)
}

type PeakService struct {
	l            *log.Logger
	peaksDao     *daos.PeaksDao
	userPeaksDao *daos.UserPeaksDao
}

func NewPeakService(
	l *log.Logger,
	peaksDao *daos.PeaksDao,
	userPeaksDao *daos.UserPeaksDao,
) *PeakService {
	return &PeakService{
		l:            l,
		peaksDao:     peaksDao,
		userPeaksDao: userPeaksDao,
	}
}

func (s *PeakService) ListPeaks(userID int64) ([]models.PeakSummited, error) {
	peaks, err := s.peaksDao.GetPeaks()
	if err != nil {
		s.l.Printf("Error calling PeaksDao: %v", err)
		return nil, err
	}

	userPeaks, err := s.userPeaksDao.GetUserPeaks()
	if err != nil {
		s.l.Printf("Error calling UserPeaksDao: %v", err)
		return nil, err
	}

	// Create a map of peaks summited by this specific user
	summitedSet := make(map[int64]bool)
	for _, up := range userPeaks {
		if up.UserID == userID {
			summitedSet[up.PeakID] = true
		}
	}

	var peaksSummited []models.PeakSummited
	for _, p := range peaks {
		isSummited := summitedSet[p.ID]

		peakSummited := models.PeakSummited{
			Peak:       p,
			IsSummited: isSummited,
		}
		peaksSummited = append(peaksSummited, peakSummited)
	}

	return peaksSummited, nil
}

func (s *PeakService) StorePeaks(resp *models.OverpassResponse) error {
	if resp == nil {
		return nil
	}

	for _, el := range resp.Elements {
		if el.Type != "node" {
			continue
		}

		// Parse standard fields
		var name string
		var elev float64
		if val, ok := el.Tags["name"]; ok {
			name = val
		}
		if val, ok := el.Tags["ele"]; ok {
			parsedElev, _ := strconv.ParseFloat(val, 64)
			elev = parsedElev
		}

		// Parse additional metadata for differentiation
		var altName, nameEN, region, wikipedia, wikidata, description string
		var prominence float64

		if val, ok := el.Tags["alt_name"]; ok {
			altName = val
		}
		if val, ok := el.Tags["name:en"]; ok {
			nameEN = val
		}
		// Try multiple tags for region info
		if val, ok := el.Tags["is_in"]; ok {
			region = val
		} else if val, ok := el.Tags["is_in:region"]; ok {
			region = val
		} else if val, ok := el.Tags["is_in:mountain_range"]; ok {
			region = val
		}
		if val, ok := el.Tags["wikipedia"]; ok {
			wikipedia = val
		}
		if val, ok := el.Tags["wikidata"]; ok {
			wikidata = val
		}
		if val, ok := el.Tags["description"]; ok {
			description = val
		}
		if val, ok := el.Tags["prominence"]; ok {
			parsedProm, _ := strconv.ParseFloat(val, 64)
			prominence = parsedProm
		}

		peak := &models.Peak{
			OsmID:           el.ID,
			Latitude:        el.Lat,
			Longitude:       el.Lon,
			Name:            name,
			ElevationMeters: elev,
			AltName:         altName,
			NameEN:          nameEN,
			Region:          region,
			Wikipedia:       wikipedia,
			Wikidata:        wikidata,
			Description:     description,
			Prominence:      prominence,
		}

		err := s.peaksDao.UpsertPeak(peak)
		if err != nil {
			s.l.Printf("Error calling PeakDao: %v", err)
			return err
		}
	}
	return nil
}
