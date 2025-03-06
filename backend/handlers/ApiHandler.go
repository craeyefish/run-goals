package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"run-goals/controllers"
)

type ApiHandler struct {
	l                  *log.Logger
	activityController *controllers.ActivityController
	userController     *controllers.UserController
}

func NewApiHandler(l *log.Logger, db *sql.DB) *ApiHandler {
	return &ApiHandler{
		l,
		controllers.NewActivityController(l, db),
		controllers.NewUserController(l, db),
	}
}

// ServeHTTP is the main entry point for the handler and satisfies the handler interface
func (handler *ApiHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// handle request to get activities
	switch r.URL.Path {
	case "/api/activities":
		handler.activityController.ListActivities(rw, r)
		return
	case "/api/progress":
		handler.progressController.
		return
	}

}
