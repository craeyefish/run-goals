package daos

import (
	"database/sql"
	"log"
	"run-goals/models"
)

type PeakDaoInterface interface {
	GetPeaks() ([]models.Peak, error)
}

type PeakDao struct {
	l  *log.Logger
	db *sql.DB
}

func NewPeakDao(logger *log.Logger, db *sql.DB) *PeakDao {
	return &PeakDao{
		l:  logger,
		db: db,
	}
}

func (dao *PeakDao) GetPeaks() ([]models.Peak, error) {
	limit := 1000
	peaks := []models.Peak{}
	sql := `
		SELECT
			id,
			osm_id,
			latitude,
			longitude,
			name,
			elevation_meters
		FROM peak
		LIMIT $1
	`
	rows, err := dao.db.Query(sql, limit)
	if err != nil {
		dao.l.Println("Error querying peak table", err)
	}
	defer rows.Close()
	for rows.Next() {
		peak := models.Peak{}
		err = rows.Scan(
			&peak.ID,
			&peak.OsmID,
			&peak.Latitude,
			&peak.Longitude,
			&peak.Name,
			&peak.ElevationMeters,
		)
		if err != nil {
			dao.l.Println("Error parsing query result", err)
		}
		peaks = append(peaks, peak)
	}
	err = rows.Err()
	if err != nil {
		dao.l.Println("Error during iteration", err)
	}

	return peaks, nil
}
