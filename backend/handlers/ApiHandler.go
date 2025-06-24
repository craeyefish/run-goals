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
	case "/api/profile":
		handler.apiController.GetUserProfile(rw, r)
		return
	case "/api/peak-summaries":
		handler.apiController.GetPeakSummaries(rw, r)
		return
	case "/api/groups":
		if r.Method == http.MethodPost {
			handler.groupsController.CreateGroup(rw, r)
			return
		}
		if r.Method == http.MethodPut {
			handler.groupsController.UpdateGroup(rw, r)
			return
		}
		if r.Method == http.MethodDelete {
			handler.groupsController.DeleteGroup(rw, r)
			return
		}
		if r.Method == http.MethodGet {
			handler.groupsController.GetUserGroups(rw, r)
			return
		}
	case "/api/group-member":
		if r.Method == http.MethodPost {
			handler.groupsController.CreateGroupMember(rw, r)
			return
		}
		if r.Method == http.MethodPut {
			handler.groupsController.UpdateGroupMember(rw, r)
			return
		}
		if r.Method == http.MethodDelete {
			handler.groupsController.DeleteGroupMember(rw, r)
			return
		}
	case "/api/group-members":
		if r.Method == http.MethodGet {
			handler.groupsController.GetGroupMembers(rw, r)
			return
		}
	case "/api/group-members-contribution":
		if r.Method == http.MethodGet {
			handler.groupsController.GetGroupMembersGoalContribution(rw, r)
			return
		}
	case "/api/group-goal":
		if r.Method == http.MethodPost {
			handler.groupsController.CreateGroupGoal(rw, r)
			return
		}
		if r.Method == http.MethodPut {
			handler.groupsController.UpdateGroupGoal(rw, r)
			return
		}
		if r.Method == http.MethodDelete {
			handler.groupsController.DeleteGroupGoal(rw, r)
			return
		}
	case "/api/group-goals":
		if r.Method == http.MethodGet {
			handler.groupsController.GetGroupGoals(rw, r)
			return
		}
	}
}
