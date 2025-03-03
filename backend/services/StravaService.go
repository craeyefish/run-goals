package services

import (
	"database/sql"
	"log"
	"net/http"
	"net/url"
	"strings"
	"strconv"
	"fmt"
	"time"
	"encoding/json"
	"run-goals/daos"
	"run-goals/models"
	"run-goals/config"
)

type StravaService struct {
	l   		*log.Logger
	config 		*config.Config
	userDao		*daos.UserDao
	activityDao *daos.ActivityDao
}

func NewStravaService(l *log.Logger, config *config.Config, db *sql.DB) *StravaService {
	userDao := daos.NewUserDao(l, db)
	activityDao := daos.NewActivityDao(l, db)
	return &StravaService{
		l:   	     	l,
		config: 		config,
		userDao: 		userDao,
		activityDao: 	activityDao,
	}
}

func (service *StravaService) FetchAndStoreUserActivities(user *models.User) error {
	// Ensure token is valid first
	if err := service.ensureValidToken(user); err != nil {
		return fmt.Errorf("token refresh error: %w", err)
	}

	page := 1
	perPage := 30 // Strava default max is 30 or 100 depending on your app scope

	for {
		stravaActivities, err := service.fetchActivitiesPage(user.AccessToken, page, perPage)
		if err != nil {
			return err
		}
		if len(stravaActivities) == 0 {
			break // no more activities
		}
		for _, stravaActivity := range stravaActivities {
			// upsert (create or update) each activity in DB
			// first convert stravaActivity into our activity model

			t, _ := time.Parse(time.RFC3339, stravaActivity.StartDate) // handle error properly
			activity := models.Activity{
				StravaActivityID: stravaActivity.ID,
				UserID:           user.ID,
				Name:             stravaActivity.Name,
				Distance:         stravaActivity.Distance, // decide if you store in m or km
				StartDate:        t,
				MapPolyline:      stravaActivity.Map.SummaryPolyline,
			}
			if err := service.activityDao.UpsertActivity(&activity); err != nil {
				log.Printf("Error upserting activity %d: %v\n", stravaActivity.ID, err)
			}
		}
		page++
	}

	return nil
}

func (service *StravaService) ensureValidToken(u *models.User) error {
	// 1. Check if token is still valid
	if time.Now().Unix() < u.ExpiresAt {
		// Token is not expired yet
		return nil
	}

	// 2. Refresh
	formData := url.Values{}
	formData.Set("client_id", service.config.Strava.ClientID)
	formData.Set("client_secret", service.config.Strava.ClientSecret)
	formData.Set("grant_type", "refresh_token")
	formData.Set("refresh_token", u.RefreshToken)

	resp, err := http.Post(
		"https://www.strava.com/oauth/token",
		"application/x-www-form-urlencoded",
		strings.NewReader(formData.Encode()),
	)
	if err != nil {
		return fmt.Errorf("failed to refresh token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("token refresh request failed with status %d", resp.StatusCode)
	}

	var trr models.TokenRefreshResponse
	if err := json.NewDecoder(resp.Body).Decode(&trr); err != nil {
		return fmt.Errorf("failed to parse refresh response: %w", err)
	}

	// 3. Update user struct
	u.AccessToken = trr.AccessToken
	u.RefreshToken = trr.RefreshToken
	u.ExpiresAt = trr.ExpiresAt

	// 4. Save to DB
	err = service.userDao.UpsertUser(u)
	if err != nil {
		return fmt.Errorf("failed to save refreshed tokens: %w", err)
	}

	log.Printf("Refreshed token for user with athlete ID %d\n", u.StravaAthleteID)
	return nil
}

func (service *StravaService) getUserDistance(u *models.User) (*float64, error) {
	// 1. Check if we have a recent value
	distanceCacheTTL, err := strconv.ParseInt(service.config.Strava.DistanceCacheTTL, 10, 64)
	if err != nil {
		fmt.Println("Error converting to int:", err)
		return nil, err
	}
	if time.Since(u.LastUpdated) < (time.Hour * time.Duration(distanceCacheTTL)) {
		// Within cache window, return cached distance
		return &u.LastDistance, nil
	}

	// 2. Otherwise, fetch from Strava
	dist, err := service.fetchUserDistance(u)
	dist = 0
	if err != nil {
		return &dist, err
	}

	// 3. Update the user’s cached values in DB
	u.LastDistance = dist
	u.LastUpdated = time.Now()

	err = service.userDao.UpsertUser(u)
	if err != nil {
		return nil, fmt.Errorf("failed to update the user's cached values: %w", err)
	}

	return &dist, nil
}

// For demonstration only; use a proper HTTP client and handle errors properly
// var httpClient = &http.Client{Timeout: 10 * time.Second}

// Simple function to fetch the total distance for a user from Strava
func (service *StravaService) fetchUserDistance(user *models.User) (float64, error) {
	if err := service.ensureValidToken(user); err != nil {
		return 0, err
	}

	// Endpoint for Strava activities: https://www.strava.com/api/v3/athlete/activities
	// Or the athlete/stats endpoint: https://www.strava.com/api/v3/athletes/{id}/stats
	//
	// For demonstration, let's pretend we call athlete/stats
	// (Need user’s athlete ID as well, so adjust logic accordingly)
	// You might also sum up distances from activities for the current year.

	// Pseudo-code (replace {athleteId} with actual ID):
	// GET https://www.strava.com/api/v3/athletes/{athleteId}/stats
	// Authorization: Bearer {user.AccessToken}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.strava.com/api/v3/athletes/%d/stats", user.StravaAthleteID), nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", "Bearer "+user.AccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// For the example, let's parse a minimal subset of the JSON
	type athleteStats struct {
		YtdRunTotals struct {
			Distance float64 `json:"distance"`
		} `json:"ytd_run_totals"`
	}
	var stats athleteStats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return 0, err
	}

	// Strava distance is in meters; convert to km
	distanceInKm := stats.YtdRunTotals.Distance / 1000.0
	return distanceInKm, nil
}

// fetchActivitiesPage calls the Strava API to fetch a single page of activities
func (service *StravaService) fetchActivitiesPage(accessToken string, page, perPage int) ([]models.StravaActivity, error) {
	url := fmt.Sprintf("https://www.strava.com/api/v3/athlete/activities?page=%d&per_page=%d", page, perPage)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch activities status %d", resp.StatusCode)
	}

	var activities []models.StravaActivity
	if err := json.NewDecoder(resp.Body).Decode(&activities); err != nil {
		return nil, err
	}

	return activities, nil
}
