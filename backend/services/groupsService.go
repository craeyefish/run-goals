package services

import (
	"log"
	"run-goals/daos"
	"run-goals/dto"
	"run-goals/models"
	"time"
)

type GroupServiceInterface interface {
	CreateGroup(request dto.CreateGroupRequest) error
	UpdateGroup(group dto.UpdateGroupRequest) error
	DeleteGroup(groupID int64) error

	CreateGroupMember(member dto.CreateGroupMemberRequest) error
	UpdateGroupMember(member dto.UpdateGroupMemberRequest) error
	DeleteGroupMember(userID int64, groupID int64) error

	CreateGroupGoal(goal dto.CreateGroupGoalRequest) error
	UpdateGroupGoal(goal dto.UpdateGroupGoalRequest) error
	DeleteGroupGoal(goalID int64) error

	GetUserGroups(userID int64) ([]models.Group, error)
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

func (s *GroupsService) CreateGroup(request dto.CreateGroupRequest) error {
	group := models.Group{
		ID:        0,
		Name:      request.Name,
		CreatedBy: request.CreatedBy,
		CreatedAt: time.Now(),
	}
	groupID, err := s.groupsDao.CreateGroup(group)
	if err != nil {
		s.l.Printf("Error calling groupsDao.CreateGroup: %v", err)
		return err
	}

	groupMember := models.GroupMember{
		GroupID:  *groupID,
		UserID:   group.CreatedBy,
		Role:     "admin",
		JoinedAt: group.CreatedAt,
	}
	err = s.groupsDao.CreateGroupMember(groupMember)
	if err != nil {
		s.l.Printf("Error calling groupsDao.CreateGroupMember: %v", err)
		return err
	}

	return nil
}

func (s *GroupsService) UpdateGroup(request dto.UpdateGroupRequest) error {
	group := models.Group{
		ID:        request.ID,
		Name:      request.Name,
		CreatedBy: request.CreatedBy,
		CreatedAt: request.CreatedAt,
	}
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

func (s *GroupsService) CreateGroupMember(request dto.CreateGroupMemberRequest) error {
	groupMember := models.GroupMember{
		ID:       0,
		GroupID:  request.GroupID,
		UserID:   request.UserID,
		Role:     request.Role,
		JoinedAt: time.Now(),
	}
	err := s.groupsDao.CreateGroupMember(groupMember)
	if err != nil {
		s.l.Printf("Error calling groupsDao.CreateMember: %v", err)
		return err
	}
	return nil
}

func (s *GroupsService) UpdateGroupMember(request dto.UpdateGroupMemberRequest) error {
	member := models.GroupMember{
		ID:       0,
		GroupID:  request.GroupID,
		UserID:   request.UserID,
		Role:     request.Role,
		JoinedAt: request.JoinedAt,
	}
	err := s.groupsDao.UpdateGroupMember(member)
	if err != nil {
		s.l.Printf("Error calling groupsDao.UpdateGroupMember: %v", err)
		return err
	}
	return nil
}

func (s *GroupsService) DeleteGroupMember(userID int64, groupID int64) error {
	err := s.groupsDao.DeleteGroupMember(userID, groupID)
	if err != nil {
		s.l.Printf("Error calling groupsDao.DeleteGroupMember: %v", err)
		return err
	}
	return nil
}

func (s *GroupsService) CreateGroupGoal(request dto.CreateGroupGoalRequest) error {
	goal := models.GroupGoal{
		ID:          0,
		GroupID:     request.GroupID,
		Name:        request.Name,
		TargetValue: request.TargetValue,
		StartDate:   request.StartDate,
		EndDate:     request.EndDate,
		CreatedAt:   time.Now(),
	}
	err := s.groupsDao.CreateGroupGoal(goal)
	if err != nil {
		s.l.Printf("Error calling groupsDao.CreateGroupGoal: %v", err)
		return err
	}
	return nil
}

func (s *GroupsService) UpdateGroupGoal(request dto.UpdateGroupGoalRequest) error {
	goal := models.GroupGoal{
		ID:          request.ID,
		GroupID:     request.GroupID,
		Name:        request.Name,
		TargetValue: request.TargetValue,
		StartDate:   request.StartDate,
		EndDate:     request.EndDate,
		CreatedAt:   time.Now(),
	}
	err := s.groupsDao.UpdateGroupGoal(goal)
	if err != nil {
		s.l.Printf("Error calling groupsDao.UpdateGroupGoal: %v", err)
		return err
	}
	return nil
}

func (s *GroupsService) DeleteGroupGoal(goalID int64) error {
	err := s.groupsDao.DeleteGroupGoal(goalID)
	if err != nil {
		s.l.Printf("Error calling groupsDao.DeleteGroupGoal: %v", err)
		return err
	}
	return nil
}

func (s *GroupsService) GetUserGroups(userID int64) ([]models.Group, error) {
	userGroups, err := s.groupsDao.GetUserGroups(userID)
	if err != nil {
		s.l.Printf("Error calling groupsDao.GetUserGroups: %v", err)
		return nil, err
	}
	return userGroups, nil
}
