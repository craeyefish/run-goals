package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"run-goals/config"
	"run-goals/daos"
	"run-goals/models"
	"strconv"
	"strings"
	"time"
)

type StravaServiceInterface interface {
	FetchAndStoreUserActivities(user *models.User) error
	EnsureValidToken(u *models.User) error
	GetUserDistance(u *models.User) (*float64, error)
	FetchUserDistance(user *models.User) (float64, error)
	FetchActivitiesPage(accessToken string, page, perPage int)
	ProcessWebhookEvent(payload models.StravaWebhookPayload)
	ProcessCallback(code string) error
}

type StravaService struct {
	l           *log.Logger
	config      *config.Config
	userDao     *daos.UserDao
	activityDao *daos.ActivityDao
}

func NewStravaService(
	l *log.Logger,
	config *config.Config,
	userDao *daos.UserDao,
	activityDao *daos.ActivityDao,
) *StravaService {
	return &StravaService{
		l:           l,
		config:      config,
		userDao:     userDao,
		activityDao: activityDao,
	}
}

func (service *StravaService) FetchAndStoreUserActivities(user *models.User) error {
	// Ensure token is valid first
	if err := service.EnsureValidToken(user); err != nil {
		return fmt.Errorf("token refresh error: %w", err)
	}

	page := 1
	perPage := 30 // Strava default max is 30 or 100 depending on your app scope

	for {
		stravaActivities, err := service.FetchActivitiesPage(user.AccessToken, page, perPage)
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
				StravaAthleteId: user.StravaAthleteID,
				UserID:          user.ID,
				Name:            stravaActivity.Name,
				Distance:        stravaActivity.Distance, // decide if you store in m or km
				StartDate:       t,
				MapPolyline:     stravaActivity.Map.SummaryPolyline,
			}
			if err := service.activityDao.UpsertActivity(&activity); err != nil {
				log.Printf("Error upserting activity %d: %v\n", stravaActivity.ID, err)
			}
		}
		page++
	}

	return nil
}

func (service *StravaService) EnsureValidToken(u *models.User) error {
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

func (service *StravaService) GetUserDistance(u *models.User) (*float64, error) {
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
	dist, err := service.FetchUserDistance(u)
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
func (service *StravaService) FetchUserDistance(user *models.User) (float64, error) {
	if err := service.EnsureValidToken(user); err != nil {
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
func (service *StravaService) FetchActivitiesPage(accessToken string, page, perPage int) ([]models.StravaActivity, error) {
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

func (s *StravaService) ProcessWebhookEvent(payload models.StravaWebhookPayload) {
	// Find the user in DB
	user, err := s.userDao.GetUserByStravaAthleteID(payload.OwnerID)
	if err != nil {
		s.l.Printf("Error calling UserDao.GetUserByStravaAthleteID: %v", err)
		return
	}

	// 'strava_athlete_id' is how we store the Strava athlete ID in the User model
	if user == nil {
		s.l.Printf("No matching  user")
	} else {
		// We found the user - now fetch updated stats
		dist, err := s.FetchUserDistance(user)
		if err != nil {
			s.l.Printf("Error fetching updated distance for user %d: %v\n", user.ID, err)
		} else {
			// Update and save
			user.LastDistance = dist
			user.LastUpdated = time.Now()

			err = s.userDao.UpsertUser(user)
			if err != nil {
				s.l.Printf("Failed to update the user's cached values: %v", err)
			} else {
				s.l.Printf("Updated user %d with new distance: %.2f km\n", user.ID, dist)
			}
		}
	}
}

func (s *StravaService) ProcessCallback(code string) error {
	// 1. Excahnge code for tokens
	tokenRes, err := s.exchangeCodeForToken(code)
	if err != nil {
		s.l.Println("Failed to exchange code", err)
		return err
	}

	// 2. Store (or update) the user in the DB
	user, err := s.userDao.GetUserByStravaAthleteID(tokenRes.Athlete.Id)
	if err != nil {
		s.l.Printf("Error calling UserDao.GetUserByStravaAthleteID: %v", err)
		return err
	}
	// Create new user if not found, update if found
	user.StravaAthleteID = tokenRes.Athlete.Id
	user.AccessToken = tokenRes.AccessToken
	user.RefreshToken = tokenRes.RefreshToken
	user.ExpiresAt = tokenRes.ExpiresAt

	err = s.userDao.UpsertUser(user)
	s.l.Printf("Upsert new user: AthleteID %d", tokenRes.Athlete.Id)

	// 3. Pull activities
	s.FetchAndStoreUserActivities(user)
	return nil
}

func (s *StravaService) exchangeCodeForToken(code string) (*models.StravaTokenResponse, error) {
	formData := url.Values{}
	formData.Set("client_id", s.config.Strava.ClientID)
	formData.Set("client_secret", s.config.Strava.ClientSecret)
	formData.Set("code", code)
	formData.Set("grant_type", "authorization_code")

	resp, err := http.Post(
		"https://www.strava.com/oauth/token",
		"application/x-www-form-urlencoded",
		strings.NewReader(formData.Encode()),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tokenRes models.StravaTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenRes); err != nil {
		return nil, err
	}
	return &tokenRes, nil
}
