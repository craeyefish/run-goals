package services

import (
	"log"
	"run-goals/daos"
	"run-goals/models"
)

type PeakServiceInterface interface {
	ListPeaks() ([]models.Peak, error)
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
