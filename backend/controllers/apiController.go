package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"run-goals/meta"
	"run-goals/models"
	"run-goals/services"
	"strconv"
	"time"
)

type ApiControllerInterface interface {
	ListActivities(rw http.ResponseWriter, r *http.Request)
	ListPeaks(rw http.ResponseWriter, r *http.Request)
	GetProgress(rw http.ResponseWriter, r *http.Request)
	GetPeakSummaries(rw http.ResponseWriter, r *http.Request)
	GetUserProfile(rw http.ResponseWriter, r *http.Request)
	GetPersonalGoals(rw http.ResponseWriter, r *http.Request)
	SavePersonalGoals(rw http.ResponseWriter, r *http.Request)
}

type ApiController struct {
	l                    *log.Logger
	activityService      *services.ActivityService
	progressService      *services.ProgressService
	peakService          *services.PeakService
	summariesService     *services.SummariesService
	userService          *services.UserService
	personalGoalsService *services.PersonalGoalsService
}

func NewApiController(
	l *log.Logger,
	activityService *services.ActivityService,
	progressService *services.ProgressService,
	peakService *services.PeakService,
	summariesService *services.SummariesService,
	userService *services.UserService,
	personalGoalsService *services.PersonalGoalsService,
) *ApiController {
	return &ApiController{
		l:                    l,
		activityService:      activityService,
		progressService:      progressService,
		peakService:          peakService,
		summariesService:     summariesService,
		userService:          userService,
		personalGoalsService: personalGoalsService,
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

	userID, _ := meta.GetUserIDFromContext(r.Context())

	response, err := c.peakService.ListPeaks(userID)
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

	userID, _ := meta.GetUserIDFromContext(r.Context())

	response, err := c.summariesService.GetPeakSummaries(userID)
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

// GetPersonalGoals returns the user's personal yearly goals
// GET /api/personal-goals?year=2025 (defaults to current year if not specified)
func (c *ApiController) GetPersonalGoals(rw http.ResponseWriter, r *http.Request) {
	c.l.Println("Handle GET PersonalGoals")

	userID, _ := meta.GetUserIDFromContext(r.Context())

	// Get year from query params, default to current year
	yearStr := r.URL.Query().Get("year")
	year := time.Now().Year()
	if yearStr != "" {
		if parsedYear, err := strconv.Atoi(yearStr); err == nil {
			year = parsedYear
		}
	}

	goal, err := c.personalGoalsService.GetGoalForYear(userID, year)
	if err != nil {
		c.l.Printf("Error fetching personal goals: %v", err)
		http.Error(rw, "Failed to fetch personal goals", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(goal); err != nil {
		log.Println("Error encoding personal goals response:", err)
	}
}

// SavePersonalGoals creates or updates the user's personal yearly goals
// POST /api/personal-goals
func (c *ApiController) SavePersonalGoals(rw http.ResponseWriter, r *http.Request) {
	c.l.Println("Handle POST PersonalGoals")

	userID, _ := meta.GetUserIDFromContext(r.Context())

	var goal models.PersonalYearlyGoal
	if err := json.NewDecoder(r.Body).Decode(&goal); err != nil {
		c.l.Printf("Error decoding request body: %v", err)
		http.Error(rw, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Ensure the goal belongs to this user
	goal.UserID = userID

	// Default to current year if not specified
	if goal.Year == 0 {
		goal.Year = time.Now().Year()
	}

	if err := c.personalGoalsService.SaveGoal(&goal); err != nil {
		c.l.Printf("Error saving personal goals: %v", err)
		http.Error(rw, "Failed to save personal goals", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(goal); err != nil {
		log.Println("Error encoding personal goals response:", err)
	}
}

// GetAllPersonalGoals returns all yearly goals for the user (for history view)
// GET /api/personal-goals/all
func (c *ApiController) GetAllPersonalGoals(rw http.ResponseWriter, r *http.Request) {
	c.l.Println("Handle GET AllPersonalGoals")

	userID, _ := meta.GetUserIDFromContext(r.Context())

	goals, err := c.personalGoalsService.GetAllGoals(userID)
	if err != nil {
		c.l.Printf("Error fetching all personal goals: %v", err)
		http.Error(rw, "Failed to fetch personal goals", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(goals); err != nil {
		log.Println("Error encoding personal goals response:", err)
	}
}
