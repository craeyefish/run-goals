package services

import (
	"database/sql"
	"log"
	"run-goals/daos"
	"run-goals/models"
)

type ActivityServiceInterface interface {
	GetActivitiesByUserID(userID int64) ([]models.Activity, error)
	UpsertActivitiesByUserId(id int, activity *models.Activity) error
}

type ActivityService struct {
	l           *log.Logger
	activityDao *daos.ActivityDao
}

func NewActivityService(l *log.Logger, db *sql.DB) *ActivityService {
	activityDao := daos.NewActivityDao(l, db)
	return &ActivityService{
		l:           l,
		activityDao: activityDao,
	}
}

func (s *ActivityService) GetActivitiesByUserID(userID int64) ([]models.Activity, error) {
	activities, err := s.activityDao.GetActivitiesByUserID(userID)
	if err != nil {
		s.l.Printf("Error calling ActivityDao: %v", err)
		return nil, err
	}
	return activities, nil
}

func (s *ActivityService) UpsertActivitiesByUserId(id int, activity *models.Activity) error {
	err := s.activityDao.UpsertActivity(activity)
	if err != nil {
		s.l.Printf("Error calling ActivityDao: %v", err)
		return err
	}
	return nil
}
