package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"gorm.io/gorm"
)

// TODO(cian): Make group setup process and store this in db.
const groupGoal = 1000.0

func handleProgress(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 1. Load all users from DB
	var users []User
	if err := DB.Find(&users).Error; err != nil {
		log.Println("Error fetching users from DB:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var total float64
	contributions := make([]UserContribution, 0, len(users))

	var wg sync.WaitGroup
	var mu sync.Mutex

	wg.Add(len(users))
	for _, user := range users {
		go func(u User) {
			defer wg.Done()

			dist, err := getUserDistance(&u)
			if err != nil {
				log.Println("Error fetching distance for user:", u.ID, err)
				return
			}

			fetchAndStoreUserActivities(&u)
			if err != nil {
				log.Println("Error fetching activities for user:", u.ID, err)
				return
			}

			mu.Lock()
			total += dist
			contributions = append(contributions, UserContribution{
				ID:            u.StravaAthleteID,
				TotalDistance: dist,
			})
			mu.Unlock()
		}(user)
	}

	wg.Wait()

	// Construct response
	response := GoalProgress{
		Goal:            groupGoal,
		CurrentProgress: total,
		Contributions:   contributions,
	}

	json.NewEncoder(w).Encode(response)
}

func handleListActivities(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userIDStr := r.URL.Query().Get("userId")
	if userIDStr == "" {
		http.Error(w, "missing userId", http.StatusBadRequest)
		return
	}

	var activities []Activity
	if err := DB.Where("user_id = ?", userIDStr).Find(&activities).Error; err != nil {
		log.Println("DB error fetching activities:", err)
		http.Error(w, "failed to fetch activities", http.StatusInternalServerError)
		return
	}

	// Return the array of activities as JSON
	if err := json.NewEncoder(w).Encode(activities); err != nil {
		log.Println("Error encoding activities:", err)
	}
}

type ResponsePeak struct {
	Peak
	IsSummited bool `json:"is_summited"`
}

func handleListPeaks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var peaks []Peak
	if err := DB.Find(&peaks).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var userPeaks []UserPeak
	if err := DB.Find(&userPeaks).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	summitedSet := make(map[uint]bool)
	for _, up := range userPeaks {
		summitedSet[up.PeakID] = true
	}

	var respPeaks []ResponsePeak
	for _, p := range peaks {
		isSummited := summitedSet[p.ID]

		respPeak := ResponsePeak{
			Peak:       p,
			IsSummited: isSummited,
		}
		respPeaks = append(respPeaks, respPeak)
	}

	json.NewEncoder(w).Encode(respPeaks)
}

func handlePeakSummaries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 1. Load all actual Peak records so we can return them even if they have 0 summits.
	var peaks []Peak
	if err := DB.Find(&peaks).Error; err != nil {
		http.Error(w, "Failed to load peaks", http.StatusInternalServerError)
		return
	}

	// 2. Fetch joined data from user_peaks + users for summits
	type userPeakJoin struct {
		PeakID     uint      // from up.peak_id
		UserID     uint      // from up.user_id
		ActivityID uint      // from up.activity_id
		SummitedAt time.Time // from up.summited_at
		UserName   uint      // from u.strava_athlete_id (as an example)
	}

	var joined []userPeakJoin

	// IMPORTANT: select up.peak_id, not up.id, so we can group by the actual peak.
	err := DB.Raw(`
		SELECT 
		   up.peak_id, 
		   up.user_id, 
		   up.activity_id, 
		   up.summited_at, 
		   u.strava_athlete_id AS user_name
		FROM user_peaks up
		JOIN users u ON up.user_id = u.id
	`).Scan(&joined).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, "Failed to load user_peaks join", http.StatusInternalServerError)
		log.Println("Error fetching user_peaks join:", err)
		return
	}

	// 3. Build a map: peak_id -> []SummitedActivity
	summitsByPeak := make(map[uint][]SummitedActivity)

	for _, row := range joined {
		// Fetch the actual Activity for more details (like StravaActivityID).
		// Note the correct chaining: DB.Where(...).First(&activity)
		var activity Activity
		if err := DB.Where("id = ?", row.ActivityID).First(&activity).Error; err != nil {
			log.Println("Error fetching activity:", err)
			continue // skip if we can't find the activity
		}

		summit := SummitedActivity{
			UserName:   strconv.Itoa(int(row.UserName)), // if you want a string
			UserID:     row.UserID,
			ActivityID: uint(activity.StravaActivityID),
			SummitedAt: row.SummitedAt,
		}

		summitsByPeak[row.PeakID] = append(summitsByPeak[row.PeakID], summit)
	}

	// 4. Build the final JSON structure by iterating over all peaks
	var results []PeakSummary

	for _, p := range peaks {
		if len(summitsByPeak[p.ID]) == 0 {
			continue
		}

		results = append(results, PeakSummary{
			PeakID:   p.ID,
			PeakName: p.Name,
			Summits:  summitsByPeak[p.ID], // may be nil or empty if no summits
		})
	}

	// 5. Encode to JSON
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}
}
