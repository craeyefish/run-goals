package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func init() {
	// Load .env file if it exists; ignore error in production
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Ensure the required env vars are set
	if os.Getenv("STRAVA_CLIENT_ID") == "" || os.Getenv("STRAVA_CLIENT_SECRET") == "" {
		log.Fatal("STRAVA_CLIENT_ID and STRAVA_CLIENT_SECRET must be set in your environment")
	}
}

var (
	stravaClientID     = os.Getenv("STRAVA_CLIENT_ID")
	stravaClientSecret = os.Getenv("STRAVA_CLIENT_SECRET")
)

type TokenResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresAt    int64  `json:"expires_at"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	Athlete      struct {
		Id int64 `json:"id"`
		// ... other fields if needed
	} `json:"athlete"`
}

func handleStravaCallback(w http.ResponseWriter, r *http.Request) {
	// 1. Get the authorization code from query params
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "missing code", http.StatusBadRequest)
		return
	}

	// 2. Exchange code for tokens
	tokenRes, err := exchangeCodeForToken(code)
	if err != nil {
		http.Error(w, "failed to exchange code", http.StatusInternalServerError)
		return
	}

	// 3. Store (or update) the user in the DB
	err = storeUserToken(tokenRes)
	if err != nil {
		http.Error(w, "failed to store user", http.StatusInternalServerError)
		return
	}

	// 4. Pull activities
	// TODO(cian): Make async. (workflow?)
	var user User
	result := DB.Where("strava_athlete_id = ?", tokenRes.Athlete.Id).First(&user)
	if result.Error != nil {
		http.Error(w, "failed to lookup user", http.StatusInternalServerError)
		return
	}
	fetchAndStoreUserActivities(&user)

	// 4. Redirect back to frontend or show success
	http.Redirect(w, r, "https://craeyebytes.com/", http.StatusFound)
}

// exchangeCodeForToken calls Strava's OAuth token endpoint
func exchangeCodeForToken(code string) (*TokenResponse, error) {
	formData := url.Values{}
	formData.Set("client_id", stravaClientID)
	formData.Set("client_secret", stravaClientSecret)
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

	var tokenRes TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenRes); err != nil {
		return nil, err
	}
	return &tokenRes, nil
}

// storeUserToken upserts a user (create if not found, update if found)
func storeUserToken(tokenRes *TokenResponse) error {
	// Find existing user by athlete ID
	var user User
	result := DB.Where("strava_athlete_id = ?", tokenRes.Athlete.Id).First(&user)
	if result.Error != nil {
		// If user not found, create new record
		if result.Error.Error() == "record not found" {
			user = User{
				StravaAthleteID: tokenRes.Athlete.Id,
				AccessToken:     tokenRes.AccessToken,
				RefreshToken:    tokenRes.RefreshToken,
				ExpiresAt:       tokenRes.ExpiresAt,
			}
			if err := DB.Create(&user).Error; err != nil {
				return err
			}
			log.Printf("Created new user: AthleteID %d", tokenRes.Athlete.Id)
			return nil
		}
		// If any other error, return it
		return result.Error
	}

	// If found, update tokens
	user.AccessToken = tokenRes.AccessToken
	user.RefreshToken = tokenRes.RefreshToken
	user.ExpiresAt = tokenRes.ExpiresAt
	if err := DB.Save(&user).Error; err != nil {
		return err
	}

	log.Printf("Updated user tokens: AthleteID %d", tokenRes.Athlete.Id)
	return nil
}
