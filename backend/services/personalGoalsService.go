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

	// If no goal exists, return empty goal with 0 values (no defaults for past years)
	// The frontend will handle showing defaults only for the current year when setting new goals
	if goal == nil {
		goal = &models.PersonalYearlyGoal{
			UserID:        userID,
			Year:          year,
			DistanceGoal:  0,
			ElevationGoal: 0,
			SummitGoal:    0,
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
