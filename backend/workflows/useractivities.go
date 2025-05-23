package workflows

import (
	"log"
	"time"

	"run-goals/daos"
	"run-goals/services"
)

type StravaActivityFetcher struct {
	stravaService *services.StravaService
	activitiesDao *daos.ActivityDao
	userDao       *daos.UserDao
	logger        *log.Logger
}

// NewStravaActivityFetcher initializes the fetcher.
func NewStravaActivityFetcher(stravaService *services.StravaService, usersDao *daos.UserDao, activitiesDao *daos.ActivityDao, logger *log.Logger) *StravaActivityFetcher {
	return &StravaActivityFetcher{
		stravaService: stravaService,
		activitiesDao: activitiesDao,
		userDao:       usersDao,
		logger:        logger,
	}
}

// FetchUserActivities pulls user activities from Strava.
func (s *StravaActivityFetcher) FetchUserActivities() {
	s.logger.Println("Starting to fetch user activities from Strava...")

	users, err := s.userDao.GetUsers()
	if err != nil {
		s.logger.Printf("Error fetching users: %v", err)
		return
	}

	for _, user := range users {
		err = s.stravaService.FetchAndStoreUserActivities(&user)
		if err != nil {
			s.logger.Printf("Error fetching activities for user %d: %v", user.ID, err)
			continue
		}

		// Fetch all of the user's activities from the database
		activities, err := s.activitiesDao.GetActivitiesByUserID(user.ID)
		if err != nil {
			s.logger.Printf("Error fetching activities from database for user %d: %v", user.ID, err)
			continue
		}

		for _, activity := range activities {
			if !activity.IsHG() {
				continue
			}

			s.logger.Printf("Fetching detailed activity for user %d, activity %d", user.ID, activity.StravaActivityId)
			s.stravaService.FetchAndStoreDetailedActivity(&user, activity.StravaActivityId)
			time.Sleep(100 * time.Millisecond)
		}
	}

	s.logger.Println("Finished fetching user activities from Strava.")
}
