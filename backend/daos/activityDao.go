package daos

import (
	"database/sql"
	"log"
	"run-goals/models"
)

type ActivityDaoInterface interface {
	UpsertActivity(activity *models.Activity) error
	GetActivitiesByUserID(userID int64) ([]models.Activity, error)
	GetActivityByID(id int64) (models.Activity, error)
	GetActivities() ([]models.Activity, error)
}

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
            strava_athlete_id,
            user_id,
            name,
            description,
            distance,
            elevation,
            moving_time,
            start_date,
            map_polyline,
            created_at,
            updated_at,
            has_summit,
            photo_url
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
        ) ON CONFLICT (
            strava_activity_id
        ) DO UPDATE
            SET
                user_id = EXCLUDED.user_id,
                name = EXCLUDED.name,
                description = EXCLUDED.description,
                distance = EXCLUDED.distance,
                elevation = EXCLUDED.elevation,
                moving_time = EXCLUDED.moving_time,
                start_date = EXCLUDED.start_date,
                map_polyline = EXCLUDED.map_polyline,
                updated_at = EXCLUDED.updated_at,
                has_summit = EXCLUDED.has_summit,
                photo_url = EXCLUDED.photo_url;
    `
	_, err := dao.db.Exec(
		sql,
		activity.StravaActivityId,
		activity.StravaAthleteId,
		activity.UserID,
		activity.Name,
		activity.Description,
		activity.Distance,
		activity.Elevation,
		activity.MovingTime,
		activity.StartDate,
		activity.MapPolyline,
		activity.CreatedAt,
		activity.UpdatedAt,
		activity.HasSummit,
		activity.PhotoURL,
	)
	if err != nil {
		dao.l.Printf("Error upserting activity: %v", err)
		return err
	}
	return nil
}

func (dao *ActivityDao) GetActivitiesByUserID(userID int64) ([]models.Activity, error) {
	activities := []models.Activity{}
	sqlQuery := `
        SELECT
            id,
            strava_activity_id,
            strava_athlete_id,
            user_id,
            name,
            description,
            distance,
            elevation,
            moving_time,
            start_date,
            map_polyline,
            photo_url
        FROM activity
        WHERE
            user_id = $1;
    `
	rows, err := dao.db.Query(sqlQuery, userID)
	if err != nil {
		dao.l.Println("Error querying activity table", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		activity := models.Activity{}
		var elevation sql.NullFloat64
		var movingTime sql.NullFloat64
		err = rows.Scan(
			&activity.ID,
			&activity.StravaActivityId,
			&activity.StravaAthleteId,
			&activity.UserID,
			&activity.Name,
			&activity.Description,
			&activity.Distance,
			&elevation,
			&movingTime,
			&activity.StartDate,
			&activity.MapPolyline,
			&activity.PhotoURL,
		)
		if err != nil {
			dao.l.Println("Error parsing query result", err)
			return nil, err
		}

		// Set elevation value, defaulting to 0 if NULL
		if elevation.Valid {
			activity.Elevation = elevation.Float64
		} else {
			activity.Elevation = 0
		}

		// Set moving_time value, defaulting to 0 if NULL
		if movingTime.Valid {
			activity.MovingTime = movingTime.Float64
		} else {
			activity.MovingTime = 0
		}

		activities = append(activities, activity)
	}
	err = rows.Err()
	if err != nil {
		dao.l.Println("Error during iteration", err)
		return nil, err
	}

	return activities, nil
}

func (dao *ActivityDao) GetActivityByID(id int64) (models.Activity, error) {
	activity := models.Activity{}
	sqlQuery := `
        SELECT
            id,
            strava_activity_id,
            strava_athlete_id,
            user_id,
            name,
            description,
            distance,
            elevation,
            moving_time,
            start_date,
            map_polyline,
            photo_url
        FROM activity
        WHERE
            id = $1;
    `
	row := dao.db.QueryRow(sqlQuery, id)
	var elevation sql.NullFloat64
	var movingTime sql.NullFloat64
	err := row.Scan(
		&activity.ID,
		&activity.StravaActivityId,
		&activity.StravaAthleteId,
		&activity.UserID,
		&activity.Name,
		&activity.Description,
		&activity.Distance,
		&elevation,
		&movingTime,
		&activity.StartDate,
		&activity.MapPolyline,
		&activity.PhotoURL,
	)
	if err != nil {
		dao.l.Println("Error querying activity table", err)
		return models.Activity{}, err
	}

	// Set elevation value, defaulting to 0 if NULL
	if elevation.Valid {
		activity.Elevation = elevation.Float64
	} else {
		activity.Elevation = 0
	}

	// Set moving_time value, defaulting to 0 if NULL
	if movingTime.Valid {
		activity.MovingTime = movingTime.Float64
	} else {
		activity.MovingTime = 0
	}

	return activity, nil
}

func (dao *ActivityDao) GetActivities() ([]models.Activity, error) {
	activities := []models.Activity{}
	sqlQuery := `
        SELECT
            id,
            strava_activity_id,
            strava_athlete_id,
            user_id,
            name,
            description,
            distance,
            elevation,
            moving_time,
            start_date,
            map_polyline,
            photo_url
        FROM activity
    `
	rows, err := dao.db.Query(sqlQuery)
	if err != nil {
		dao.l.Println("Error querying activity table", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		activity := models.Activity{}
		var elevation sql.NullFloat64
		var movingTime sql.NullFloat64
		err = rows.Scan(
			&activity.ID,
			&activity.StravaActivityId,
			&activity.StravaAthleteId,
			&activity.UserID,
			&activity.Name,
			&activity.Description,
			&activity.Distance,
			&elevation,
			&movingTime,
			&activity.StartDate,
			&activity.MapPolyline,
			&activity.PhotoURL,
		)
		if err != nil {
			dao.l.Println("Error parsing query result", err)
			return nil, err
		}

		// Set elevation value, defaulting to 0 if NULL
		if elevation.Valid {
			activity.Elevation = elevation.Float64
		} else {
			activity.Elevation = 0
		}

		// Set moving_time value, defaulting to 0 if NULL
		if movingTime.Valid {
			activity.MovingTime = movingTime.Float64
		} else {
			activity.MovingTime = 0
		}

		activities = append(activities, activity)
	}
	err = rows.Err()
	if err != nil {
		dao.l.Println("Error during iteration", err)
		return nil, err
	}

	return activities, nil
}
