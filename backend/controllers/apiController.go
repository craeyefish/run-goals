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
	l                       *log.Logger
	activityService         *services.ActivityService
	progressService         *services.ProgressService
	peakService             *services.PeakService
	summariesService        *services.SummariesService
	userService             *services.UserService
	personalGoalsService    *services.PersonalGoalsService
	summitFavouritesService *services.SummitFavouritesService
}

func NewApiController(
	l *log.Logger,
	activityService *services.ActivityService,
	progressService *services.ProgressService,
	peakService *services.PeakService,
	summariesService *services.SummariesService,
	userService *services.UserService,
	personalGoalsService *services.PersonalGoalsService,
	summitFavouritesService *services.SummitFavouritesService,
) *ApiController {
	return &ApiController{
		l:                       l,
		activityService:         activityService,
		progressService:         progressService,
		peakService:             peakService,
		summariesService:        summariesService,
		userService:             userService,
		personalGoalsService:    personalGoalsService,
		summitFavouritesService: summitFavouritesService,
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

// GetSummitFavourites returns all favourite peak IDs for the user
// GET /api/summit-favourites
func (c *ApiController) GetSummitFavourites(rw http.ResponseWriter, r *http.Request) {
	c.l.Println("Handle GET SummitFavourites")

	userID, _ := meta.GetUserIDFromContext(r.Context())

	peakIDs, err := c.summitFavouritesService.GetFavourites(userID)
	if err != nil {
		c.l.Printf("Error fetching summit favourites: %v", err)
		http.Error(rw, "Failed to fetch summit favourites", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(peakIDs); err != nil {
		log.Println("Error encoding summit favourites response:", err)
	}
}

// AddSummitFavourite adds a peak to user's favourites
// POST /api/summit-favourites
func (c *ApiController) AddSummitFavourite(rw http.ResponseWriter, r *http.Request) {
	c.l.Println("Handle POST AddSummitFavourite")

	userID, _ := meta.GetUserIDFromContext(r.Context())

	var req struct {
		PeakID int64 `json:"peak_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.l.Printf("Error decoding request body: %v", err)
		http.Error(rw, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := c.summitFavouritesService.AddFavourite(userID, req.PeakID); err != nil {
		c.l.Printf("Error adding summit favourite: %v", err)
		http.Error(rw, "Failed to add summit favourite", http.StatusInternalServerError)
		return
	}

	// Return updated list
	peakIDs, err := c.summitFavouritesService.GetFavourites(userID)
	if err != nil {
		c.l.Printf("Error fetching summit favourites: %v", err)
		http.Error(rw, "Failed to fetch summit favourites", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(peakIDs); err != nil {
		log.Println("Error encoding summit favourites response:", err)
	}
}

// RemoveSummitFavourite removes a peak from user's favourites
// DELETE /api/summit-favourites/:peak_id
func (c *ApiController) RemoveSummitFavourite(rw http.ResponseWriter, r *http.Request) {
	c.l.Println("Handle DELETE RemoveSummitFavourite")

	userID, _ := meta.GetUserIDFromContext(r.Context())

	// Get peak_id from query params
	peakIDStr := r.URL.Query().Get("peak_id")
	peakID, err := strconv.ParseInt(peakIDStr, 10, 64)
	if err != nil {
		c.l.Printf("Invalid peak_id: %v", err)
		http.Error(rw, "Invalid peak_id", http.StatusBadRequest)
		return
	}

	if err := c.summitFavouritesService.RemoveFavourite(userID, peakID); err != nil {
		c.l.Printf("Error removing summit favourite: %v", err)
		http.Error(rw, "Failed to remove summit favourite", http.StatusInternalServerError)
		return
	}

	// Return updated list
	peakIDs, err := c.summitFavouritesService.GetFavourites(userID)
	if err != nil {
		c.l.Printf("Error fetching summit favourites: %v", err)
		http.Error(rw, "Failed to fetch summit favourites", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(peakIDs); err != nil {
		log.Println("Error encoding summit favourites response:", err)
	}
}
