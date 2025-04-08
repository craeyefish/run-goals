package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"run-goals/dto"
	"run-goals/services"
	"strconv"
)

type GroupsControllerInterface interface {
	CreateGroup(rw http.ResponseWriter, r *http.Request)
	UpdateGroup(rw http.ResponseWriter, r *http.Request)
	DeleteGroup(rw http.ResponseWriter, r *http.Request)
	GetUserGroups(rw http.ResponseWriter, r *http.Request)

	CreateGroupMember(rw http.ResponseWriter, r *http.Request)
	UpdateGroupMember(rw http.ResponseWriter, r *http.Request)
	DeleteGroupMember(rw http.ResponseWriter, r *http.Request)

	CreateGroupGoal(rw http.ResponseWriter, r *http.Request)
	UpdateGroupGoal(rw http.ResponseWriter, r *http.Request)
	DeleteGroupGoal(rw http.ResponseWriter, r *http.Request)

	// todo: get goal progress
}

type GroupsController struct {
	l             *log.Logger
	groupsService *services.GroupsService
}

func NewGroupsController(
	l *log.Logger,
	groupsService *services.GroupsService,
) *GroupsController {
	return &GroupsController{
		l:             l,
		groupsService: groupsService,
	}
}

func (c *GroupsController) CreateGroup(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle POST groups - creating new group")

	// read and decode the json body
	var request dto.CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		c.l.Printf("Error unmarshalling data: %v", err)
		http.Error(rw, "Error unmarshalling data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err := c.groupsService.CreateGroup(request)
	if err != nil {
		c.l.Printf("Error creating group: %v", err)
		http.Error(rw, "Failed to create group", http.StatusInternalServerError)
		return
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

	err := c.groupsService.UpdateGroup(request)
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

	// read and decode the json body
	var request dto.CreateGroupMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		c.l.Printf("Error unmarshalling data: %v", err)
		http.Error(rw, "Error unmarshalling data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err := c.groupsService.CreateGroupMember(request)
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

	err := c.groupsService.UpdateGroupMember(request)
	if err != nil {
		c.l.Printf("Error updating group member: %v", err)
		http.Error(rw, "Failed to update group member", http.StatusInternalServerError)
		return
	}
}

func (c *GroupsController) DeleteGroupMember(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle DELETE group-member - deleting existing group member")

	// extract userID from url
	strUserID := r.URL.Query().Get("userID")
	if strUserID == "" {
		http.Error(rw, "missing userID", http.StatusBadRequest)
		return
	}
	userID, err := strconv.ParseInt(strUserID, 10, 64)
	if err != nil {
		http.Error(rw, "Invalid user ID", http.StatusBadRequest)
		return
	}
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

	err := c.groupsService.CreateGroupGoal(request)
	if err != nil {
		c.l.Printf("Error creating group goal: %v", err)
		http.Error(rw, "Failed to create group goal", http.StatusInternalServerError)
		return
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

	// extract userID from url
	strID := r.URL.Query().Get("userID")
	if strID == "" {
		http.Error(rw, "missing userID", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		http.Error(rw, "Invalid user ID", http.StatusBadRequest)
		return
	}

	groups, err := c.groupsService.GetUserGroups(id)
	response := dto.GetUserGroupsResponse{
		groups: groups,
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
