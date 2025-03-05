package daos

import (
	"database/sql"
	"log"
	"run-goals/models"
)

type ActivityDao struct {
	l  *log.Logger
	db *sql.DB
}

func NewActivityDao(logger *log.Logger, db *sql.DB) *ActivityDao {
	return &ActivityDao{
		l:  logger,
		db: db,
	}
}

func (dao *ActivityDao) UpsertActivity(activity *models.Activity) error {
	sql := `
		INSERT INTO activity (
			strava_activity_id,
			user_id,
			name,
			distance,
			start_date,
			map_polyline,
			created_at,
			updated_at,
			has_summit
		) VALUES (
			($1, $2, $3, $4, $5, $6, $7, $8, $9)
		) ON CONFLICT (
			strava_activity_id
		) DO UPDATE
			SET
				user_id = EXCLUDED.user_id,
    			name = EXCLUDED.name,
       			distance = EXCLUDED.distance,
          		start_date = EXCLUDED.start_date,
            	map_polyline = EXCLUDED.map_polyline,
             	updated_at = EXCLUDED.updated_at,
              	has_summit = EXCLUDED.has_summit;
	`
	_, err := dao.db.Exec(
		sql,
		activity.StravaActivityID,
		activity.UserID,
		activity.Name,
		activity.Distance,
		activity.StartDate,
		activity.MapPolyline,
		activity.CreatedAt,
		activity.UpdatedAt,
		activity.HasSummit,
	)
	if err != nil {
		dao.l.Printf("Error upserting activity: %v", err)
		return err
	}
	return nil
}

func (dao *ActivityDao) GetActivitiesByUserID(userID int64) ([]models.Activity, error) {
	activities := []models.Activity{}
	sql := `
		SELECT
			id,
			strava_activity_id,
			user_id,
			name,
			distance,
			start_date,
			map_polyline
		FROM activity
		WHERE
			user_id = $1;
	`
	rows, err := dao.db.Query(sql, userID)
	if err != nil {
		dao.l.Println("Error querying activity table", err)
	}
	defer rows.Close()
	for rows.Next() {
		activity := models.Activity{}
		err = rows.Scan(
			&activity.ID,
			&activity.StravaActivityID,
			&activity.UserID,
			&activity.Name,
			&activity.Distance,
			&activity.StartDate,
			&activity.MapPolyline,
		)
		if err != nil {
			dao.l.Println("Error parsing query result", err)
		}
		activities = append(activities, activity)
	}
	err = rows.Err()
	if err != nil {
		dao.l.Println("Error during iteration", err)
	}

	return activities, nil
}
