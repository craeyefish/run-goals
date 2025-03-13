package daos

import (
	"database/sql"
	"log"
	"run-goals/models"
)

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
			($1, $2, $3, $4, $5, $6, $7, $8)
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
	}
	defer rows.Close()
	for rows.Next() {
		user := models.User{}
		err = rows.Scan(
			&user.ID,
			&user.StravaAthleteID,
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
		}
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		dao.l.Println("Error during iteration", err)
	}

	return users, nil
}

func (dao *UserDao) GetUserByStravaAthleteID(id int64) (*models.User, error) {
	user := models.User{}
	sql := `
		SELECT
			id,
			strava_athlete_id,
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
	row := dao.db.QueryRow(sql, id)
	err := row.Scan(
		&user.ID,
		&user.StravaAthleteID,
		&user.AccessToken,
		&user.RefreshToken,
		&user.ExpiresAt,
		&user.LastDistance,
		&user.LastUpdated,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		dao.l.Println("Error querying user table", err)
	}

	return &user, nil
}
