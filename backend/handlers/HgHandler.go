package handlers

import (
	"log"
	"net/http"
	"run-goals/controllers"
)

type HgHandler struct {
	l            *log.Logger
	HgController *controllers.HgController
}

func NewHgHandler(
	l *log.Logger,
	hgController *controllers.HgController,
) *HgHandler {
	return &HgHandler{
		l,
		hgController,
	}
}

// ServeHTTP is the main entry point for the handler and satisfies the handler interface
func (handler *HgHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// handle request to get activities
	switch r.URL.Path {
	case "/hg/activities":
		handler.HgController.ListHikeGangActivities(rw, r)
		return
	default:
		http.Error(rw, "Not Found", http.StatusNotFound)
		return
	}
}
