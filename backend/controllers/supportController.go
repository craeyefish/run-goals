package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"run-goals/daos"
	"run-goals/meta"
	"run-goals/services"
	"strconv"
	"strings"
)

type SupportController struct {
	l           *log.Logger
	userService *services.UserService
}

func NewSupportController(
	l *log.Logger,
	userService *services.UserService,
) *SupportController {
	return &SupportController{
		l:           l,
		userService: userService,
	}
}

func (c *SupportController) DeleteUserAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract strava athlete ID from URL path
	// Expected path: /support/delete-account/{stravaAthleteID}
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/support/delete-account/"), "/")
	if len(pathParts) == 0 || pathParts[0] == "" {
		c.l.Printf("Missing strava athlete ID in delete request")
		http.Error(w, "Strava athlete ID is required", http.StatusBadRequest)
		return
	}

	stravaAthleteIDStr := pathParts[0]
	stravaAthleteID, err := strconv.ParseInt(stravaAthleteIDStr, 10, 64)
	if err != nil {
		c.l.Printf("Invalid strava athlete ID format: %s", stravaAthleteIDStr)
		http.Error(w, "Invalid strava athlete ID format", http.StatusBadRequest)
		return
	}

	userID, _ := meta.GetUserIDFromContext(r.Context())
	user, err := c.userService.GetUserByID(userID)
	if err != nil {
		c.l.Printf("Error fetching user by ID %d: %v", userID, err)
		http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
		return
	}

	if user.StravaAthleteID != stravaAthleteID {
		c.l.Printf("User %d attempted to delete account with mismatched strava athlete ID: %d", userID, stravaAthleteID)
		http.Error(w, "You are not authorized to delete this account", http.StatusForbidden)
		return
	}

	c.l.Printf("Processing account deletion request for strava_athlete_id: %d", stravaAthleteID)

	// Delete the user account
	err = c.userService.DeleteUserAccount(stravaAthleteID)
	if err != nil {
		if err == daos.ErrUserNotFound {
			c.l.Printf("User not found for deletion: strava_athlete_id=%d", stravaAthleteID)
			response := map[string]string{
				"message": "No account found with that Strava Athlete ID",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(response)
			return
		}

		c.l.Printf("Error deleting user account: %v", err)
		response := map[string]string{
			"message": "Failed to delete account. Please try again later.",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Success response
	response := map[string]string{
		"message": "Account successfully deleted",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	c.l.Printf("Successfully processed account deletion for strava_athlete_id: %d", stravaAthleteID)
}
