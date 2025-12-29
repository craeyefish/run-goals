package daos

import (
	"database/sql"
	"log"
	"run-goals/models"
	"time"

	"github.com/lib/pq"
)

type PersonalYearlyGoalDao struct {
	l  *log.Logger
	db *sql.DB
}

func NewPersonalYearlyGoalDao(logger *log.Logger, db *sql.DB) *PersonalYearlyGoalDao {
	return &PersonalYearlyGoalDao{
		l:  logger,
		db: db,
	}
}

// GetByUserAndYear retrieves a user's goal for a specific year
func (dao *PersonalYearlyGoalDao) GetByUserAndYear(userID int64, year int) (*models.PersonalYearlyGoal, error) {
	goal := &models.PersonalYearlyGoal{}
	query := `
		SELECT id, user_id, year, distance_goal, elevation_goal, summit_goal, 
			   COALESCE(target_summits, '{}') as target_summits, created_at, updated_at
		FROM personal_yearly_goals
		WHERE user_id = $1 AND year = $2
	`
	var targetSummits pq.Int64Array
	err := dao.db.QueryRow(query, userID, year).Scan(
		&goal.ID,
		&goal.UserID,
		&goal.Year,
		&goal.DistanceGoal,
		&goal.ElevationGoal,
		&goal.SummitGoal,
		&targetSummits,
		&goal.CreatedAt,
		&goal.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		dao.l.Printf("Error getting personal yearly goal: %v", err)
		return nil, err
	}
	goal.TargetSummits = targetSummits
	return goal, nil
}

// GetByUser retrieves all goals for a user (for history view)
func (dao *PersonalYearlyGoalDao) GetByUser(userID int64) ([]models.PersonalYearlyGoal, error) {
	goals := []models.PersonalYearlyGoal{}
	query := `
		SELECT id, user_id, year, distance_goal, elevation_goal, summit_goal,
			   COALESCE(target_summits, '{}') as target_summits, created_at, updated_at
		FROM personal_yearly_goals
		WHERE user_id = $1
		ORDER BY year DESC
	`
	rows, err := dao.db.Query(query, userID)
	if err != nil {
		dao.l.Printf("Error getting personal yearly goals: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var goal models.PersonalYearlyGoal
		var targetSummits pq.Int64Array
		err := rows.Scan(
			&goal.ID,
			&goal.UserID,
			&goal.Year,
			&goal.DistanceGoal,
			&goal.ElevationGoal,
			&goal.SummitGoal,
			&targetSummits,
			&goal.CreatedAt,
			&goal.UpdatedAt,
		)
		if err != nil {
			dao.l.Printf("Error scanning personal yearly goal: %v", err)
			return nil, err
		}
		goal.TargetSummits = targetSummits
		goals = append(goals, goal)
	}

	return goals, nil
}

// Upsert creates or updates a user's goal for a specific year
func (dao *PersonalYearlyGoalDao) Upsert(goal *models.PersonalYearlyGoal) error {
	query := `
		INSERT INTO personal_yearly_goals (
			user_id, year, distance_goal, elevation_goal, summit_goal, target_summits, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (user_id, year) DO UPDATE SET
			distance_goal = EXCLUDED.distance_goal,
			elevation_goal = EXCLUDED.elevation_goal,
			summit_goal = EXCLUDED.summit_goal,
			target_summits = EXCLUDED.target_summits,
			updated_at = EXCLUDED.updated_at
		RETURNING id
	`
	now := time.Now()
	if goal.CreatedAt.IsZero() {
		goal.CreatedAt = now
	}
	goal.UpdatedAt = now

	var targetSummits pq.Int64Array = goal.TargetSummits
	if targetSummits == nil {
		targetSummits = pq.Int64Array{}
	}

	err := dao.db.QueryRow(
		query,
		goal.UserID,
		goal.Year,
		goal.DistanceGoal,
		goal.ElevationGoal,
		goal.SummitGoal,
		targetSummits,
		goal.CreatedAt,
		goal.UpdatedAt,
	).Scan(&goal.ID)

	if err != nil {
		dao.l.Printf("Error upserting personal yearly goal: %v", err)
		return err
	}
	return nil
}

// Delete removes a user's goal for a specific year
func (dao *PersonalYearlyGoalDao) Delete(userID int64, year int) error {
	query := `DELETE FROM personal_yearly_goals WHERE user_id = $1 AND year = $2`
	_, err := dao.db.Exec(query, userID, year)
	if err != nil {
		dao.l.Printf("Error deleting personal yearly goal: %v", err)
		return err
	}
	return nil
}
