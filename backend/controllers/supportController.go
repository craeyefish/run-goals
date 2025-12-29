package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"run-goals/daos"
	"run-goals/meta"
	"run-goals/services"
	"strconv"
	"strings"
)

type SupportController struct {
	l               *log.Logger
	userService     *services.UserService
	peakService     *services.PeakService
	overpassService *services.OverpassService
	activityDao     *daos.ActivityDao
	userPeaksDao    *daos.UserPeaksDao
}

func NewSupportController(
	l *log.Logger,
	userService *services.UserService,
	peakService *services.PeakService,
	overpassService *services.OverpassService,
	activityDao *daos.ActivityDao,
	userPeaksDao *daos.UserPeaksDao,
) *SupportController {
	return &SupportController{
		l:               l,
		userService:     userService,
		peakService:     peakService,
		overpassService: overpassService,
		activityDao:     activityDao,
		userPeaksDao:    userPeaksDao,
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

// RefreshPeaks fetches fresh peak data from OpenStreetMap and optionally
// recalculates summit data for all activities.
// Query params:
//   - recalculate=true: Also clear user_peaks and reset summits_calculated flags
//
// Note: This endpoint is unauthenticated but requires admin_key query param
func (c *SupportController) RefreshPeaks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Simple admin key check (set ADMIN_KEY env var)
	adminKey := r.URL.Query().Get("admin_key")
	expectedKey := os.Getenv("ADMIN_KEY")
	if expectedKey == "" {
		expectedKey = "dev-admin-key" // Default for local development
	}
	if adminKey != expectedKey {
		c.l.Printf("Unauthorized refresh-peaks attempt")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	recalculate := r.URL.Query().Get("recalculate") == "true"

	c.l.Printf("Starting peak data refresh (recalculate=%v)...", recalculate)

	// Step 1: Force fetch fresh peaks from Overpass API (bypasses "already stored" check)
	peaks, err := c.overpassService.ForceFetchPeaks()
	if err != nil {
		c.l.Printf("Error fetching peaks from Overpass: %v", err)
		http.Error(w, "Failed to fetch peaks from OpenStreetMap", http.StatusInternalServerError)
		return
	}

	if peaks == nil || len(peaks.Elements) == 0 {
		c.l.Printf("No peaks returned from Overpass")
		http.Error(w, "No peaks returned from OpenStreetMap", http.StatusInternalServerError)
		return
	}

	// Step 2: Store/update peaks (upsert based on osm_id)
	err = c.peakService.StorePeaks(peaks)
	if err != nil {
		c.l.Printf("Error storing peaks: %v", err)
		http.Error(w, "Failed to store peaks", http.StatusInternalServerError)
		return
	}

	c.l.Printf("Stored/updated %d peaks", len(peaks.Elements))

	result := map[string]interface{}{
		"peaksUpdated": len(peaks.Elements),
	}

	// Step 3 (optional): Recalculate summits
	if recalculate {
		c.l.Printf("Clearing user_peaks and resetting summits_calculated flags...")

		// Clear user_peaks table
		err = c.userPeaksDao.ClearUserPeaks()
		if err != nil {
			c.l.Printf("Error clearing user_peaks: %v", err)
			http.Error(w, "Failed to clear user peaks", http.StatusInternalServerError)
			return
		}

		// Reset summits_calculated flags on all activities
		rowsAffected, err := c.activityDao.ResetSummitsCalculated()
		if err != nil {
			c.l.Printf("Error resetting summits_calculated: %v", err)
			http.Error(w, "Failed to reset summit calculations", http.StatusInternalServerError)
			return
		}

		c.l.Printf("Reset %d activities for summit recalculation", rowsAffected)
		result["activitiesReset"] = rowsAffected
		result["message"] = "Peak data refreshed. User peaks cleared. Activities will recalculate summits on next sync."
	} else {
		result["message"] = "Peak data refreshed successfully"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
