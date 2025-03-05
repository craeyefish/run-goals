package controllers

import (
	"database/sql"
	"log"
	"run-goals/services"
)

type UserController struct {
	l       *log.Logger
	service *services.ActivityService
}

func NewUserController(l *log.Logger, db *sql.DB) *UserController {
	return &UserController{
		l:       l,
		service: services.NewUserService(l, db),
	}
}
