package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type TokenRefreshResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresAt    int64  `json:"expires_at"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func ensureValidToken(u *User) error {
	// 1. Check if token is still valid
	if time.Now().Unix() < u.ExpiresAt {
		// Token is not expired yet
		return nil
	}

	// 2. Refresh
	formData := url.Values{}
	formData.Set("client_id", stravaClientID)
	formData.Set("client_secret", stravaClientSecret)
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

	var trr TokenRefreshResponse
	if err := json.NewDecoder(resp.Body).Decode(&trr); err != nil {
		return fmt.Errorf("failed to parse refresh response: %w", err)
	}

	// 3. Update user struct
	u.AccessToken = trr.AccessToken
	u.RefreshToken = trr.RefreshToken
	u.ExpiresAt = trr.ExpiresAt

	// 4. Save to DB
	if err := DB.Save(u).Error; err != nil {
		return fmt.Errorf("failed to save refreshed tokens: %w", err)
	}

	log.Printf("Refreshed token for user with athlete ID %d\n", u.StravaAthleteID)
	return nil
}

const DistanceCacheTTL = 24 * time.Hour // or however long you want to cache

func getUserDistance(u *User) (float64, error) {
	// 1. Check if we have a recent value
	if time.Since(u.LastUpdated) < DistanceCacheTTL {
		// Within cache window, return cached distance
		return u.LastDistance, nil
	}

	// 2. Otherwise, fetch from Strava
	dist, err := fetchUserDistance(*u)
	if err != nil {
		return 0, err
	}

	// 3. Update the user’s cached values in DB
	u.LastDistance = dist
	u.LastUpdated = time.Now()

	if err := DB.Save(u).Error; err != nil {
		return 0, err
	}

	return dist, nil
}

// For demonstration only; use a proper HTTP client and handle errors properly
var httpClient = &http.Client{Timeout: 10 * time.Second}

// Simple function to fetch the total distance for a user from Strava
func fetchUserDistance(user User) (float64, error) {
	if err := ensureValidToken(&user); err != nil {
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

	resp, err := httpClient.Do(req)
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

func fetchAndStoreUserActivities(user *User) error {
	// Ensure token is valid first
	if err := ensureValidToken(user); err != nil {
		return fmt.Errorf("token refresh error: %w", err)
	}

	page := 1
	perPage := 30 // Strava default max is 30 or 100 depending on your app scope

	for {
		activities, err := fetchActivitiesPage(user.AccessToken, page, perPage)
		if err != nil {
			return err
		}
		if len(activities) == 0 {
			break // no more activities
		}

		for _, act := range activities {
			// upsert (create or update) each activity in DB
			if err := upsertActivity(&act, user); err != nil {
				log.Printf("Error upserting activity %d: %v\n", act.ID, err)
			}
		}

		page++
	}
	return nil
}

// fetchActivitiesPage calls the Strava API to fetch a single page of activities
func fetchActivitiesPage(accessToken string, page, perPage int) ([]StravaActivity, error) {
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

	var activities []StravaActivity
	if err := json.NewDecoder(resp.Body).Decode(&activities); err != nil {
		return nil, err
	}

	return activities, nil
}

// StravaActivity is a partial struct to parse activity JSON from Strava
type StravaActivity struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Distance  float64 `json:"distance"` // in meters
	StartDate string  `json:"start_date_local"`
	Map       struct {
		SummaryPolyline string `json:"summary_polyline"`
	} `json:"map"`
	// add other fields if needed
}
