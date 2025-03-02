package services

import (
	"database/sql"
	"log"
	"run-goals/daos"
	"run-goals/models"
)

type ActivityService struct {
	l   *log.Logger
	dao *daos.ActivityDao
}

func NewActivityService(l *log.Logger, db *sql.DB) *ActivityService {
	customerDao := daos.NewActivityDao(l, db)
	return &ActivityService{
		l:   l,
		dao: customerDao,
	}
}

// func (service *ActivityService) UpsertActivity(id int, activity *models.Activity) error {
// 	err := service.dao.UpsertActivity(activity)
// 	if err != nil {
// 		service.l.Printf("Error calling ActivityDao: %v", err)
// 		return err
// 	}
// 	return nil
// }

func (service *ActivityService) GetActivitiesByUserId(id int, activity *models.Activity) error {
	err := service.dao.UpsertActivity(activity)
	if err != nil {
		service.l.Printf("Error calling ActivityDao: %v", err)
		return err
	}
	return nil
}
