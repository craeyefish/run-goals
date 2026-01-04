package daos

import (
	"database/sql"
	"errors"
	"log"
	"run-goals/models"
)

var ErrUserNotFound = errors.New("user not found")

type UserDaoInterface interface {
	UpsertUser(user *models.User) error
	GetUsers() ([]models.User, error)
}

type UserDao struct {
	l  *log.Logger
	db *sql.DB
}

func NewUserDao(logger *log.Logger, db *sql.DB) *UserDao {
	return &UserDao{
		l:  logger,
		db: db,
	}
}

func (dao *UserDao) UpsertUser(user *models.User) error {
	sql := `
		INSERT INTO users (
			strava_athlete_id,
			access_token,
			refresh_token,
			expires_at,
			last_distance,
			last_updated,
			created_at,
			updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		) ON CONFLICT (
			strava_athlete_id
		) DO UPDATE
			SET
				strava_athlete_id = EXCLUDED.strava_athlete_id,
				access_token = EXCLUDED.access_token,
				refresh_token = EXCLUDED.refresh_token,
				expires_at = EXCLUDED.expires_at,
				last_distance = EXCLUDED.last_distance,
				last_updated = EXCLUDED.last_updated,
				created_at = EXCLUDED.created_at,
				updated_at = EXCLUDED.updated_at;
	`
	_, err := dao.db.Exec(
		sql,
		user.StravaAthleteID,
		user.AccessToken,
		user.RefreshToken,
		user.ExpiresAt,
		user.LastDistance,
		user.LastUpdated,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		dao.l.Printf("Error upserting user: %v", err)
		return err
	}
	return nil
}

func (dao *UserDao) GetUsers() ([]models.User, error) {
	users := []models.User{}
	sql := `
		SELECT
			id,
			strava_athlete_id,
			username,
			is_admin,
			access_token,
			refresh_token,
			expires_at,
			last_distance,
			last_updated,
			created_at,
			updated_at
		FROM users;
	`
	rows, err := dao.db.Query(sql)
	if err != nil {
		dao.l.Println("Error querying user table", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		user := models.User{}
		err = rows.Scan(
			&user.ID,
			&user.StravaAthleteID,
			&user.Username,
			&user.IsAdmin,
			&user.AccessToken,
			&user.RefreshToken,
			&user.ExpiresAt,
			&user.LastDistance,
			&user.LastUpdated,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			dao.l.Println("Error parsing query result", err)
			return nil, err
		}
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		dao.l.Println("Error during iteration", err)
		return nil, err
	}

	return users, nil
}

func (dao *UserDao) GetUserByID(id int64) (*models.User, error) {
	user := models.User{}
	query := `
		SELECT
			id,
			strava_athlete_id,
			username,
			is_admin,
			access_token,
			refresh_token,
			expires_at,
			last_distance,
			last_updated,
			created_at,
			updated_at
		FROM users
		WHERE
			id = $1;
	`
	row := dao.db.QueryRow(query, id)
	err := row.Scan(
		&user.ID,
		&user.StravaAthleteID,
		&user.Username,
		&user.IsAdmin,
		&user.AccessToken,
		&user.RefreshToken,
		&user.ExpiresAt,
		&user.LastDistance,
		&user.LastUpdated,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		dao.l.Printf("No user found with id=%d", id)
		return nil, ErrUserNotFound
	} else if err != nil {
		dao.l.Println("Error querying user table", err)
		return nil, err
	}

	return &user, nil
}

func (dao *UserDao) GetUserByStravaAthleteID(id int64) (*models.User, error) {
	user := models.User{}
	query := `
		SELECT
			id,
			strava_athlete_id,
			username,
			is_admin,
			access_token,
			refresh_token,
			expires_at,
			last_distance,
			last_updated,
			created_at,
			updated_at
		FROM users
		WHERE
			strava_athlete_id = $1;
	`
	row := dao.db.QueryRow(query, id)
	err := row.Scan(
		&user.ID,
		&user.StravaAthleteID,
		&user.Username,
		&user.IsAdmin,
		&user.AccessToken,
		&user.RefreshToken,
		&user.ExpiresAt,
		&user.LastDistance,
		&user.LastUpdated,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		dao.l.Printf("No user found with strava_athlete_id=%d", id)
		return nil, ErrUserNotFound
	} else if err != nil {
		dao.l.Println("Error querying user table", err)
		return nil, err
	}

	return &user, nil
}

// UpdateUsername updates just the username for a user
func (dao *UserDao) UpdateUsername(userID int64, username string) error {
	query := `UPDATE users SET username = $1, updated_at = NOW() WHERE id = $2`
	result, err := dao.db.Exec(query, username, userID)
	if err != nil {
		dao.l.Printf("Error updating username for user_id=%d: %v", userID, err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		dao.l.Printf("Error getting rows affected: %v", err)
		return err
	}

	if rowsAffected == 0 {
		dao.l.Printf("No user found with id=%d", userID)
		return ErrUserNotFound
	}

	return nil
}

func (dao *UserDao) DeleteUserByStravaAthleteID(stravaAthleteID int64) error {
	// Due to CASCADE DELETE constraints, this will automatically delete:
	// - activities (via strava_athlete_id FK)
	// - user_peaks (via user_id FK)
	// - group_members (via user_id FK)
	query := `DELETE FROM users WHERE strava_athlete_id = $1`
	result, err := dao.db.Exec(query, stravaAthleteID)
	if err != nil {
		dao.l.Printf("Error deleting user with strava_athlete_id=%d: %v", stravaAthleteID, err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		dao.l.Printf("Error getting rows affected: %v", err)
		return err
	}

	if rowsAffected == 0 {
		dao.l.Printf("No user found with strava_athlete_id=%d", stravaAthleteID)
		return ErrUserNotFound
	}

	dao.l.Printf("Successfully deleted user with strava_athlete_id=%d", stravaAthleteID)
	return nil
}
