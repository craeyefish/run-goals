package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"run-goals/meta"
	"run-goals/services"
)

type ApiControllerInterface interface {
	ListActivities(rw http.ResponseWriter, r *http.Request)
	ListPeaks(rw http.ResponseWriter, r *http.Request)
	GetProgress(rw http.ResponseWriter, r *http.Request)
	GetPeakSummaries(rw http.ResponseWriter, r *http.Request)
	GetUserProfile(rw http.ResponseWriter, r *http.Request)
}

type ApiController struct {
	l                *log.Logger
	activityService  *services.ActivityService
	progressService  *services.ProgressService
	peakService      *services.PeakService
	summariesService *services.SummariesService
	userService      *services.UserService
}

func NewApiController(
	l *log.Logger,
	activityService *services.ActivityService,
	progressService *services.ProgressService,
	peakService *services.PeakService,
	summariesService *services.SummariesService,
	userService *services.UserService,
) *ApiController {
	return &ApiController{
		l:                l,
		activityService:  activityService,
		progressService:  progressService,
		peakService:      peakService,
		summariesService: summariesService,
		userService:      userService,
	}
}

func (c *ApiController) ListActivities(rw http.ResponseWriter, r *http.Request) {
	c.l.Println("Handle GET ListActivities")

	userID, _ := meta.GetUserIDFromContext(r.Context())

	// Call activityService to return list of activities by userId
	response, err := c.activityService.GetActivitiesByUserID(userID)
	if err != nil {
		c.l.Println("Error fetching activities", err)
		http.Error(rw, "Failed to fetch activities", http.StatusInternalServerError)
		return
	}

	// Return the array of activities as JSON
	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(response); err != nil {
		log.Println("Error encoding activities:", err)
	}
}

func (c *ApiController) ListPeaks(rw http.ResponseWriter, r *http.Request) {
	c.l.Println("Handle GET ListPeaks")

	response, err := c.peakService.ListPeaks()
	if err != nil {
		c.l.Println("Error listing peaks", err)
		http.Error(rw, "Failed to list peaks", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(response); err != nil {
		log.Println("Error encoding listPeaksResponse:", err)
	}
}

func (c *ApiController) GetProgress(rw http.ResponseWriter, r *http.Request) {
	c.l.Println("Handle GET Progress")

	response, err := c.progressService.GetUsersProgress()
	if err != nil {
		c.l.Println("Error listing peaks", err)
		http.Error(rw, "Failure calling progressService", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(response); err != nil {
		log.Println("Error encoding GoalProgress response:", err)
	}
}

func (c *ApiController) GetPeakSummaries(rw http.ResponseWriter, r *http.Request) {
	c.l.Println("Handle GET PeakSummaries")

	response, err := c.summariesService.GetPeakSummaries()
	if err != nil {
		c.l.Println("Error listing peaks", err)
		http.Error(rw, "Failure calling summariesService", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(response); err != nil {
		log.Println("Error encoding Peak Summaries response:", err)
	}
}

func (c *ApiController) GetUserProfile(rw http.ResponseWriter, r *http.Request) {
	c.l.Println("Handle GET UserProfile")

	userID, _ := meta.GetUserIDFromContext(r.Context())

	response, err := c.userService.GetUserProfile(userID)
	if err != nil {
		c.l.Println("Error fetching user profile", err)
		http.Error(rw, "Failed to fetch user profile", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(response); err != nil {
		log.Println("Error encoding user profile response:", err)
	}
}
