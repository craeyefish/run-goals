package handlers

import (
	"database/sql"
	"log"
	"run-goals/controllers"
)


type ActivityHandler struct {
	l  *log.Logger
	controller *controllers.ActivityController
}

func NewActivityHandler(l *log.Logger, db *sql.DB) *ActivityHandler {
	return &ActivityHandler{
		l,
		controllers.NewActivityController(l, db),
	}
}
