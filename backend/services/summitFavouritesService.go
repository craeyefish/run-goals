package services

import (
	"log"
	"run-goals/daos"
)

type SummitFavouritesService struct {
	logger *log.Logger
	dao    *daos.SummitFavouritesDao
}

func NewSummitFavouritesService(logger *log.Logger, dao *daos.SummitFavouritesDao) *SummitFavouritesService {
	return &SummitFavouritesService{
		logger: logger,
		dao:    dao,
	}
}

// GetFavourites returns all favourite peak IDs for a user
func (s *SummitFavouritesService) GetFavourites(userID int64) ([]int64, error) {
	return s.dao.GetAllByUser(userID)
}

// AddFavourite adds a peak to user's favourites
func (s *SummitFavouritesService) AddFavourite(userID int64, peakID int64) error {
	return s.dao.Add(userID, peakID)
}

// RemoveFavourite removes a peak from user's favourites
func (s *SummitFavouritesService) RemoveFavourite(userID int64, peakID int64) error {
	return s.dao.Remove(userID, peakID)
}

// IsFavourite checks if a peak is in user's favourites
func (s *SummitFavouritesService) IsFavourite(userID int64, peakID int64) (bool, error) {
	return s.dao.IsFavourite(userID, peakID)
}
