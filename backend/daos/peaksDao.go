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
			COALESCE(name, ''),
			COALESCE(elevation_meters, 0),
			COALESCE(alt_name, ''),
			COALESCE(name_en, ''),
			COALESCE(region, ''),
			COALESCE(wikipedia, ''),
			COALESCE(wikidata, ''),
			COALESCE(description, ''),
			COALESCE(prominence, 0)
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
			&peak.AltName,
			&peak.NameEN,
			&peak.Region,
			&peak.Wikipedia,
			&peak.Wikidata,
			&peak.Description,
			&peak.Prominence,
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
			elevation_meters,
			alt_name,
			name_en,
			region,
			wikipedia,
			wikidata,
			description,
			prominence
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		) ON CONFLICT (
			osm_id
		) DO UPDATE
			SET
				latitude = EXCLUDED.latitude,
				longitude = EXCLUDED.longitude,
				name = EXCLUDED.name,
				elevation_meters = EXCLUDED.elevation_meters,
				alt_name = EXCLUDED.alt_name,
				name_en = EXCLUDED.name_en,
				region = EXCLUDED.region,
				wikipedia = EXCLUDED.wikipedia,
				wikidata = EXCLUDED.wikidata,
				description = EXCLUDED.description,
				prominence = EXCLUDED.prominence;
	`
	_, err := dao.db.Exec(
		sql,
		peak.OsmID,
		peak.Latitude,
		peak.Longitude,
		peak.Name,
		peak.ElevationMeters,
		peak.AltName,
		peak.NameEN,
		peak.Region,
		peak.Wikipedia,
		peak.Wikidata,
		peak.Description,
		peak.Prominence,
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
			COALESCE(name, ''),
			COALESCE(elevation_meters, 0),
			COALESCE(alt_name, ''),
			COALESCE(name_en, ''),
			COALESCE(region, ''),
			COALESCE(wikipedia, ''),
			COALESCE(wikidata, ''),
			COALESCE(description, ''),
			COALESCE(prominence, 0)
		FROM peaks
		WHERE
			latitude BETWEEN $1 AND $2
			AND longitude BETWEEN $3 AND $4
	`
	rows, err := dao.db.Query(sql, minLat, maxLat, minLon, maxLon)
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
			&peak.AltName,
			&peak.NameEN,
			&peak.Region,
			&peak.Wikipedia,
			&peak.Wikidata,
			&peak.Description,
			&peak.Prominence,
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
