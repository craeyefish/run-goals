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
	"run-goals/middleware"
	"run-goals/services"
	"run-goals/workflows"
	"time"
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
	groupsDao := daos.NewGroupsDao(logger, db)

	// initialise services
	jwtService := services.NewJWTService(logger, config)
	stravaService := services.NewStravaService(logger, config, userDao, activityDao)
	activityService := services.NewActivityService(logger, activityDao)
	peakService := services.NewPeakService(logger, peaksDao, userPeaksDao)
	summariesService := services.NewSummariesService(logger, peaksDao, userPeaksDao, activityDao)
	progressService := services.NewProgressService(logger, userDao, stravaService)
	goalProgressService := services.NewGoalProgressService(logger, groupsDao, activityDao, userPeaksDao)
	groupsService := services.NewGroupsService(logger, groupsDao)
	userService := services.NewUserService(logger, userDao)

	// once off data population
	// TODO(cian): Update to sync processes
	summitService := services.NewSummitService(logger, config, peaksDao, userPeaksDao, activityDao)
	summitService.PopulateSummitedPeaks()
	overpassService := services.NewOverpassService(logger, peaksDao)
	peaks, _ := overpassService.FetchPeaks()
	peakService.StorePeaks(peaks)

	// initialise controllers
	apiController := controllers.NewApiController(
		logger,
		activityService,
		progressService,
		peakService,
		summariesService,
		userService,
	)
	authController := controllers.NewAuthController(logger, jwtService)
	groupsController := controllers.NewGroupsController(logger, groupsService, goalProgressService)
	
	// backgorund jobs
	// TODO(cian): Move out of server.
	fetcher := workflows.NewStravaActivityFetcher(stravaService, userDao, activityDao, logger)
	
	hgController := controllers.NewHgController(logger, activityService, userDao, fetcher)
	stravaController := controllers.NewStravaController(logger, jwtService, stravaService)
	supportController := controllers.NewSupportController(logger, userService)

	// initialise handlers
	apiHandler := handlers.NewApiHandler(logger, apiController, groupsController)
	authHandler := handlers.NewAuthHandler(logger, authController, stravaController)
	hgHandler := handlers.NewHgHandler(logger, hgController)
	stravaHandler := handlers.NewStravaHandler(logger, stravaController)
	supportHandler := handlers.NewSupportHandler(logger, supportController)
	// fetcher.FetchUserActivities()
	go func() {
		for {
			time.Sleep(24 * time.Hour)
			fetcher.FetchUserActivities()
		}
	}()

	// create new serve mux and register handlers
	mux := http.NewServeMux()
	mux.Handle("/api/", middleware.JWT(jwtService, apiHandler))
	mux.Handle("/webhook/", stravaHandler)
	mux.Handle("/auth/", authHandler)
	mux.Handle("/hikegang/", hgHandler)
	mux.Handle("/support/", middleware.JWT(jwtService, supportHandler))

	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
}
