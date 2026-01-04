package handlers

import (
	"log"
	"net/http"
	"run-goals/controllers"
)

type ApiHandler struct {
	l                    *log.Logger
	apiController        *controllers.ApiController
	groupsController     *controllers.GroupsController
	challengesController *controllers.ChallengesController
}

func NewApiHandler(
	l *log.Logger,
	apiController *controllers.ApiController,
	groupsController *controllers.GroupsController,
	challengesController *controllers.ChallengesController,
) *ApiHandler {
	return &ApiHandler{
		l,
		apiController,
		groupsController,
		challengesController,
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
		if r.Method == http.MethodPut || r.Method == http.MethodPatch {
			handler.apiController.UpdateUserProfile(rw, r)
			return
		}
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
	case "/api/group-goal-progress":
		if r.Method == http.MethodGet {
			handler.groupsController.GetGroupGoalProgress(rw, r)
			return
		}
	case "/api/personal-goals":
		if r.Method == http.MethodGet {
			handler.apiController.GetPersonalGoals(rw, r)
			return
		}
		if r.Method == http.MethodPost {
			handler.apiController.SavePersonalGoals(rw, r)
			return
		}
	case "/api/personal-goals/all":
		if r.Method == http.MethodGet {
			handler.apiController.GetAllPersonalGoals(rw, r)
			return
		}
	case "/api/summit-favourites":
		if r.Method == http.MethodGet {
			handler.apiController.GetSummitFavourites(rw, r)
			return
		}
		if r.Method == http.MethodPost {
			handler.apiController.AddSummitFavourite(rw, r)
			return
		}
		if r.Method == http.MethodDelete {
			handler.apiController.RemoveSummitFavourite(rw, r)
			return
		}

	// ==================== Challenge Routes ====================
	case "/api/challenges":
		if r.Method == http.MethodPost {
			handler.challengesController.CreateChallenge(rw, r)
			return
		}
		if r.Method == http.MethodGet {
			handler.challengesController.GetUserChallenges(rw, r)
			return
		}
	case "/api/challenge":
		if r.Method == http.MethodGet {
			handler.challengesController.GetChallenge(rw, r)
			return
		}
		if r.Method == http.MethodPut {
			handler.challengesController.UpdateChallenge(rw, r)
			return
		}
		if r.Method == http.MethodDelete {
			handler.challengesController.DeleteChallenge(rw, r)
			return
		}
	case "/api/challenges/featured":
		if r.Method == http.MethodGet {
			handler.challengesController.GetFeaturedChallenges(rw, r)
			return
		}
	case "/api/challenges/public":
		if r.Method == http.MethodGet {
			handler.challengesController.GetPublicChallenges(rw, r)
			return
		}
	case "/api/challenges/search":
		if r.Method == http.MethodGet {
			handler.challengesController.SearchChallenges(rw, r)
			return
		}
	case "/api/challenge-peaks":
		if r.Method == http.MethodGet {
			handler.challengesController.GetChallengePeaks(rw, r)
			return
		}
		if r.Method == http.MethodPut {
			handler.challengesController.SetChallengePeaks(rw, r)
			return
		}
	case "/api/challenge-join":
		if r.Method == http.MethodPost {
			handler.challengesController.JoinChallenge(rw, r)
			return
		}
	case "/api/challenge-join-by-code":
		if r.Method == http.MethodPost {
			handler.challengesController.JoinChallengeByCode(rw, r)
			return
		}
	case "/api/challenge-leave":
		if r.Method == http.MethodDelete {
			handler.challengesController.LeaveChallenge(rw, r)
			return
		}
	case "/api/challenge-lock":
		if r.Method == http.MethodPost {
			handler.challengesController.LockChallenge(rw, r)
			return
		}
	case "/api/challenge-participants":
		if r.Method == http.MethodGet {
			handler.challengesController.GetParticipants(rw, r)
			return
		}
	case "/api/challenge-leaderboard":
		if r.Method == http.MethodGet {
			handler.challengesController.GetLeaderboard(rw, r)
			return
		}
	case "/api/challenge-summit-log":
		if r.Method == http.MethodGet {
			handler.challengesController.GetSummitLog(rw, r)
			return
		}
	case "/api/challenge-activities":
		if r.Method == http.MethodGet {
			handler.challengesController.GetChallengeActivities(rw, r)
			return
		}
	case "/api/challenge-summit":
		if r.Method == http.MethodPost {
			handler.challengesController.RecordSummit(rw, r)
			return
		}
	case "/api/challenge-group":
		if r.Method == http.MethodPost {
			handler.challengesController.AddGroupToChallenge(rw, r)
			return
		}
		if r.Method == http.MethodDelete {
			handler.challengesController.RemoveGroupFromChallenge(rw, r)
			return
		}
	case "/api/group-challenges":
		if r.Method == http.MethodGet {
			handler.challengesController.GetGroupChallenges(rw, r)
			return
		}
	}
}
