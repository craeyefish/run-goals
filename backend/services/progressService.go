package services

import (
	"database/sql"
	"log"
	"run-goals/daos"
	"run-goals/models"
	"sync"
)

type ProgressServiceInterface interface {
	GetUsersProgress() (*models.GoalProgress, error)
}

type ProgressService struct {
	l             *log.Logger
	userDao       *daos.UserDao
	stravaService *StravaService
}

func NewProgressService(
	l *log.Logger,
	db *sql.DB,
	stravaService *StravaService,
) *ProgressService {
	userDao := daos.NewUserDao(l, db)
	return &ProgressService{
		l:             l,
		userDao:       userDao,
		stravaService: stravaService,
	}
}

// TODO(cian): Make group setup process and store this in db.
const groupGoal = 1000.0

func (s *ProgressService) GetUsersProgress() (*models.GoalProgress, error) {
	// load all users from db
	users, err := s.userDao.GetUsers()
	if err != nil {
		s.l.Printf("Error calling UserDao: %v", err)
		return nil, err
	}

	var total float64
	contributions := make([]models.UserContribution, 0, len(users))

	var wg sync.WaitGroup
	var mu sync.Mutex

	wg.Add(len(users))
	for _, user := range users {
		go func(u models.User) {
			defer wg.Done()

			dist, err := s.stravaService.GetUserDistance(&u)
			if err != nil {
				log.Println("Error fetching distance for user:", u.ID, err)
				return
			}

			s.stravaService.FetchAndStoreUserActivities(&u)
			if err != nil {
				log.Println("Error fetching activities for user:", u.ID, err)
				return
			}

			mu.Lock()
			total += *dist
			contributions = append(contributions, models.UserContribution{
				ID:            u.StravaAthleteID,
				TotalDistance: *dist,
			})
			mu.Unlock()
		}(user)
	}

	wg.Wait()

	// Construct response
	goalProgress := models.GoalProgress{
		Goal:            groupGoal,
		CurrentProgress: total,
		Contributions:   contributions,
	}

	return &goalProgress, nil
}
