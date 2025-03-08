package server

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"run-goals/config"
	"run-goals/controllers"
	"run-goals/database"
	"run-goals/handlers"
	"run-goals/services"
)

type Server struct {
	Mux *http.ServeMux
}

func NewServer() (*http.Server, error) {
	// create logger
	logger := log.New(os.Stdout, "app", log.LstdFlags)

	// config
	config := config.NewConfig()

	// database setup
	db := database.OpenPG(config, logger)

	// initialise services
	activityService := services.NewActivityService(logger, db)
	stravaService := services.NewStravaService(logger, config, db)
	peakService := services.NewPeakService(logger, db)
	progressService := services.NewProgressService(logger, db, stravaService)

	// initialise controllers
	apiController := controllers.NewApiController(l*log.Logger, db*sql.DB)
	stravaController := controllers.NewStravaController(l*log.Logger, config*config.Config, db*sql.DB)

	// initialise handlers
	apiHandler := handlers.NewApiHandler(logger, db)
	stravaHandler := handlers.NewStravaHandler(logger, config, db)

	// create new serve mux and register handlers
	mux := http.NewServeMux()
	mux.Handle("/api/", apiHandler)
	mux.Handle("/webhook/", stravaHandler)
	mux.Handle("/auth/", stravaHandler)

	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}, nil
}
