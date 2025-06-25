package services

import (
	"fmt"
	"log"
	"run-goals/daos"
	"run-goals/models"
)

type GoalProgressService struct {
	l             *log.Logger
	groupsDao     *daos.GroupsDao
	activitiesDao *daos.ActivityDao
	userPeaksDao  *daos.UserPeaksDao
}

func NewGoalProgressService(
	l *log.Logger,
	groupsDao *daos.GroupsDao,
	activitiesDao *daos.ActivityDao,
	userPeaksDao *daos.UserPeaksDao,
) *GoalProgressService {
	return &GoalProgressService{
		l:             l,
		groupsDao:     groupsDao,
		activitiesDao: activitiesDao,
		userPeaksDao:  userPeaksDao,
	}
}

func (s *GoalProgressService) CalculateGoalProgress(goal models.GroupGoal) (float64, error) {
	switch goal.GoalType {
	case "distance":
		return s.calculateDistanceProgress(goal)
	case "elevation":
		return s.calculateElevationProgress(goal)
	case "summit_count":
		return s.calculateSummitCountProgress(goal)
	case "specific_summits":
		return s.calculateSpecificSummitsProgress(goal)
	default:
		return 0, fmt.Errorf("unknown goal type: %s", goal.GoalType)
	}
}

func (s *GoalProgressService) calculateDistanceProgress(goal models.GroupGoal) (float64, error) {
	// Get group members
	members, err := s.groupsDao.GetGroupMembers(goal.GroupID)
	if err != nil {
		return 0, err
	}

	var totalDistance float64
	for _, member := range members {
		memberContribution, err := s.groupsDao.GetGroupMembersGoalContribution(
			goal.GroupID,
			goal.StartDate,
			goal.EndDate,
		)
		if err != nil {
			return 0, err
		}

		for _, contribution := range memberContribution {
			if contribution.UserID == member.UserID {
				totalDistance += contribution.TotalDistance / 1000 // Convert meters to km
				break
			}
		}
	}

	if goal.TargetValue == 0 {
		return 0, nil
	}

	progress := (totalDistance / goal.TargetValue) * 100
	if progress > 100 {
		progress = 100
	}

	return progress, nil
}

func (s *GoalProgressService) calculateElevationProgress(goal models.GroupGoal) (float64, error) {
	// Similar to distance but sum elevation gain from activities
	// You'll need to add elevation tracking to your GroupMemberGoalContribution
	// For now, return 0 as placeholder
	return 0, nil
}

func (s *GoalProgressService) calculateSpecificSummitsProgress(goal models.GroupGoal) (float64, error) {
	// Count how many of the specific target summits have been completed by group members
	if len(goal.TargetSummits) == 0 {
		return 0, nil
	}

	members, err := s.groupsDao.GetGroupMembers(goal.GroupID)
	if err != nil {
		return 0, err
	}

	// Track which summits have been completed
	completedSummits := make(map[int64]bool)

	for _, member := range members {
		userSummits, err := s.userPeaksDao.GetUserSummitsInDateRange(
			member.UserID,
			goal.TargetSummits,
			goal.StartDate,
			goal.EndDate,
		)
		if err != nil {
			return 0, err
		}

		for _, summit := range userSummits {
			completedSummits[summit.PeakID] = true
		}
	}

	progress := (float64(len(completedSummits)) / float64(len(goal.TargetSummits))) * 100
	if progress > 100 {
		progress = 100
	}

	return progress, nil
}

func (s *GoalProgressService) calculateSummitCountProgress(goal models.GroupGoal) (float64, error) {
	members, err := s.groupsDao.GetGroupMembers(goal.GroupID)
	if err != nil {
		return 0, err
	}

	var totalSummits int64
	for _, member := range members {
		// Get all summits for this user in the date range
		userSummits, err := s.userPeaksDao.GetUserSummitsInDateRangeAll(
			member.UserID,
			goal.StartDate,
			goal.EndDate,
		)
		if err != nil {
			return 0, err
		}

		totalSummits += int64(len(userSummits))
	}

	if goal.TargetValue == 0 {
		return 0, nil
	}

	progress := (float64(totalSummits) / goal.TargetValue) * 100
	if progress > 100 {
		progress = 100
	}

	return progress, nil
}
