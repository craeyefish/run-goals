package services

import (
	"database/sql"
	"log"
	"run-goals/daos"
	"run-goals/models"
)

type PeakServiceInterface interface {
	ListPeaks() ([]models.Peak, error)
}

type PeakService struct {
	l       *log.Logger
	peakDao *daos.PeakDao
}

func NewPeakService(
	l *log.Logger,
	db *sql.DB,
) *PeakService {
	peakDao := daos.NewPeakDao(l, db)
	return &PeakService{
		l:       l,
		peakDao: peakDao,
	}
}
