package daos

import (
	"database/sql"
	"log"
	"run-goals/models"
)

type UserDao struct {
	l        *log.Logger
	database *sql.DB
}

func NewUserDao(logger *log.Logger, db *sql.DB) *UserDao {
	return &UserDao{
		l:        logger,
		database: db,
	}
}

func (dao *UserDao) UpsertUser(user *models.User) error {
	sql := `
		INSERT INTO user (
			strava_athlete_id,
			access_token,
			refresh_token,
			expires_at,
			last_distance,
			last_update,
			created_at,
			updated_at
		) VALUES (
			($1, $2, $3, $4, $5, $6, $7, $8)
		) ON CONFLICT (

		) DO UPDATE
			SET
				strava_athlete_id = EXCLUDED.strava_athlete_id,
				access_token = EXCLUDED.access_token,
				refresh_token = EXCLUDED.refresh_token,
				expires_at = EXCLUDED.expires_at,
				last_distance = EXCLUDED.last_distance,
				created_at = EXCLUDED.created_at,
				updated_at = EXCLUDED.updated_at;
	`
	_, err := dao.database.Exec(
		sql,
		user.StravaAthleteID,
		user.AccessToken,
		user.RefreshToken,
		user.ExpiresAt,
		user.LastDistance,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		dao.l.Printf("Error upserting user: %v", err)
		return err
	}
	return nil
}
