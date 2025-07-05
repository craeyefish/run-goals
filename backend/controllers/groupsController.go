package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"run-goals/dto"
	"run-goals/meta"
	"run-goals/services"
	"strconv"
	"time"
)

type GroupsControllerInterface interface {
	CreateGroup(rw http.ResponseWriter, r *http.Request)
	UpdateGroup(rw http.ResponseWriter, r *http.Request)
	DeleteGroup(rw http.ResponseWriter, r *http.Request)
	GetUserGroups(rw http.ResponseWriter, r *http.Request)

	CreateGroupMember(rw http.ResponseWriter, r *http.Request)
	UpdateGroupMember(rw http.ResponseWriter, r *http.Request)
	DeleteGroupMember(rw http.ResponseWriter, r *http.Request)
	GetGroupMembers(rw http.ResponseWriter, r *http.Request)
	GetGroupMembersGoalContribution(rw http.ResponseWriter, r *http.Request)

	CreateGroupGoal(rw http.ResponseWriter, r *http.Request)
	UpdateGroupGoal(rw http.ResponseWriter, r *http.Request)
	DeleteGroupGoal(rw http.ResponseWriter, r *http.Request)
	GetGroupGoals(rw http.ResponseWriter, r *http.Request)
	GetGroupGoalProgress(rw http.ResponseWriter, r *http.Request)
}

type GroupsController struct {
	l                   *log.Logger
	groupsService       *services.GroupsService
	goalProgressService *services.GoalProgressService
}

func NewGroupsController(
	l *log.Logger,
	groupsService *services.GroupsService,
	goalProgressService *services.GoalProgressService,
) *GroupsController {
	return &GroupsController{
		l:                   l,
		groupsService:       groupsService,
		goalProgressService: goalProgressService,
	}
}

