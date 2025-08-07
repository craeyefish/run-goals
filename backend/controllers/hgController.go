package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"run-goals/daos"
	"run-goals/models"
	"run-goals/services"
	"run-goals/workflows"
	"strings"
)

type HgControllerInterface interface {
	ListHikeGangActivities(rw http.ResponseWriter, r *http.Request)
}

type HgController struct {
	l               *log.Logger
	activityService *services.ActivityService
	userDao         *daos.UserDao
	activityFetcher *workflows.StravaActivityFetcher
}

func NewHgController(
	l *log.Logger,
	activityService *services.ActivityService,
	userDao *daos.UserDao,
	activityFetcher *workflows.StravaActivityFetcher,
) *HgController {
	return &HgController{
		l:               l,
		activityService: activityService,
		userDao:         userDao,
		activityFetcher: activityFetcher,
	}
}

// POST /auth/refresh
func (c *HgController) ListHikeGangActivities(rw http.ResponseWriter, r *http.Request) {
	c.l.Println("Handle GET ListHikeGangActivities")

	u, err := c.userDao.GetUserByStravaAthleteID(int64(3630433))
	if err != nil {
		c.l.Println("Error fetching user", err)
		http.Error(rw, "Failed to fetch user", http.StatusInternalServerError)
		return
	}

	// Call activityService to return list of activities by userId
	activities, err := c.activityService.GetActivitiesByUserID(u.ID)
	if err != nil {
		c.l.Println("Error fetching activities", err)
		http.Error(rw, "Failed to fetch activities", http.StatusInternalServerError)
		return
	}

	// Filter out activities that don't have #hg in the title
	var hgActivities []models.Activity
	for _, activity := range activities {
		titleWords := strings.Split(activity.Name, " ")
		for _, word := range titleWords {
			if strings.ToLower(word) == "#hg" {
				hgActivities = append(hgActivities, activity)
				break
			}
		}
	}

	// Return the array of activities as JSON
	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(hgActivities); err != nil {
		log.Println("Error encoding activities:", err)
	}
}

// TriggerActivitySync manually triggers the activity sync workflow for testing
func (c *HgController) TriggerActivitySync(rw http.ResponseWriter, r *http.Request) {
	c.l.Println("Handle POST TriggerActivitySync - manually triggering activity sync")

	if c.activityFetcher == nil {
		c.l.Println("Error: activityFetcher is nil")
		http.Error(rw, "Activity fetcher not available", http.StatusInternalServerError)
		return
	}

	// Run the activity fetch workflow in a goroutine to avoid blocking
	go func() {
		c.l.Println("Starting manual activity sync...")
		c.activityFetcher.FetchUserActivities()
		c.l.Println("Manual activity sync completed")
	}()

	rw.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"message": "Activity sync triggered successfully",
		"status":  "running",
	}
	json.NewEncoder(rw).Encode(response)
}

// DiagnosticsActivities provides diagnostic information about #hg activities
func (c *HgController) DiagnosticsActivities(rw http.ResponseWriter, r *http.Request) {
	c.l.Println("Handle GET DiagnosticsActivities")

	u, err := c.userDao.GetUserByStravaAthleteID(int64(3630433))
	if err != nil {
		c.l.Println("Error fetching user", err)
		http.Error(rw, "Failed to fetch user", http.StatusInternalServerError)
		return
	}

	// Get all activities for the user
	activities, err := c.activityService.GetActivitiesByUserID(u.ID)
	if err != nil {
		c.l.Println("Error fetching activities", err)
		http.Error(rw, "Failed to fetch activities", http.StatusInternalServerError)
		return
	}

	// Count different types of activities
	totalActivities := len(activities)
	hgActivities := 0
	hgActivitiesWithDescriptions := 0
	totalWithDescriptions := 0

	for _, activity := range activities {
		if activity.Description != "" {
			totalWithDescriptions++
		}
		
		if activity.IsHG() {
			hgActivities++
			if activity.Description != "" {
				hgActivitiesWithDescriptions++
			}
		}
	}

	diagnostics := map[string]interface{}{
		"user_id":                       u.ID,
		"total_activities":              totalActivities,
		"hg_activities":                 hgActivities,
		"hg_activities_with_descriptions": hgActivitiesWithDescriptions,
		"total_activities_with_descriptions": totalWithDescriptions,
		"message": "Diagnostic data for activity descriptions",
	}

	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(diagnostics); err != nil {
		log.Println("Error encoding diagnostics:", err)
	}
}
