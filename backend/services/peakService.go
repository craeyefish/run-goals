package services

import (
	"log"
	"run-goals/daos"
	"run-goals/models"
	"strconv"
)

type PeakServiceInterface interface {
	ListPeaks() ([]models.Peak, error)
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

func (s *PeakService) ListPeaks() ([]models.PeakSummited, error) {
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

	summitedSet := make(map[int64]bool)
	for _, up := range userPeaks {
		summitedSet[up.PeakID] = true
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

		var name string
		var elev float64
		if val, ok := el.Tags["name"]; ok {
			name = val
		}
		if val, ok := el.Tags["ele"]; ok {
			parsedElev, _ := strconv.ParseFloat(val, 64)
			elev = parsedElev
		}

		peak := &models.Peak{
			OsmID:           el.ID,
			Latitude:        el.Lat,
			Longitude:       el.Lon,
			Name:            name,
			ElevationMeters: elev,
		}

		err := s.peaksDao.UpsertPeak(peak)
		if err != nil {
			s.l.Printf("Error calling PeakDao: %v", err)
			return err
		}
	}
	return nil
}
