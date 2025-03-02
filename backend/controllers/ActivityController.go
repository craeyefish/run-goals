package controllers

import (
	"log"
	"database/sql"
	"run-goals/services"
)

type ActivityController struct {
	l       *log.Logger
	service *services.ActivityService
}

func NewActivityController(l *log.Logger, db *sql.DB) *ActivityController {
	return &ActivityController{
		l:       l,
		service: services.NewActivityService(l, db),
	}
}
