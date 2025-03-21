package services

import (
	"log"
	"run-goals/daos"
	"run-goals/models"
)

type GroupServiceInterface interface {
	CreateGroup(group models.Group) error
	UpdateGroup(group models.Group) error
	DeleteGroup(groupID int64) error

	// CreateGroupMember(member models.GroupMember) error
	// UpdateGroupMember(member models.GroupMember) error
	// DeleteGroupMember(userID int64) error

	// CreateGroupGoal(goal models.GroupGoal) error
	// UpdateGroupGoal(goal models.GroupGoal) error
	// DeleteGroupGoal(goalID int64) error

	// GetUserGroups(userID int64) ([]models.Group, error)
	// GetGroupMembers(groupID int64) ([]models.GroupMember, error)
}

type GroupsService struct {
	l         *log.Logger
	groupsDao *daos.GroupsDao
}

func NewGroupsService(
	l *log.Logger,
	groupsDao *daos.GroupsDao,
) *GroupsService {
	return &GroupsService{
		l:         l,
		groupsDao: groupsDao,
	}
}

func (s *GroupsService) CreateGroup(group models.Group) error {
	err := s.groupsDao.CreateGroup(group)
	if err != nil {
		s.l.Printf("Error calling groupsDao.CreateGroup: %v", err)
		return err
	}
	return nil
}

func (s *GroupsService) UpdateGroup(group models.Group) error {
	err := s.groupsDao.UpdateGroup(group)
	if err != nil {
		s.l.Printf("Error calling groupsDao.UpdateGroup: %v", err)
		return err
	}
	return nil
}

func (s *GroupsService) DeleteGroup(groupID int64) error {
	err := s.groupsDao.DeleteGroup(groupID)
	if err != nil {
		s.l.Printf("Error calling groupsDao.DeleteGroup: %v", err)
		return err
	}
	return nil
}
