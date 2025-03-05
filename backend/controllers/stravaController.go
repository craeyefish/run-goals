package controllers

import (
	"database/sql"
	"log"
	"run-goals/config"
	"run-goals/services"
)

type StravaController struct {
	l       *log.Logger
	service *services.StravaService
}

func NewStravaController(l *log.Logger, config *config.Config, db *sql.DB) *StravaController {
	return &StravaController{
		l:       l,
		service: services.NewStravaService(l, config, db),
	}
}
