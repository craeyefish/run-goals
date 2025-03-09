package main

import (
	"log"
	"net/http"
)

func main() {
	InitDB()

	http.HandleFunc("/webhook/strava", handleStravaWebhookEvents)
	http.HandleFunc("/api/progress", handleProgress)
	http.HandleFunc("/auth/strava/callback", handleStravaCallback)
	http.HandleFunc("/api/activities", handleListActivities)
	http.HandleFunc("/api/peaks", handleListPeaks)
	http.HandleFunc("/api/peak-summaries", handlePeakSummaries)

	// TODO(cian): Move these to async processes.
	PopulateSummitedPeaks()
	FetchPeaks()

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

	// TODO(steven): refactor
	// refactor start

	// // create server
	// server := server.NewServer()
	// // start the server
	// go func() {
	// 	log.Println("Starting server on http://localhost:8080")
	// 	err := server.ListenAndServe()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }()
	// // trap sigterm or interrupt and gracefully shutdown the server
	// // create channel 'sigChan'
	// sigChan := make(chan os.Signal, 1)
	// // register channel to receive notifications for both interrupt and kill signals
	// signal.Notify(sigChan, os.Interrupt)
	// signal.Notify(sigChan, os.Kill)
	// // wait until signal is received and log
	// sig := <-sigChan
	// log.Println("Received terminate, graceful shutdown", sig)
	// // create context used to allow server time to finish processing ongoing requests before shutting down
	// ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	// // shutdown server with created context - gracefully shutting down
	// server.Shutdown(ctx)

	// refactor end
}

// TODO(cian):
//
// 1. Add activity syncing
//  - pull activity on webhook event.
//  - sync all activities when user joins.
//  - maybe sync past 24 hours of activities once a day (not sure if we need a catching loop?)
// 2. Peaks outside Western Cape ? Less peaks ?
// 3. Check summited peaks in activity (workflow?) - PopulateSummitedPeaks()
// 4. Make a group setup process.
