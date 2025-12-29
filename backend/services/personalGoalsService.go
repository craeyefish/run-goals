package services

import (
	"log"
	"run-goals/daos"
	"run-goals/models"
	"time"
)

type PersonalGoalsService struct {
	l   *log.Logger
	dao *daos.PersonalYearlyGoalDao
}

func NewPersonalGoalsService(l *log.Logger, dao *daos.PersonalYearlyGoalDao) *PersonalGoalsService {
	return &PersonalGoalsService{
		l:   l,
		dao: dao,
	}
}

// GetCurrentYearGoal gets the user's goal for the current year, or creates a default one
func (s *PersonalGoalsService) GetCurrentYearGoal(userID int64) (*models.PersonalYearlyGoal, error) {
	currentYear := time.Now().Year()
	return s.GetGoalForYear(userID, currentYear)
}

// GetGoalForYear gets the user's goal for a specific year
func (s *PersonalGoalsService) GetGoalForYear(userID int64, year int) (*models.PersonalYearlyGoal, error) {
	goal, err := s.dao.GetByUserAndYear(userID, year)
	if err != nil {
		return nil, err
	}
	
	// If no goal exists, return a default one (not saved yet)
	if goal == nil {
		goal = &models.PersonalYearlyGoal{
			UserID:        userID,
			Year:          year,
			DistanceGoal:  1000, // Default 1000km
			ElevationGoal: 50000, // Default 50,000m
			SummitGoal:    20,   // Default 20 summits
			TargetSummits: []int64{},
		}
	}
	
	return goal, nil
}

// GetAllGoals gets all yearly goals for a user (history)
func (s *PersonalGoalsService) GetAllGoals(userID int64) ([]models.PersonalYearlyGoal, error) {
	return s.dao.GetByUser(userID)
}

// SaveGoal creates or updates a user's yearly goal
func (s *PersonalGoalsService) SaveGoal(goal *models.PersonalYearlyGoal) error {
	return s.dao.Upsert(goal)
}

// DeleteGoal removes a user's goal for a specific year
func (s *PersonalGoalsService) DeleteGoal(userID int64, year int) error {
	return s.dao.Delete(userID, year)
}

// AddTargetSummit adds a peak to the user's target summit list for the current year
func (s *PersonalGoalsService) AddTargetSummit(userID int64, peakID int64) error {
	goal, err := s.GetCurrentYearGoal(userID)
	if err != nil {
		return err
	}
	
	// Check if peak is already in list
	for _, id := range goal.TargetSummits {
		if id == peakID {
			return nil // Already exists
		}
	}
	
	goal.TargetSummits = append(goal.TargetSummits, peakID)
	return s.dao.Upsert(goal)
}

// RemoveTargetSummit removes a peak from the user's target summit list
func (s *PersonalGoalsService) RemoveTargetSummit(userID int64, peakID int64) error {
	goal, err := s.GetCurrentYearGoal(userID)
	if err != nil {
		return err
	}
	
	// Filter out the peak
	newTargets := []int64{}
	for _, id := range goal.TargetSummits {
		if id != peakID {
			newTargets = append(newTargets, id)
		}
	}
	
	goal.TargetSummits = newTargets
	return s.dao.Upsert(goal)
}
