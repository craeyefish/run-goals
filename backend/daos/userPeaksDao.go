package daos

import (
	"database/sql"
	"log"
	"run-goals/models"
	"time"

	"github.com/lib/pq"
)

type UserPeaksDaoInterface interface {
	GetUserPeaks() ([]models.UserPeak, error)
	GetUserPeaksJoin() ([]models.UserPeakJoin, error)
	GetUserPeaksJoinByUserID(userID int64) ([]models.UserPeakJoin, error)
	UpsertUserPeak(userPeak *models.UserPeak) error
	ClearUserPeaks() error
	GetUserSummitsInDateRange(userID int64, peakIDs []int64, startDate time.Time, endDate time.Time) ([]models.UserPeak, error)
	GetUserSummitsInDateRangeAll(userID int64, startDate time.Time, endDate time.Time) ([]models.UserPeak, error)
}

type UserPeaksDao struct {
	l  *log.Logger
	db *sql.DB
}

func NewUserPeaksDao(logger *log.Logger, db *sql.DB) *UserPeaksDao {
	return &UserPeaksDao{
		l:  logger,
		db: db,
	}
}

func (dao *UserPeaksDao) GetUserPeaks() ([]models.UserPeak, error) {
	userPeaks := []models.UserPeak{}
	sql := `
		SELECT
			id,
			user_id,
			peak_id,
			activity_id,
			summited_at
		FROM user_peaks;
	`
	rows, err := dao.db.Query(sql)
	if err != nil {
		dao.l.Println("Error querying user_peaks table", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		userPeak := models.UserPeak{}
		err = rows.Scan(
			&userPeak.ID,
			&userPeak.UserID,
			&userPeak.PeakID,
			&userPeak.ActivityID,
			&userPeak.SummitedAt,
		)
		if err != nil {
			dao.l.Println("Error parsing query result", err)
			return nil, err
		}
		userPeaks = append(userPeaks, userPeak)
	}
	err = rows.Err()
	if err != nil {
		dao.l.Println("Error during iteration", err)
		return nil, err
	}

	return userPeaks, nil
}

func (dao *UserPeaksDao) GetUserPeaksJoin() ([]models.UserPeakJoin, error) {
	userPeaksJoin := []models.UserPeakJoin{}
	sql := `
		SELECT
  			up.peak_id,
     		up.user_id,
       		up.activity_id,
         	up.summited_at,
          	u.strava_athlete_id AS user_name
		FROM user_peaks up
		LEFT JOIN users u ON up.user_id = u.id;
	`
	rows, err := dao.db.Query(sql)
	if err != nil {
		dao.l.Println("Error querying user_peaks table", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		userPeakJoin := models.UserPeakJoin{}
		err = rows.Scan(
			&userPeakJoin.PeakID,
			&userPeakJoin.UserID,
			&userPeakJoin.ActivityID,
			&userPeakJoin.SummitedAt,
			&userPeakJoin.UserName,
		)
		if err != nil {
			dao.l.Println("Error parsing query result", err)
			return nil, err
		}
		userPeaksJoin = append(userPeaksJoin, userPeakJoin)
	}
	err = rows.Err()
	if err != nil {
		dao.l.Println("Error during iteration", err)
		return nil, err
	}

	return userPeaksJoin, nil
}

func (dao *UserPeaksDao) GetUserPeaksJoinByUserID(userID int64) ([]models.UserPeakJoin, error) {
	userPeaksJoin := []models.UserPeakJoin{}
	sql := `
		SELECT
			up.peak_id,
			up.user_id,
			up.activity_id,
			up.summited_at,
			u.strava_athlete_id AS user_name
		FROM user_peaks up
		LEFT JOIN users u ON up.user_id = u.id
		WHERE up.user_id = $1;
	`
	rows, err := dao.db.Query(sql, userID)
	if err != nil {
		dao.l.Println("Error querying user_peaks table", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		userPeakJoin := models.UserPeakJoin{}
		err = rows.Scan(
			&userPeakJoin.PeakID,
			&userPeakJoin.UserID,
			&userPeakJoin.ActivityID,
			&userPeakJoin.SummitedAt,
			&userPeakJoin.UserName,
		)
		if err != nil {
			dao.l.Println("Error parsing query result", err)
			return nil, err
		}
		userPeaksJoin = append(userPeaksJoin, userPeakJoin)
	}
	err = rows.Err()
	if err != nil {
		dao.l.Println("Error during iteration", err)
		return nil, err
	}

	return userPeaksJoin, nil
}

func (dao *UserPeaksDao) UpsertUserPeak(userPeak *models.UserPeak) error {
	sql := `
        INSERT INTO user_peaks (
            user_id,
            peak_id,
            activity_id,
            summited_at
        ) VALUES (
            $1, $2, $3, $4
        ) ON CONFLICT (user_id, peak_id, activity_id) 
        DO UPDATE SET summited_at = EXCLUDED.summited_at;
    `
	_, err := dao.db.Exec(
		sql,
		userPeak.UserID,
		userPeak.PeakID,
		userPeak.ActivityID,
		userPeak.SummitedAt,
	)
	if err != nil {
		dao.l.Printf("Error upserting userPeak: %v", err)
		return err
	}
	return nil
}

func (dao *UserPeaksDao) ClearUserPeaks() error {
	sql := `
		DELETE FROM user_peaks;
	`
	_, err := dao.db.Exec(sql)
	if err != nil {
		dao.l.Printf("Error deleting records from user_peaks: %v", err)
		return err
	}
	return nil
}

func (dao *UserPeaksDao) GetUserSummitsInDateRange(userID int64, peakIDs []int64, startDate time.Time, endDate time.Time) ([]models.UserPeak, error) {
	userPeaks := []models.UserPeak{}

	// If no specific peaks are provided, return empty result
	if len(peakIDs) == 0 {
		return userPeaks, nil
	}

	sql := `
        SELECT
            id,
            user_id,
            peak_id,
            activity_id,
            summited_at
        FROM user_peaks
        WHERE 
            user_id = $1
            AND peak_id = ANY($2)
            AND summited_at >= $3
            AND summited_at <= $4
        ORDER BY summited_at DESC
    `

	rows, err := dao.db.Query(sql, userID, pq.Array(peakIDs), startDate, endDate)
	if err != nil {
		dao.l.Printf("Error querying user summits in date range: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		userPeak := models.UserPeak{}
		err = rows.Scan(
			&userPeak.ID,
			&userPeak.UserID,
			&userPeak.PeakID,
			&userPeak.ActivityID,
			&userPeak.SummitedAt,
		)
		if err != nil {
			dao.l.Printf("Error parsing user summit result: %v", err)
			return nil, err
		}
		userPeaks = append(userPeaks, userPeak)
	}

	err = rows.Err()
	if err != nil {
		dao.l.Printf("Error during user summits iteration: %v", err)
		return nil, err
	}

	return userPeaks, nil
}

func (dao *UserPeaksDao) GetUserSummitsInDateRangeAll(userID int64, startDate time.Time, endDate time.Time) ([]models.UserPeak, error) {
	userPeaks := []models.UserPeak{}

	sql := `
        SELECT
            id,
            user_id,
            peak_id,
            activity_id,
            summited_at
        FROM user_peaks
        WHERE 
            user_id = $1
            AND summited_at >= $2
            AND summited_at <= $3
        ORDER BY summited_at DESC
    `

	rows, err := dao.db.Query(sql, userID, startDate, endDate)
	if err != nil {
		dao.l.Printf("Error querying all user summits: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		userPeak := models.UserPeak{}
		err = rows.Scan(
			&userPeak.ID,
			&userPeak.UserID,
			&userPeak.PeakID,
			&userPeak.ActivityID,
			&userPeak.SummitedAt,
		)
		if err != nil {
			dao.l.Printf("Error parsing user summit result: %v", err)
			return nil, err
		}
		userPeaks = append(userPeaks, userPeak)
	}

	err = rows.Err()
	if err != nil {
		dao.l.Printf("Error during user summits iteration: %v", err)
		return nil, err
	}

	return userPeaks, nil
}
