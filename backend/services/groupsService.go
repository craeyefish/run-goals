package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"run-goals/daos"
	"run-goals/dto"
	"run-goals/models"
	"time"
)

type GroupServiceInterface interface {
	CreateGroup(request dto.CreateGroupRequest) (int64, error)
	UpdateGroup(group dto.UpdateGroupRequest) error
	DeleteGroup(groupID int64) error

	CreateGroupMember(member dto.CreateGroupMemberRequest) error
	UpdateGroupMember(member dto.UpdateGroupMemberRequest) error
	DeleteGroupMember(userID int64, groupID int64) error
	GetGroupMembers(groupID int64) ([]models.GroupMember, error)
	GetGroupMembersGoalContribution(groupID int64, startDate time.Time, endDate time.Time) ([]models.GroupMemberGoalContribution, error)

	CreateGroupGoal(goal dto.CreateGroupGoalRequest) (int64, error)
	UpdateGroupGoal(goal dto.UpdateGroupGoalRequest) error
	DeleteGroupGoal(goalID int64) error

	GetUserGroups(userID int64) ([]models.Group, error)
	GetGroupGoals(groupID int64) ([]models.GroupGoal, error)
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

func (s *GroupsService) CreateGroup(request dto.CreateGroupRequest) (*int64, error) {
	code, err := s.GenerateGroupCode()
	if err != nil {
		return nil, err
	}
	group := models.Group{
		ID:        0,
		Name:      request.Name,
		Code:      code,
		CreatedBy: request.CreatedBy,
		CreatedAt: time.Now(),
	}
	groupID, err := s.groupsDao.CreateGroup(group)
	if err != nil {
		s.l.Printf("Error calling groupsDao.CreateGroup: %v", err)
		return nil, err
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
		return nil, err
	}

	return groupID, nil
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
	id, err := s.groupsDao.GetGroupIDFromCode(request.GroupCode)
	if err != nil {
		s.l.Printf("Error calling groupsDao.GetGroupIDFromCode: %v", err)
		return err
	}
	groupMember := models.GroupMember{
		ID:       0,
		GroupID:  *id,
		UserID:   request.UserID,
		Role:     request.Role,
		JoinedAt: time.Now(),
	}
	err = s.groupsDao.CreateGroupMember(groupMember)
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

func (s *GroupsService) CreateGroupGoal(request dto.CreateGroupGoalRequest) (*int64, error) {
	goal := models.GroupGoal{
		ID:          0,
		GroupID:     request.GroupID,
		Name:        request.Name,
		TargetValue: request.TargetValue,
		StartDate:   request.StartDate,
		EndDate:     request.EndDate,
		CreatedAt:   time.Now(),
	}
	goalID, err := s.groupsDao.CreateGroupGoal(goal)
	if err != nil {
		s.l.Printf("Error calling groupsDao.CreateGroupGoal: %v", err)
		return nil, err
	}
	return goalID, nil
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

func (s *GroupsService) GetGroupGoals(groupID int64) ([]models.GroupGoal, error) {
	groupGoals, err := s.groupsDao.GetGroupGoals(groupID)
	if err != nil {
		s.l.Printf("Error calling groupsDao.GetGroupGoals: %v", err)
		return nil, err
	}
	return groupGoals, nil
}

func (s *GroupsService) GetGroupMembersGoalContribution(groupID int64, startDate time.Time, endDate time.Time) ([]models.GroupMemberGoalContribution, error) {
	groupMembersContribution, err := s.groupsDao.GetGroupMembersGoalContribution(groupID, startDate, endDate)
	if err != nil {
		s.l.Printf("Error calling groupsDao.GetGroupMembersGoalContribution: %v", err)
		return nil, err
	}
	return groupMembersContribution, nil
}

func (s *GroupsService) GetGroupMembers(groupID int64) ([]models.GroupMember, error) {
	groupMembers, err := s.groupsDao.GetGroupMembers(groupID)
	if err != nil {
		s.l.Printf("Error calling groupsDao.GetGroupMembers: %v", err)
		return nil, err
	}
	return groupMembers, nil
}

func (s *GroupsService) GenerateGroupCode() (string, error) {
	maxAttempts := 5
	for i := 0; i < maxAttempts; i++ {
		bytes := make([]byte, 3)
		_, err := rand.Read(bytes)
		if err != nil {
			return "", err
		}
		code := hex.EncodeToString(bytes)
		exists, err := s.CheckGroupCodeExists(code)
		if err != nil {
			return "", err
		}
		if !exists {
			return code, nil
		}
	}
	return "", fmt.Errorf("Failed to generate unique group code after reaching max attemps: %v", maxAttempts)
}

func (s *GroupsService) CheckGroupCodeExists(code string) (bool, error) {
	count, err := s.groupsDao.CheckGroupCodeExists(code)
	if err != nil {
		s.l.Printf("Error calling groupsDao.CheckGroupCodeExists: %v", err)
		return false, err
	}
	return *count > 0, nil
}
