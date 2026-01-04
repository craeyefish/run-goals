package daos

import (
	"database/sql"
	"log"
	"run-goals/models"
	"time"
)

type ActivityDaoInterface interface {
	UpsertActivity(activity *models.Activity) error
	GetActivitiesByUserID(userID int64) ([]models.Activity, error)
	GetActivitiesByUserIDAndDateRange(userID int64, startDate, endDate *time.Time) ([]models.Activity, error)
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
            activity_type,
            sport_type,
            description,
            distance,
            elevation,
            moving_time,
            start_date,
            map_polyline,
            created_at,
            updated_at,
            has_summit,
            summits_calculated,
            photo_url
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
        ) ON CONFLICT (
            strava_activity_id
        ) DO UPDATE
            SET
                user_id = EXCLUDED.user_id,
                name = EXCLUDED.name,
                activity_type = EXCLUDED.activity_type,
                sport_type = EXCLUDED.sport_type,
                description = EXCLUDED.description,
                distance = EXCLUDED.distance,
                elevation = EXCLUDED.elevation,
                moving_time = EXCLUDED.moving_time,
                start_date = EXCLUDED.start_date,
                map_polyline = EXCLUDED.map_polyline,
                updated_at = EXCLUDED.updated_at,
                has_summit = EXCLUDED.has_summit,
                summits_calculated = EXCLUDED.summits_calculated,
                photo_url = EXCLUDED.photo_url;
    `
	_, err := dao.db.Exec(
		sql,
		activity.StravaActivityId,
		activity.StravaAthleteId,
		activity.UserID,
		activity.Name,
		activity.Type,
		activity.SportType,
		activity.Description,
		activity.Distance,
		activity.Elevation,
		activity.MovingTime,
		activity.StartDate,
		activity.MapPolyline,
		activity.CreatedAt,
		activity.UpdatedAt,
		activity.HasSummit,
		activity.SummitsCalculated,
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

// GetActivitiesPendingSummitCalculation returns activities that haven't had summit detection run
func (dao *ActivityDao) GetActivitiesPendingSummitCalculation() ([]models.Activity, error) {
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
            photo_url,
            has_summit,
            COALESCE(summits_calculated, false) as summits_calculated
        FROM activity
        WHERE summits_calculated IS NULL OR summits_calculated = false
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
			&activity.HasSummit,
			&activity.SummitsCalculated,
		)
		if err != nil {
			dao.l.Println("Error parsing query result", err)
			return nil, err
		}

		if elevation.Valid {
			activity.Elevation = elevation.Float64
		}
		if movingTime.Valid {
			activity.MovingTime = movingTime.Float64
		}

		activities = append(activities, activity)
	}
	if err = rows.Err(); err != nil {
		dao.l.Println("Error during iteration", err)
		return nil, err
	}

	return activities, nil
}

// GetActivityByStravaID fetches a single activity by its Strava ID
func (dao *ActivityDao) GetActivityByStravaID(stravaActivityID int64) (*models.Activity, error) {
	activity := &models.Activity{}
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
            photo_url,
            has_summit,
            COALESCE(summits_calculated, false) as summits_calculated
        FROM activity
        WHERE strava_activity_id = $1
    `
	var elevation sql.NullFloat64
	var movingTime sql.NullFloat64
	err := dao.db.QueryRow(sqlQuery, stravaActivityID).Scan(
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
		&activity.HasSummit,
		&activity.SummitsCalculated,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		dao.l.Printf("Error querying activity by strava ID: %v", err)
		return nil, err
	}

	if elevation.Valid {
		activity.Elevation = elevation.Float64
	}
	if movingTime.Valid {
		activity.MovingTime = movingTime.Float64
	}

	return activity, nil
}

// DeleteNonAllowedActivityTypes removes activities that aren't Run, Walk, or Hike
// Returns the count of deleted activities
func (dao *ActivityDao) DeleteNonAllowedActivityTypes() (int64, error) {
	allowedTypes := []string{"Run", "Walk", "Hike"}
	sqlQuery := `
        DELETE FROM activity
        WHERE activity_type IS NULL 
           OR activity_type NOT IN ($1, $2, $3)
        RETURNING id;
    `
	result, err := dao.db.Exec(sqlQuery, allowedTypes[0], allowedTypes[1], allowedTypes[2])
	if err != nil {
		dao.l.Printf("Error deleting non-allowed activity types: %v", err)
		return 0, err
	}
	count, _ := result.RowsAffected()
	dao.l.Printf("Deleted %d activities with non-allowed types", count)
	return count, nil
}

// ResetSummitsCalculated resets summits_calculated flag on all activities
// This forces re-calculation of summit detection
func (dao *ActivityDao) ResetSummitsCalculated() (int64, error) {
	sqlQuery := `
		UPDATE activity
		SET summits_calculated = false, has_summit = false
		WHERE summits_calculated = true
	`
	result, err := dao.db.Exec(sqlQuery)
	if err != nil {
		dao.l.Printf("Error resetting summits_calculated: %v", err)
		return 0, err
	}
	count, _ := result.RowsAffected()
	dao.l.Printf("Reset summits_calculated on %d activities", count)
	return count, nil
}

// GetActivitiesByUserIDAndDateRange returns activities for a user within a date range
func (dao *ActivityDao) GetActivitiesByUserIDAndDateRange(userID int64, startDate, endDate *time.Time) ([]models.Activity, error) {
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
        WHERE user_id = $1
          AND ($2::timestamp IS NULL OR start_date >= $2)
          AND ($3::timestamp IS NULL OR start_date <= $3)
        ORDER BY start_date DESC;
    `
	rows, err := dao.db.Query(sqlQuery, userID, startDate, endDate)
	if err != nil {
		dao.l.Println("Error querying activity table with date range", err)
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
