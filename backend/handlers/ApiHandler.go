package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"run-goals/controllers"
)

type ApiHandler struct {
	l             *log.Logger
	apiController *controllers.ApiController
}

func NewApiHandler(l *log.Logger, db *sql.DB) *ApiHandler {
	return &ApiHandler{
		l,
		controllers.NewApiController(l, db),
	}
}

// ServeHTTP is the main entry point for the handler and satisfies the handler interface
func (handler *ApiHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// handle request to get activities
	switch r.URL.Path {
	case "/api/activities":
		handler.apiController.ListActivities(rw, r)
		return
	case "/api/peaks":
		handler.apiController.ListActivities(rw, r)
		return
	case "/api/progress":
		handler.apiController.ListActivities(rw, r)
		return
	case "/api/peak-summaries":
		handler.apiController.ListActivities(rw, r)
		return
	}
}
