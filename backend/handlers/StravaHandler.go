package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"run-goals/config"
	"run-goals/controllers"
)

type StravaHandler struct {
	l          *log.Logger
	controller *controllers.StravaController
}

func NewStravaHandler(l *log.Logger, config *config.Config, db *sql.DB) *StravaHandler {
	return &StravaHandler{
		l,
		controllers.NewStravaController(l, config, db),
	}
}

// ServeHTTP is the main entry point for the handler and satisfies the handler interface
func (handler *StravaHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// handle request to get activities
	switch r.URL.Path {
	case "/webhook/strava":
		handler.controller.ListActivities(rw, r)
		return
	case "/auth/strava/callback":
		handler.controller.ListActivities(rw, r)
		return
	}
}
