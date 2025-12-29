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
	CalculateSummitsForActivity(activity *models.Activity) error
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
		s.l.Printf("Failed to decode polyline: %v", err)
		return nil, fmt.Errorf("invalid polyline data: %w", err)
	}

	if len(coords) == 0 {
		return nil, errors.New("polyline decoded to empty coordinates")
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
	// Fetch only activities that haven't been processed yet
	activities, err := s.activityDao.GetActivitiesPendingSummitCalculation()
	if err != nil {
		return fmt.Errorf("failed to fetch pending activities: %w", err)
	}

	if len(activities) == 0 {
		s.l.Println("No activities pending summit calculation")
		return nil
	}

	s.l.Printf("Processing summit detection for %d activities", len(activities))

	for _, activity := range activities {
		if err := s.CalculateSummitsForActivity(&activity); err != nil {
			s.l.Printf("Failed to calculate summits for activity %d: %v", activity.ID, err)
			// Continue with other activities even if one fails
		}
	}

	return nil
}

// CalculateSummitsForActivity processes a single activity for summit detection
func (s *SummitService) CalculateSummitsForActivity(activity *models.Activity) error {
	summitThresholdMeters, err := strconv.ParseFloat(s.config.Summit.SummitThresholdMeters, 64)
	if err != nil {
		return fmt.Errorf("invalid summit threshold config: %w", err)
	}

	// Check if the route (MapPolyline) is empty
	if activity.MapPolyline == "" {
		s.l.Printf("Skipping activity %d for user %d: no route provided", activity.ID, activity.UserID)
		// Mark as calculated even though no route - nothing to do
		activity.SummitsCalculated = true
		return s.activityDao.UpsertActivity(activity)
	}

	// Fetch candidate peaks
	peaks, err := s.CandidatePeaks(activity.MapPolyline)
	if err != nil {
		s.l.Printf("Failed to fetch candidate peaks for activity %d: %v", activity.ID, err)
		// Still mark as calculated to avoid retrying bad polylines
		activity.SummitsCalculated = true
		return s.activityDao.UpsertActivity(activity)
	}

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
				s.l.Printf("Failed to mark summit for user=%d peak=%d: %v", activity.UserID, peak.ID, err)
			} else {
				s.l.Printf("Summit detected! user=%d peak=%d (%s) activity=%d", activity.UserID, peak.ID, peak.Name, activity.ID)
			}
			hasSummit = true
		}
	}

	activity.HasSummit = hasSummit
	activity.SummitsCalculated = true
	return s.activityDao.UpsertActivity(activity)
}
