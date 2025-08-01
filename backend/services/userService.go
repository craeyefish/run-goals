package services

import (
	"log"
	"run-goals/daos"
	"run-goals/models"
)

type UserServiceInterface interface {
	GetUserByID(userID int64) (*models.User, error)
	GetUserProfile(userID int64) (*models.User, error)
	DeleteUserAccount(stravaAthleteID int64) error
}

type UserService struct {
	l       *log.Logger
	userDao *daos.UserDao
}

func NewUserService(
	l *log.Logger,
	userDao *daos.UserDao,
) *UserService {
	return &UserService{
		l:       l,
		userDao: userDao,
	}
}

func (s *UserService) GetUserByID(userID int64) (*models.User, error) {
	user, err := s.userDao.GetUserByID(userID)
	if err != nil {
		s.l.Printf("Error calling userDao.GetUserByID: %v", err)
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetUserProfile(userID int64) (*models.User, error) {
	user, err := s.userDao.GetUserByID(userID)
	if err != nil {
		s.l.Printf("Error getting user profile for user %d: %v", userID, err)
		return nil, err
	}

	// Clear sensitive fields before returning
	user.AccessToken = ""
	user.RefreshToken = ""

	return user, nil
}

func (s *UserService) DeleteUserAccount(stravaAthleteID int64) error {
	err := s.userDao.DeleteUserByStravaAthleteID(stravaAthleteID)
	if err != nil {
		s.l.Printf("Error deleting user account for strava_athlete_id %d: %v", stravaAthleteID, err)
		return err
	}
	s.l.Printf("Successfully deleted user account for strava_athlete_id %d", stravaAthleteID)
	return nil
}
