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
	personalYearlyGoalDao := daos.NewPersonalYearlyGoalDao(logger, db)
	summitFavouritesDao := daos.NewSummitFavouritesDao(logger, db)
	challengeDao := daos.NewChallengeDao(logger, db)

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
	personalGoalsService := services.NewPersonalGoalsService(logger, personalYearlyGoalDao)
	summitFavouritesService := services.NewSummitFavouritesService(logger, summitFavouritesDao)
	challengeService := services.NewChallengeService(logger, challengeDao)

	// Services for background jobs
	summitService := services.NewSummitService(logger, config, peaksDao, userPeaksDao, activityDao)
	overpassService := services.NewOverpassService(logger, peaksDao)

	// One-time peak data fetch on startup (peaks don't change often)
	go func() {
		logger.Println("Background: Starting peak data fetch from OpenStreetMap...")
		peaks, err := overpassService.FetchPeaks()
		if err != nil {
			logger.Printf("Background: Failed to fetch peaks: %v", err)
		} else if peaks != nil {
			peakService.StorePeaks(peaks)
			logger.Printf("Background: Stored %d peaks", len(peaks.Elements))
		}
	}()

	// initialise controllers
	apiController := controllers.NewApiController(
		logger,
		activityService,
		progressService,
		peakService,
		summariesService,
		userService,
		personalGoalsService,
		summitFavouritesService,
	)
	authController := controllers.NewAuthController(logger, jwtService)
	groupsController := controllers.NewGroupsController(logger, groupsService, goalProgressService)
	challengesController := controllers.NewChallengesController(logger, challengeService)

	// background jobs
	// TODO(cian): Move out of server.
	fetcher := workflows.NewStravaActivityFetcher(stravaService, summitService, userDao, activityDao, logger)

	hgController := controllers.NewHgController(logger, activityService, userDao, fetcher)
	stravaController := controllers.NewStravaController(logger, jwtService, stravaService, summitService, activityDao)
	supportController := controllers.NewSupportController(logger, userService, peakService, overpassService, activityDao, userPeaksDao)

	// initialise handlers
	apiHandler := handlers.NewApiHandler(logger, apiController, groupsController, challengesController)
	authHandler := handlers.NewAuthHandler(logger, authController, stravaController)
	hgHandler := handlers.NewHgHandler(logger, hgController)
	stravaHandler := handlers.NewStravaHandler(logger, stravaController)
	supportHandler := handlers.NewSupportHandler(logger, supportController)

	// background sync job - disabled via DISABLE_SYNC_JOB=true for local development
	if os.Getenv("DISABLE_SYNC_JOB") != "true" {
		// Daily sync - fetches only recent activities (last 30 days)
		go func() {
			logger.Println("Starting initial recent activity sync...")
			fetcher.FetchRecentUserActivities()
			logger.Println("Initial recent sync complete. Daily sync scheduled.")
			for {
				time.Sleep(24 * time.Hour)
				fetcher.FetchRecentUserActivities()
			}
		}()

		// Weekly full sync - fetches all activities (runs on Sundays)
		go func() {
			for {
				// Sleep until next Sunday at 3am
				now := time.Now()
				daysUntilSunday := (7 - int(now.Weekday())) % 7
				if daysUntilSunday == 0 && now.Hour() >= 3 {
					daysUntilSunday = 7 // Already past Sunday 3am, wait for next week
				}
				nextSunday := time.Date(now.Year(), now.Month(), now.Day()+daysUntilSunday, 3, 0, 0, 0, now.Location())
				sleepDuration := time.Until(nextSunday)
				logger.Printf("Weekly full sync scheduled for %s (in %s)", nextSunday.Format(time.RFC1123), sleepDuration.Round(time.Hour))
				time.Sleep(sleepDuration)
				logger.Println("Starting weekly full activity sync...")
				fetcher.FetchUserActivities()
				logger.Println("Weekly full sync complete.")
			}
		}()
	} else {
		logger.Println("Sync job disabled via DISABLE_SYNC_JOB environment variable")
	}

	// create new serve mux and register handlers
	mux := http.NewServeMux()
	mux.Handle("/api/", middleware.JWT(jwtService, apiHandler))
	mux.Handle("/webhook/", stravaHandler)
	mux.Handle("/auth/", authHandler)
	mux.Handle("/hikegang/", hgHandler)
	mux.Handle("/support/", middleware.JWT(jwtService, supportHandler))
	// Admin endpoints - no JWT, uses admin_key query param
	mux.HandleFunc("/admin/refresh-peaks", supportController.RefreshPeaks)

	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
}
