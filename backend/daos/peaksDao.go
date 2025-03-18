package daos

import (
	"database/sql"
	"log"
	"run-goals/models"
)

type PeaksDaoInterface interface {
	GetPeaks() ([]models.Peak, error)
	UpsertPeak(models.Peak) error
	GetPeaksBetweenLatLon(minLat float64, maxLat float64, minLon float64, maxLon float64) ([]models.Peak, error)
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
		FROM peaks
	`
	rows, err := dao.db.Query(sql)
	if err != nil {
		dao.l.Println("Error querying peaks table", err)
		return nil, err
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
			return nil, err
		}
		peaks = append(peaks, peak)
	}
	err = rows.Err()
	if err != nil {
		dao.l.Println("Error during iteration", err)
		return nil, err
	}

	return peaks, nil
}

func (dao *PeaksDao) UpsertPeak(peak *models.Peak) error {
	sql := `
		INSERT INTO peaks (
			osm_id,
			latitude,
			longitude,
			name,
			elevation_meters
		) VALUES (
			$1, $2, $3, $4, $5
		) ON CONFLICT (
			osm_id
		) DO UPDATE
			SET
				osm_id = EXCLUDED.osm_id,
				latitude = EXCLUDED.latitude,
				longitude = EXCLUDED.longitude,
				name = EXCLUDED.name,
				elevation_meters = EXCLUDED.elevation_meters;
	`
	_, err := dao.db.Exec(
		sql,
		peak.OsmID,
		peak.Latitude,
		peak.Longitude,
		peak.Name,
		peak.ElevationMeters,
	)
	if err != nil {
		dao.l.Printf("Error upserting peak: %v", err)
		return err
	}
	return nil
}

func (dao *PeaksDao) GetPeaksBetweenLatLon(minLat float64, maxLat float64, minLon float64, maxLon float64) ([]models.Peak, error) {
	peaks := []models.Peak{}
	sql := `
		SELECT
			id,
			osm_id,
			latitude,
			longitude,
			name,
			elevation_meters
		FROM peaks
		WHERE
			lat BETWEEN ? AND ?
			AND lon BETWEEN ? AND ?
	`
	rows, err := dao.db.Query(sql)
	if err != nil {
		dao.l.Println("Error querying peak table", err)
		return nil, err
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
			return nil, err
		}
		peaks = append(peaks, peak)
	}
	err = rows.Err()
	if err != nil {
		dao.l.Println("Error during iteration", err)
		return nil, err
	}

	return peaks, nil
}
