package handlers

import (
	"log"
	"net/http"
	"run-goals/controllers"
)

type ApiHandler struct {
	l                *log.Logger
	apiController    *controllers.ApiController
	groupsController *controllers.GroupsController
}

func NewApiHandler(
	l *log.Logger,
	apiController *controllers.ApiController,
	groupsController *controllers.GroupsController,
) *ApiHandler {
	return &ApiHandler{
		l,
		apiController,
		groupsController,
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
		handler.apiController.ListPeaks(rw, r)
		return
	case "/api/progress":
		handler.apiController.GetProgress(rw, r)
		return
	case "/api/peak-summaries":
		handler.apiController.GetPeakSummaries(rw, r)
		return
	case "/api/groups":
		if r.Method == http.MethodPost {
			handler.groupsController.CreateGroup(rw, r)
			return
		}
	}
}
