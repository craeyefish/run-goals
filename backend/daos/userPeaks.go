package daos

import (
	"database/sql"
	"log"
	"run-goals/models"
)

type UserPeaksDaoInterface interface {
	GetUserPeaksJoin() ([]models.UserPeakJoin, error)
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

func (dao *PeakDao) GetUserPeaksJoin() ([]models.UserPeakJoin, error) {
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
	`
	rows, err := dao.db.Query(sql)
	if err != nil {
		dao.l.Println("Error querying user_peaks table", err)
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
		}
		userPeaksJoin = append(userPeaksJoin, userPeakJoin)
	}
	err = rows.Err()
	if err != nil {
		dao.l.Println("Error during iteration", err)
	}

	return userPeaksJoin, nil
}
