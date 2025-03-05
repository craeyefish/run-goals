package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"run-goals/services"
	"strconv"
)

type ActivityController struct {
	l       *log.Logger
	service *services.ActivityService
}

func NewActivityController(l *log.Logger, db *sql.DB) *ActivityController {
	return &ActivityController{
		l:       l,
		service: services.NewActivityService(l, db),
	}
}

func (c *ActivityController) ListActivities(rw http.ResponseWriter, r *http.Request) {
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
	activities, err := c.service.GetActivitiesByUserID(userID)
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
