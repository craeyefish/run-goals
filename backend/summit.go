package main

import (
	"errors"
	"fmt"
	"log"
	"math"

	"github.com/twpayne/go-polyline"
	"gorm.io/gorm"
)

func candidatePeaks(route string) ([]Peak, error) {
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

	var candidatePeaks []Peak
	err = DB.Where("lat BETWEEN ? AND ? AND lon BETWEEN ? AND ?",
		minLat, maxLat, minLon, maxLon,
	).Find(&candidatePeaks).Error
	if err != nil {
		return nil, err
	}

	return candidatePeaks, nil
}

func isPeakVisited(route string, peakLat float64, peakLon float64, thresholdMeters float64) bool {
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

// The distance threshold (in meters) to consider a summit "bagged"
const SummitThresholdMeters = 0.0007 // ~ 70 m

// Retroactively populate summited peaks data for all activities in the DB
func PopulateSummitedPeaks() error {
	if err := DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&UserPeak{}).Error; err != nil {
		return fmt.Errorf("failed to clear user_peaks: %w", err)
	}

	var activities []Activity
	if err := DB.Find(&activities).Error; err != nil {
		return fmt.Errorf("failed to fetch activities: %w", err)
	}

	for _, activity := range activities {
		peaks, err := candidatePeaks(activity.MapPolyline)
		if err != nil {
			fmt.Println("failed to fetch candidate peaks: %w", err)
			continue
		}

		var hasSummit bool
		for _, peak := range peaks {
			if isPeakVisited(activity.MapPolyline, peak.Lat, peak.Lon, SummitThresholdMeters) {
				err = markUserSummitedPeak(activity.UserID, peak.ID, activity.ID, activity.StartDate)
				if err != nil {
					log.Printf("Failed to mark summit for user=%d peak=%d: %v\n", activity.UserID, peak.ID, err)
				}

				hasSummit = true
			}
		}

		activity.HasSummit = hasSummit
		err = DB.Save(&activity).Error
		if err != nil {
			log.Printf("Failed to update activity=%d: %v\n", activity.ID, err)
		}
	}

	return nil
}
