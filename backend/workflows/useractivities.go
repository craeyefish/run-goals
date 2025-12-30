package workflows

import (
	"log"
	"time"

	"run-goals/daos"
	"run-goals/services"
)

type StravaActivityFetcher struct {
	stravaService *services.StravaService
	summitService *services.SummitService
	activitiesDao *daos.ActivityDao
	userDao       *daos.UserDao
	logger        *log.Logger
}

// NewStravaActivityFetcher initializes the fetcher.
func NewStravaActivityFetcher(
	stravaService *services.StravaService,
	summitService *services.SummitService,
	usersDao *daos.UserDao,
	activitiesDao *daos.ActivityDao,
	logger *log.Logger,
) *StravaActivityFetcher {
	return &StravaActivityFetcher{
		stravaService: stravaService,
		summitService: summitService,
		activitiesDao: activitiesDao,
		userDao:       usersDao,
		logger:        logger,
	}
}

// FetchUserActivities pulls ALL user activities from Strava (full sync).
// Use this for weekly syncs or initial data population.
func (s *StravaActivityFetcher) FetchUserActivities() {
	s.logger.Println("Starting FULL sync of user activities from Strava...")
	s.fetchActivities(false)
	s.logger.Println("Finished FULL sync of user activities from Strava.")
}

// FetchRecentUserActivities pulls only recent (last 30 days) activities from Strava.
// Use this for daily syncs to minimize API calls.
func (s *StravaActivityFetcher) FetchRecentUserActivities() {
	s.logger.Println("Starting RECENT sync of user activities from Strava (last 30 days)...")
	s.fetchActivities(true)
	s.logger.Println("Finished RECENT sync of user activities from Strava.")
}

// fetchActivities is the internal method that handles both sync modes
func (s *StravaActivityFetcher) fetchActivities(recentOnly bool) {
	users, err := s.userDao.GetUsers()
	if err != nil {
		s.logger.Printf("Error fetching users: %v", err)
		return
	}

	for _, user := range users {
		var fetchErr error
		if recentOnly {
			fetchErr = s.stravaService.FetchAndStoreRecentUserActivities(&user)
		} else {
			fetchErr = s.stravaService.FetchAndStoreUserActivities(&user)
		}
		if fetchErr != nil {
			s.logger.Printf("Error fetching activities for user %d: %v", user.ID, fetchErr)
			continue
		}

		// Fetch detailed data only for #hg activities that don't have it yet
		activities, err := s.activitiesDao.GetActivitiesByUserID(user.ID)
		if err != nil {
			s.logger.Printf("Error fetching activities from database for user %d: %v", user.ID, err)
			continue
		}

		for _, activity := range activities {
			if !activity.IsHG() {
				continue
			}

			// Skip if we already have detailed data (has moving_time and elevation from detailed endpoint)
			// We check MovingTime > 0 because the list endpoint doesn't include these fields
			if activity.MovingTime > 0 && activity.Elevation > 0 {
				continue
			}

			s.logger.Printf("Fetching detailed activity for user %d, activity %d", user.ID, activity.StravaActivityId)
			s.stravaService.FetchAndStoreDetailedActivity(&user, activity.StravaActivityId)
			time.Sleep(100 * time.Millisecond)
		}
	}

	// After syncing activities, run summit detection for any new activities
	s.logger.Println("Running summit detection for new activities...")
	if err := s.summitService.PopulateSummitedPeaks(); err != nil {
		s.logger.Printf("Error during summit detection: %v", err)
	} else {
		s.logger.Println("Summit detection complete")
	}
}
