package server

import (
	"log"
	"net/http"
	"os"
	"run-goals/config"
	"run-goals/controllers"
	"run-goals/daos"
	"run-goals/database"
	"run-goals/handlers"
	"run-goals/services"
)

type Server struct {
	Mux *http.ServeMux
}

func NewServer() *http.Server {
	// create logger
	logger := log.New(os.Stdout, "app", log.LstdFlags)

	// config
	config := config.NewConfig()

	// database setup
	db := database.OpenPG(config, logger)

	// intialise daos
	activityDao := daos.NewActivityDao(logger, db)
	peaksDao := daos.NewPeaksDao(logger, db)
	userDao := daos.NewUserDao(logger, db)
	userPeaksDao := daos.NewUserPeaksDao(logger, db)

	// initialise services
	stravaService := services.NewStravaService(logger, config, userDao, activityDao)
	activityService := services.NewActivityService(logger, activityDao)
	peakService := services.NewPeakService(logger, peaksDao, userPeaksDao)
	summariesService := services.NewSummariesService(logger, peaksDao, userPeaksDao, activityDao)
	progressService := services.NewProgressService(logger, userDao, stravaService)

	// initialise controllers
	apiController := controllers.NewApiController(
		logger,
		activityService,
		progressService,
		peakService,
		summariesService,
	)
	stravaController := controllers.NewStravaController(logger, stravaService)

	// initialise handlers
	apiHandler := handlers.NewApiHandler(logger, apiController)
	stravaHandler := handlers.NewStravaHandler(logger, stravaController)

	// create new serve mux and register handlers
	mux := http.NewServeMux()
	mux.Handle("/api/", apiHandler)
	mux.Handle("/webhook/", stravaHandler)
	mux.Handle("/auth/", stravaHandler)

	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
}
