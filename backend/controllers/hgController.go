package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"run-goals/daos"
	"run-goals/models"
	"run-goals/services"
	"strings"
)

type HgControllerInterface interface {
	ListHikeGangActivities(rw http.ResponseWriter, r *http.Request)
}

type HgController struct {
	l               *log.Logger
	activityService *services.ActivityService
	userDao         *daos.UserDao
}

func NewHgController(
	l *log.Logger,
	activityService *services.ActivityService,
	userDao *daos.UserDao,
) *HgController {
	return &HgController{
		l:               l,
		activityService: activityService,
		userDao:         userDao,
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
