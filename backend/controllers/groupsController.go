package controllers

import (
	"log"
	"net/http"
	"run-goals/models"
	"run-goals/services"
	"strconv"
)

type GroupsControllerInterface interface {
	CreateGroup(rw http.ResponseWriter, r *http.Request)
	UpdateGroup(rw http.ResponseWriter, r *http.Request)
	DeleteGroup(rw http.ResponseWriter, r *http.Request)
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
	var group models.Group
	err := group.FromJSON(r.Body)
	if err != nil {
		c.l.Printf("Error unmarshalling data: %v", err)
		http.Error(rw, "Error unmarshalling data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	c.groupsService.CreateGroup(group)
}

func (c *GroupsController) UpdateGroup(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle PUT groups - updating existing group")

	// read and decode the json body
	var group models.Group
	err := group.FromJSON(r.Body)
	if err != nil {
		c.l.Printf("Error unmarshalling data: %v", err)
		http.Error(rw, "Error unmarshalling data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	c.groupsService.UpdateGroup(group)
}

func (c *GroupsController) DeleteGroup(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle DELETE groups - deleting existing group")

	// extract group id from url
	idStr := r.URL.Path[len("/api/groups"):]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(rw, "Invalid group ID", http.StatusBadRequest)
		return
	}

	c.groupsService.DeleteGroup(id)
}
