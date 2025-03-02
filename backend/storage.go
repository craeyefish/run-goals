package main

import (
	"errors"
	"log"
	"strconv"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Test

var DB *gorm.DB // global DB handle, or you can pass it around as needed

func InitDB() {
	db, err := gorm.Open(sqlite.Open("myapp.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// Auto-migrate the schema (creates tables if they don't exist)
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Activity{})
	db.AutoMigrate(&Peak{})
	db.AutoMigrate(&UserPeak{})

	DB = db
}

func upsertActivity(stravaAct *StravaActivity, user *User) error {
	// Convert start date
	t, _ := time.Parse(time.RFC3339, stravaAct.StartDate) // handle error properly

	// Check if we already have this activity
	var existing Activity
	result := DB.Where("strava_activity_id = ?", stravaAct.ID).First(&existing)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Create a new record
		newActivity := Activity{
			StravaActivityID: stravaAct.ID,
			UserID:           user.ID,
			Name:             stravaAct.Name,
			Distance:         stravaAct.Distance, // decide if you store in m or km
			StartDate:        t,
			MapPolyline:      stravaAct.Map.SummaryPolyline,
		}
		return DB.Create(&newActivity).Error
	} else if result.Error != nil {
		// some other DB error
		return result.Error
	}

	// If found, update fields as needed
	existing.Name = stravaAct.Name
	existing.Distance = stravaAct.Distance
	existing.StartDate = t
	existing.MapPolyline = stravaAct.Map.SummaryPolyline
	return DB.Save(&existing).Error
}

func storePeaks(resp *OverpassResponse) error {
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

		peak := Peak{
			OsmID: el.ID,
			Lat:   el.Lat,
			Lon:   el.Lon,
			Name:  name,
			ElevM: elev,
		}

		// upsert logic (create if not exists, update if found)
		result := DB.Where("osm_id = ?", peak.OsmID).First(&Peak{})
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			if err := DB.Create(&peak).Error; err != nil {
				log.Printf("Error creating peak: %v", err)
			}
		} else if result.Error == nil {
			// if found, maybe update fields if changed
			if err := DB.Model(&Peak{}).
				Where("osm_id = ?", peak.OsmID).
				Updates(peak).Error; err != nil {
				log.Printf("Error updating peak: %v", err)
			}
		}
	}
	return nil
}

// markUserSummitedPeak upserts a UserPeak record
func markUserSummitedPeak(userID, peakID, activityID uint, summitTime time.Time) error {
	// Create new
	up := UserPeak{
		UserID:     userID,
		PeakID:     peakID,
		ActivityID: activityID,
		SummitedAt: summitTime,
	}
	return DB.Create(&up).Error
}
