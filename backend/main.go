package main

import (
	"log"
	"net/http"
	"os"
	"run-goals/config"
	"run-goals/database"
	"run-goals/handlers"
)

func main() {
	// create logger
	logger := log.New(os.Stdout, "app", log.LstdFlags)

	// config
	config := config.NewConfig()

	// database setup
	db := database.OpenPG(config, logger)

	// create handlers
	apiHandler := handlers.NewApiHandler(logger, db)
	// webhookHandler
	// authHandler

	// create new serve mux and register handlers
	mux := http.NewServeMux()
	mux.Handle("/api/", apiHandler)
	// mux.Handle("/webhook/", webhookHandler)
	// mux.Handle("/auth/", authHandler)

	InitDB()

	http.HandleFunc("/webhook/strava", handleStravaWebhookEvents)
	http.HandleFunc("/api/progress", handleProgress)
	http.HandleFunc("/auth/strava/callback", handleStravaCallback)
	http.HandleFunc("/api/activities", handleListActivities)
	http.HandleFunc("/api/peaks", handleListPeaks)
	http.HandleFunc("/api/peak-summaries", handlePeakSummaries)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// TODO(cian):
// 1. Add activity syncing
//  - pull activity on webhook event.
//  - sync all activities when user joins.
//  - maybe sync past 24 hours of activities once a day (not sure if we need a catching loop?)
// 2. Peaks outside Western Cape ? Less peaks ?
// 3. Check summited peaks in activity (workflow?) - PopulateSummitedPeaks()
// 4. Make a group setup process.
