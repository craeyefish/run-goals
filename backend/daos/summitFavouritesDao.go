package daos

import (
	"database/sql"
	"log"
)

type SummitFavouritesDao struct {
	l  *log.Logger
	db *sql.DB
}

func NewSummitFavouritesDao(logger *log.Logger, db *sql.DB) *SummitFavouritesDao {
	return &SummitFavouritesDao{
		l:  logger,
		db: db,
	}
}

// GetAllByUser retrieves all favourite peaks for a user
func (dao *SummitFavouritesDao) GetAllByUser(userID int64) ([]int64, error) {
	query := `
		SELECT peak_id
		FROM summit_favourites
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := dao.db.Query(query, userID)
	if err != nil {
		dao.l.Printf("Error getting summit favourites: %v", err)
		return nil, err
	}
	defer rows.Close()

	var peakIDs []int64
	for rows.Next() {
		var peakID int64
		if err := rows.Scan(&peakID); err != nil {
			dao.l.Printf("Error scanning summit favourite: %v", err)
			continue
		}
		peakIDs = append(peakIDs, peakID)
	}

	if peakIDs == nil {
		peakIDs = []int64{}
	}

	return peakIDs, nil
}

// Add adds a peak to user's favourites
func (dao *SummitFavouritesDao) Add(userID int64, peakID int64) error {
	query := `
		INSERT INTO summit_favourites (user_id, peak_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, peak_id) DO NOTHING
	`
	_, err := dao.db.Exec(query, userID, peakID)
	if err != nil {
		dao.l.Printf("Error adding summit favourite: %v", err)
		return err
	}
	return nil
}

// Remove removes a peak from user's favourites
func (dao *SummitFavouritesDao) Remove(userID int64, peakID int64) error {
	query := `
		DELETE FROM summit_favourites
		WHERE user_id = $1 AND peak_id = $2
	`
	_, err := dao.db.Exec(query, userID, peakID)
	if err != nil {
		dao.l.Printf("Error removing summit favourite: %v", err)
		return err
	}
	return nil
}

// IsFavourite checks if a peak is in user's favourites
func (dao *SummitFavouritesDao) IsFavourite(userID int64, peakID int64) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM summit_favourites
			WHERE user_id = $1 AND peak_id = $2
		)
	`
	var exists bool
	err := dao.db.QueryRow(query, userID, peakID).Scan(&exists)
	if err != nil {
		dao.l.Printf("Error checking summit favourite: %v", err)
		return false, err
	}
	return exists, nil
}