func (c *GroupsController) CreateGroup(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle POST groups - creating new group")

	userID, _ := meta.GetUserIDFromContext(r.Context())

	// read and decode the json body
	var request dto.CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		c.l.Printf("Error unmarshalling data: %v", err)
		http.Error(rw, "Error unmarshalling data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	groupID, err := c.groupsService.CreateGroup(request.Name, userID)
	if err != nil {
		c.l.Printf("Error creating group: %v", err)
		http.Error(rw, "Failed to create group", http.StatusInternalServerError)
		return
	}

	response := dto.CreateGroupResponse{
		GroupID: *groupID,
	}

	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(response); err != nil {
		log.Println("Error encoding CreateGroupResponse:", err)
	}
}

func (c *GroupsController) UpdateGroup(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle PUT groups - updating existing group")

	// read and decode the json body
	var request dto.UpdateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		c.l.Printf("Error unmarshalling data: %v", err)
		http.Error(rw, "Error unmarshalling data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err := c.groupsService.UpdateGroup(request.ID, request.Name)
	if err != nil {
		c.l.Printf("Error updating group: %v", err)
		http.Error(rw, "Failed to update group", http.StatusInternalServerError)
		return
	}
}

func (c *GroupsController) DeleteGroup(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle DELETE groups - deleting existing group")

	// extract group id from url
	strID := r.URL.Query().Get("groupID")
	if strID == "" {
		http.Error(rw, "missing groupID", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		http.Error(rw, "Invalid group ID", http.StatusBadRequest)
		return
	}

	err = c.groupsService.DeleteGroup(id)
	if err != nil {
		c.l.Printf("Error deleting group: %v", err)
		http.Error(rw, "Failed to delete group", http.StatusInternalServerError)
		return
	}
}

func (c *GroupsController) CreateGroupMember(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle POST group-member - creating new group member")

	userID, _ := meta.GetUserIDFromContext(r.Context())

	// read and decode the json body
	var request dto.CreateGroupMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		c.l.Printf("Error unmarshalling data: %v", err)
		http.Error(rw, "Error unmarshalling data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err := c.groupsService.CreateGroupMember(request.GroupCode, userID, request.Role)
	if err != nil {
		c.l.Printf("Error creating group member: %v", err)
		http.Error(rw, "Failed to create group member", http.StatusInternalServerError)
		return
	}
}

func (c *GroupsController) UpdateGroupMember(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle PUT group-member - updating existing group member")

	// read and decode the json body
	var request dto.UpdateGroupMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		c.l.Printf("Error unmarshalling data: %v", err)
		http.Error(rw, "Error unmarshalling data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err := c.groupsService.UpdateGroupMember(request.GroupID, request.UserID, request.Role)
	if err != nil {
		c.l.Printf("Error updating group member: %v", err)
		http.Error(rw, "Failed to update group member", http.StatusInternalServerError)
		return
	}
}

func (c *GroupsController) DeleteGroupMember(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle DELETE group-member - deleting existing group member")

	userID, _ := meta.GetUserIDFromContext(r.Context())

	// extract groupID from url
	strGroupID := r.URL.Query().Get("groupID")
	if strGroupID == "" {
		http.Error(rw, "missing groupID", http.StatusBadRequest)
		return
	}
	groupID, err := strconv.ParseInt(strGroupID, 10, 64)
	if err != nil {
		http.Error(rw, "Invalid group ID", http.StatusBadRequest)
		return
	}

	err = c.groupsService.DeleteGroupMember(userID, groupID)
	if err != nil {
		c.l.Printf("Error deleting group member: %v", err)
		http.Error(rw, "Failed to delete group member", http.StatusInternalServerError)
		return
	}
}

func (c *GroupsController) CreateGroupGoal(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle POST group-goal - creating new group goal")

	// read and decode the json body
	var request dto.CreateGroupGoalRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		c.l.Printf("Error unmarshalling data: %v", err)
		http.Error(rw, "Error unmarshalling data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	goalID, err := c.groupsService.CreateGroupGoal(request)
	if err != nil {
		c.l.Printf("Error creating group goal: %v", err)
		http.Error(rw, "Failed to create group goal", http.StatusInternalServerError)
		return
	}

	response := dto.CreateGroupGoalResponse{
		GoalID: *goalID,
	}

	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(response); err != nil {
		log.Println("Error encoding CreateGroupGoalResponse:", err)
	}
}

func (c *GroupsController) UpdateGroupGoal(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle PUT group-goal - updating existing group goal")

	// read and decode the json body
	var request dto.UpdateGroupGoalRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		c.l.Printf("Error unmarshalling data: %v", err)
		http.Error(rw, "Error unmarshalling data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err := c.groupsService.UpdateGroupGoal(request)
	if err != nil {
		c.l.Printf("Error updating group goal %v", err)
		http.Error(rw, "Failed to update group goal", http.StatusInternalServerError)
		return
	}
}

func (c *GroupsController) DeleteGroupGoal(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle DELETE group-goal - deleting existing group goal")

	// extract goal id from url
	strID := r.URL.Query().Get("goalID")
	if strID == "" {
		http.Error(rw, "missing goalID", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		http.Error(rw, "Invalid group ID", http.StatusBadRequest)
		return
	}

	err = c.groupsService.DeleteGroupGoal(id)
	if err != nil {
		c.l.Printf("Error deleting group goal: %v", err)
		http.Error(rw, "Failed to delete group goal", http.StatusInternalServerError)
		return
	}
}

func (c *GroupsController) GetUserGroups(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle GET groups - get user groups")

	userID, _ := meta.GetUserIDFromContext(r.Context())

	groups, err := c.groupsService.GetUserGroups(userID)
	response := dto.GetUserGroupsResponse{
		Groups: groups,
	}
	if err != nil {
		c.l.Printf("Error getting user groups %v", err)
		http.Error(rw, "Failed to get user groups", http.StatusInternalServerError)
		return
	}

	// Return the array of Groups as JSON
	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(response); err != nil {
		log.Println("Error encoding groups:", err)
	}
}

func (c *GroupsController) GetGroupGoals(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle GET group-goals - get group goals")

	// extract groupID from url
	strID := r.URL.Query().Get("groupID")
	if strID == "" {
		http.Error(rw, "missing groupID", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		http.Error(rw, "Invalid group ID", http.StatusBadRequest)
		return
	}
	goals, err := c.groupsService.GetGroupGoals(id)
	response := dto.GetGroupGoalsResponse{
		Goals: goals,
	}
	if err != nil {
		c.l.Printf("Error getting user groups %v", err)
		http.Error(rw, "Failed to get user groups", http.StatusInternalServerError)
		return
	}

	// Return the array of Goals as JSON
	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(response); err != nil {
		log.Println("Error encoding groups:", err)
	}
}

func (c *GroupsController) GetGroupMembersGoalContribution(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle GET group-members-contribution - get group member goal contributions")

	// extract groupID from url
	str := r.URL.Query().Get("groupID")
	if str == "" {
		http.Error(rw, "missing groupID", http.StatusBadRequest)
		return
	}
	groupID, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		http.Error(rw, "Invalid group ID", http.StatusBadRequest)
		return
	}

	// extract startDate from url
	str = r.URL.Query().Get("startDate")
	if str == "" {
		http.Error(rw, "missing startDate", http.StatusBadRequest)
		return
	}
	startDate, err := time.Parse("2006-01-02", str)
	if err != nil {
		http.Error(rw, "Invalid start date", http.StatusBadRequest)
		return
	}

	// extract endDate from url
	str = r.URL.Query().Get("endDate")
	if str == "" {
		http.Error(rw, "missing endDate", http.StatusBadRequest)
		return
	}
	endDate, err := time.Parse("2006-01-02", str)
	if err != nil {
		http.Error(rw, "Invalid start date", http.StatusBadRequest)
		return
	}

	groupMembersGoalContribution, err := c.groupsService.GetGroupMembersGoalContribution(groupID, startDate, endDate)
	response := dto.GetGroupMembersGoalContributionResponse{
		Members: groupMembersGoalContribution,
	}
	if err != nil {
		c.l.Printf("Error getting group members goal contribution %v", err)
		http.Error(rw, "Failed to get group members goal contribution", http.StatusInternalServerError)
		return
	}

	// Return the array of Goals as JSON
	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(response); err != nil {
		log.Println("Error encoding groups:", err)
	}
}

func (c *GroupsController) GetGroupMembers(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle GET group-members - get group members")

	// extract groupID from url
	str := r.URL.Query().Get("groupID")
	if str == "" {
		http.Error(rw, "missing groupID", http.StatusBadRequest)
		return
	}
	groupID, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		http.Error(rw, "Invalid group ID", http.StatusBadRequest)
		return
	}

	groupMembers, err := c.groupsService.GetGroupMembers(groupID)
	response := dto.GetGroupMembersResponse{
		Members: groupMembers,
	}
	if err != nil {
		c.l.Printf("Error getting group members %v", err)
		http.Error(rw, "Failed to get group members", http.StatusInternalServerError)
		return
	}

	// Return the array of Goals as JSON
	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(response); err != nil {
		log.Println("Error encoding groups:", err)
	}
}

func (c *GroupsController) GetGroupGoalProgress(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle GET group-goal-progress - get individual goal progress")

	// extract goalID from url
	str := r.URL.Query().Get("goalID")
	if str == "" {
		http.Error(rw, "missing goalID", http.StatusBadRequest)
		return
	}
	goalID, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		http.Error(rw, "Invalid goal ID", http.StatusBadRequest)
		return
	}

	// Get the goal directly by ID
	targetGoal, err := c.groupsService.GetGroupGoalByID(goalID)
	if err != nil {
		c.l.Printf("Error getting goal by ID: %v", err)
		http.Error(rw, "Failed to get goal", http.StatusInternalServerError)
		return
	}

	if targetGoal == nil {
		http.Error(rw, "Goal not found", http.StatusNotFound)
		return
	}

	// Calculate progress
	progress, err := c.goalProgressService.CalculateGoalProgress(*targetGoal)
	if err != nil {
		c.l.Printf("Error calculating goal progress: %v", err)
		http.Error(rw, "Failed to calculate goal progress", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"goalID":   goalID,
		"progress": progress,
	}

	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(response); err != nil {
		log.Println("Error encoding goal progress response:", err)
	}
}
