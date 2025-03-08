package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"run-goals/services"
	"strconv"
)

type ApiControllerInterface interface {
	ListActivities(rw http.ResponseWriter, r *http.Request)
	ListPeaks(rw http.ResponseWriter, r *http.Request)
	Progress(rw http.ResponseWriter, r *http.Request)
	PeakSummaries(rw http.ResponseWriter, r *http.Request)
}

type ApiController struct {
	l               *log.Logger
	activityService *services.ActivityService
	progressService *services.ProgressService
}

func NewApiController(l *log.Logger, db *sql.DB) *ApiController {
	return &ApiController{
		l:               l,
		activityService: services.NewActivityService(l, db),
		progressService: services.NewProgressService(l, db),
	}
}

func (c *ApiController) ListActivities(rw http.ResponseWriter, r *http.Request) {
	c.l.Println("Handle GET ListActivities")

	// extract user id from url
	userIDStr := r.URL.Query().Get("userId")
	if userIDStr == "" {
		http.Error(rw, "missing userId", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(rw, "Invalid userId", http.StatusBadRequest)
		return
	}

	// Call activityService to return list of activies by userId
	activities, err := c.activityService.GetActivitiesByUserID(userID)
	if err != nil {
		c.l.Println("Error fetching activities", err)
		http.Error(rw, "Failed to fetch activities", http.StatusInternalServerError)
		return
	}

	// Return the array of activities as JSON
	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(activities); err != nil {
		log.Println("Error encoding activities:", err)
	}
}

func (c *ApiController) ListPeaks(rw http.ResponseWriter, r *http.Request) {
	c.l.Println("Handle GET ListPeaks")
}

func (c *ApiController) Progress(rw http.ResponseWriter, r *http.Request) {
	c.l.Println("Handle GET Progress")

	rw.Header().Set("Content-Type", "application/json")
}

func (c *ApiController) PeakSummaries(rw http.ResponseWriter, r *http.Request) {
	c.l.Println("Handle GET PeakSummaries")
}
