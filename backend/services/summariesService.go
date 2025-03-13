package services

import (
	"log"
	"run-goals/daos"
	"run-goals/models"
	"strconv"
)

type SummariesServiceInterface interface {
	GetPeakSummaries() ([]models.PeakSummary, error)
}

type SummariesService struct {
	l            *log.Logger
	peaksDao     *daos.PeaksDao
	userPeaksDao *daos.UserPeaksDao
	activityDao  *daos.ActivityDao
}

func NewSummariesService(
	l *log.Logger,
	peaksDao *daos.PeaksDao,
	userPeaksDao *daos.UserPeaksDao,
	activityDao *daos.ActivityDao,
) *SummariesService {
	return &SummariesService{
		l:            l,
		peaksDao:     peaksDao,
		userPeaksDao: userPeaksDao,
		activityDao:  activityDao,
	}
}

func (s *SummariesService) GetPeakSummaries() ([]models.PeakSummary, error) {
	// 1. Load all actual Peak records so we can return them even if they have 0 summits.
	peaks, err := s.peaksDao.GetPeaks()
	if err != nil {
		s.l.Printf("Error calling peaksDao: %v", err)
		return nil, err
	}

	// 2. Fetch joined data from user_peaks + users for summits
	userPeakJoined, err := s.userPeaksDao.GetUserPeaksJoin()
	if err != nil {
		s.l.Printf("Error calling userPeaksDao: %v", err)
		return nil, err
	}

	// 3. Build a map: peak_id -> []SummitedActivity
	summitsByPeak := make(map[int64][]models.SummitedActivity)
	for _, row := range userPeakJoined {
		// Fetch the actual Activity for more details (like StravaActivityID).
		activity, err := s.activityDao.GetActivityByID(row.ActivityID)
		if err != nil {
			s.l.Printf("Error calling activityDao.GetActivitiesByID: %v", err)
			continue
		}
		summit := models.SummitedActivity{
			UserName:   strconv.Itoa(int(row.UserName)), // if you want a string
			UserID:     row.UserID,
			ActivityID: activity.ID,
			SummitedAt: row.SummitedAt,
		}
		summitsByPeak[row.PeakID] = append(summitsByPeak[row.PeakID], summit)
	}

	// 4. Build the final JSON structure by iterating over all peaks
	var peakSummaries []models.PeakSummary
	for _, p := range peaks {
		if len(summitsByPeak[p.ID]) == 0 {
			continue
		}
		peakSummaries = append(peakSummaries, models.PeakSummary{
			PeakID:   p.ID,
			PeakName: p.Name,
			Summits:  summitsByPeak[p.ID], // may be nil or empty if no summits
		})
	}
	return peakSummaries, nil
}
