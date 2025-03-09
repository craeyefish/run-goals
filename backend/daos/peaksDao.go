package daos

import (
	"database/sql"
	"log"
	"run-goals/models"
)

type PeaksDaoInterface interface {
	GetPeaks() ([]models.Peak, error)
}

type PeaksDao struct {
	l  *log.Logger
	db *sql.DB
}

func NewPeaksDao(logger *log.Logger, db *sql.DB) *PeaksDao {
	return &PeaksDao{
		l:  logger,
		db: db,
	}
}

func (dao *PeaksDao) GetPeaks() ([]models.Peak, error) {
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
	`
	rows, err := dao.db.Query(sql)
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
