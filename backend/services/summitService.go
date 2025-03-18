package services

import (
	"errors"
	"fmt"
	"log"
	"math"
	"run-goals/config"
	"run-goals/daos"
	"run-goals/models"
	"strconv"

	"github.com/twpayne/go-polyline"
)

type SummitServiceInterface interface {
	CandidatePeaks(route string) ([]models.Peak, error)
	IsPeakVisited(route string, peakLat float64, peakLon float64, thresholdMeters float64) bool
	PopulateSummitedPeaks() error
}

type SummitService struct {
	l            *log.Logger
	config       *config.Config
	peaksDao     *daos.PeaksDao
	userPeaksDao *daos.UserPeaksDao
	activityDao  *daos.ActivityDao
}

func NewSummitService(
	l *log.Logger,
	config *config.Config,
	peaksDao *daos.PeaksDao,
	userPeaksDao *daos.UserPeaksDao,
	activityDao *daos.ActivityDao,
) *SummitService {
	return &SummitService{
		l:            l,
		config:       config,
		peaksDao:     peaksDao,
		userPeaksDao: userPeaksDao,
		activityDao:  activityDao,
	}
}

func (s *SummitService) CandidatePeaks(route string) ([]models.Peak, error) {
	var minLat, maxLat, minLon, maxLon float64
	if route == "" {
		return nil, errors.New("no route")
	}

	// After decoding your route
	coords, _, err := polyline.DecodeCoords([]byte(route))
	if err != nil {
		log.Fatalf("Failed to decode polyline: %v", err)
	}

	// Initialize min/max to first point
	minLat, maxLat = coords[0][0], coords[0][0]
	minLon, maxLon = coords[0][1], coords[0][1]

	for _, c := range coords {
		lat := c[0]
		lon := c[1]
		if lat < minLat {
			minLat = lat
		}
		if lat > maxLat {
			maxLat = lat
		}
		if lon < minLon {
			minLon = lon
		}
		if lon > maxLon {
			maxLon = lon
		}
	}

	// You can also expand this box slightly if you want a small buffer
	buffer := 0.01 // ~1 km, depends on latitude/scale
	minLat -= buffer
	maxLat += buffer
	minLon -= buffer
	maxLon += buffer

	candidatePeaks, err := s.peaksDao.GetPeaksBetweenLatLon(minLat, maxLat, minLon, maxLon)
	if err != nil {
		s.l.Printf("Error calling PeaksDao: %v", err)
		return nil, err
	}

	return candidatePeaks, nil
}

func (s *SummitService) IsPeakVisited(route string, peakLat float64, peakLon float64, thresholdMeters float64) bool {
	// DecodeCoords returns a slice of [][2]float64:
	//   coords[i][0] = latitude
	//   coords[i][1] = longitude
	coords, _, err := polyline.DecodeCoords([]byte(route))
	if err != nil {
		log.Fatalf("Failed to decode polyline: %v", err)
	}

	minDist := math.MaxFloat64
	for i := 0; i < len(coords)-1; i++ {
		// distanceToSegment(peakLat, peakLon, route[i], route[i+1])
		segDist := distancePointToSegment(peakLat, peakLon, coords[i][0], coords[i][1], coords[i+1][0], coords[i+1][1])
		if segDist < minDist {
			minDist = segDist
		}
		if minDist < thresholdMeters {
			return true
		}
	}
	return false
}

func distancePointToSegment(px, py, ax, ay, bx, by float64) float64 {
	// Vector AB
	ABx := bx - ax
	ABy := by - ay

	// Avoid a divide-by-zero if A and B are the same point
	lenABsq := ABx*ABx + ABy*ABy
	if lenABsq == 0 {
		// Distance from P to A (which is the same as B)
		return distance(px, py, ax, ay)
	}

	// Vector AP
	APx := px - ax
	APy := py - ay

	// Project AP onto AB, computing parameter t
	dot := APx*ABx + APy*ABy
	t := dot / lenABsq

	// Clamp t to [0, 1]
	if t < 0 {
		// Closest to A
		return distance(px, py, ax, ay)
	} else if t > 1 {
		// Closest to B
		return distance(px, py, bx, by)
	}

	// Projection point
	projX := ax + t*ABx
	projY := ay + t*ABy

	// Distance P->Projection
	return distance(px, py, projX, projY)
}

func distance(x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return math.Sqrt(dx*dx + dy*dy)
}

func (s *SummitService) PopulateSummitedPeaks() error {
	summitThresholdMeters, err := strconv.ParseFloat(s.config.Summit.SummitThresholdMeters, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	// Only do this if there is no data in user peaks.
	// TODO(cian): Replace with a loop that syncs things up daily,
	// it should get populated as activities get added.
	peaks, err := s.userPeaksDao.GetUserPeaks()
	if err != nil {
		return fmt.Errorf("failed to clear user_peaks: %w", err)
	}

	if len(peaks) > 0 {
		return nil
	}

	// fetch all activities
	activities, err := s.activityDao.GetActivities()
	if err != nil {
		return fmt.Errorf("failed to fetch activities: %w", err)
	}

	// loop through activities
	for _, activity := range activities {
		peaks, err := s.CandidatePeaks(activity.MapPolyline)
		if err != nil {
			fmt.Println("failed to fetch candidate peaks: %w", err)
			continue
		}

		// a) get candidate peaks
		// b) is peak visited
		// c) mark user summited peak
		//
		// upsert activity

		var hasSummit bool
		for _, peak := range peaks {
			if s.IsPeakVisited(activity.MapPolyline, peak.Latitude, peak.Longitude, summitThresholdMeters) {
				userPeak := models.UserPeak{
					UserID:     activity.UserID,
					PeakID:     peak.ID,
					ActivityID: activity.ID,
					SummitedAt: activity.StartDate,
				}
				err = s.userPeaksDao.UpsertUserPeak(&userPeak)
				if err != nil {
					log.Printf("Failed to mark summit for user=%d peak=%d: %v\n", activity.UserID, peak.ID, err)
				}

				hasSummit = true
			}
		}

		activity.HasSummit = hasSummit
		err = s.activityDao.UpsertActivity(&activity)
		if err != nil {
			log.Printf("Failed to update activity=%d: %v\n", activity.ID, err)
		}
	}

	return nil
}
