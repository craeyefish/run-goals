package daos

import (
	"database/sql"
	"log"
	"run-goals/models"
	"time"
)

type ChallengeDaoInterface interface {
	// Challenge CRUD
	CreateChallenge(challenge models.Challenge) (*int64, error)
	GetChallengeByID(id int64) (*models.Challenge, error)
	UpdateChallenge(challenge models.Challenge) error
	DeleteChallenge(id int64) error

	// Challenge queries
	GetChallengesByUser(userID int64) ([]models.ChallengeWithProgress, error)
	GetFeaturedChallenges() ([]models.Challenge, error)
	GetPublicChallenges(region *string, limit int, offset int) ([]models.Challenge, error)
	SearchChallenges(query string, limit int) ([]models.Challenge, error)

	// Challenge peaks
	AddChallengePeak(challengeID int64, peakID int64, sortOrder int) error
	RemoveChallengePeak(challengeID int64, peakID int64) error
	GetChallengePeaks(challengeID int64) ([]models.ChallengePeakWithDetails, error)
	SetChallengePeaks(challengeID int64, peakIDs []int64) error

	// Participants
	JoinChallenge(challengeID int64, userID int64) error
	LeaveChallenge(challengeID int64, userID int64) error
	GetChallengeParticipants(challengeID int64) ([]models.ChallengeParticipantWithUser, error)
	GetChallengeLeaderboard(challengeID int64) ([]models.LeaderboardEntry, error)
	UpdateParticipantProgress(challengeID int64, userID int64, peaksCompleted int, totalPeaks int) error
	MarkParticipantCompleted(challengeID int64, userID int64) error
	IsUserParticipant(challengeID int64, userID int64) (bool, error)

	// Groups
	AddGroupToChallenge(challengeID int64, groupID int64, deadlineOverride *time.Time) error
	RemoveGroupFromChallenge(challengeID int64, groupID int64) error
	GetChallengeGroups(challengeID int64) ([]models.ChallengeGroupWithDetails, error)
	GetGroupChallenges(groupID int64) ([]models.Challenge, error)

	// Summit log
	LogSummit(log models.ChallengeSummitLog) error
	GetChallengeSummitLog(challengeID int64, userID *int64) ([]models.ChallengeSummitLogWithDetails, error)
	HasUserSummitedPeakForChallenge(challengeID int64, userID int64, peakID int64) (bool, error)
}

type ChallengeDao struct {
	l  *log.Logger
	db *sql.DB
}

func NewChallengeDao(logger *log.Logger, db *sql.DB) *ChallengeDao {
	return &ChallengeDao{
		l:  logger,
		db: db,
	}
}

// ==================== Challenge CRUD ====================

func (dao *ChallengeDao) CreateChallenge(challenge models.Challenge) (*int64, error) {
	var id int64
	query := `
		INSERT INTO challenges (
			name, description, challenge_type, competition_mode, visibility,
			start_date, deadline, created_by_user_id, created_by_group_id,
			target_count, region, difficulty, is_featured
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		)
		RETURNING id;
	`
	err := dao.db.QueryRow(query,
		challenge.Name, challenge.Description, challenge.ChallengeType, challenge.CompetitionMode, challenge.Visibility,
		challenge.StartDate, challenge.Deadline, challenge.CreatedByUserID, challenge.CreatedByGroupID,
		challenge.TargetCount, challenge.Region, challenge.Difficulty, challenge.IsFeatured,
	).Scan(&id)
	if err != nil {
		dao.l.Printf("Error creating challenge: %v", err)
		return nil, err
	}
	return &id, nil
}

func (dao *ChallengeDao) GetChallengeByID(id int64) (*models.Challenge, error) {
	query := `
		SELECT
			id, name, description, challenge_type, competition_mode, visibility,
			start_date, deadline, created_by_user_id, created_by_group_id,
			target_count, region, difficulty, is_featured, created_at, updated_at
		FROM challenges
		WHERE id = $1;
	`
	var c models.Challenge
	err := dao.db.QueryRow(query, id).Scan(
		&c.ID, &c.Name, &c.Description, &c.ChallengeType, &c.CompetitionMode, &c.Visibility,
		&c.StartDate, &c.Deadline, &c.CreatedByUserID, &c.CreatedByGroupID,
		&c.TargetCount, &c.Region, &c.Difficulty, &c.IsFeatured, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		dao.l.Printf("Error getting challenge by ID: %v", err)
		return nil, err
	}
	return &c, nil
}

func (dao *ChallengeDao) UpdateChallenge(challenge models.Challenge) error {
	query := `
		UPDATE challenges SET
			name = $2,
			description = $3,
			challenge_type = $4,
			competition_mode = $5,
			visibility = $6,
			start_date = $7,
			deadline = $8,
			target_count = $9,
			region = $10,
			difficulty = $11,
			is_featured = $12,
			updated_at = NOW()
		WHERE id = $1;
	`
	_, err := dao.db.Exec(query,
		challenge.ID, challenge.Name, challenge.Description, challenge.ChallengeType, challenge.CompetitionMode,
		challenge.Visibility, challenge.StartDate, challenge.Deadline, challenge.TargetCount,
		challenge.Region, challenge.Difficulty, challenge.IsFeatured,
	)
	if err != nil {
		dao.l.Printf("Error updating challenge: %v", err)
		return err
	}
	return nil
}

func (dao *ChallengeDao) DeleteChallenge(id int64) error {
	query := `DELETE FROM challenges WHERE id = $1;`
	_, err := dao.db.Exec(query, id)
	if err != nil {
		dao.l.Printf("Error deleting challenge: %v", err)
		return err
	}
	return nil
}

// ==================== Challenge Queries ====================

func (dao *ChallengeDao) GetChallengesByUser(userID int64) ([]models.ChallengeWithProgress, error) {
	query := `
		SELECT
			c.id, c.name, c.description, c.challenge_type, c.competition_mode, c.visibility,
			c.start_date, c.deadline, c.created_by_user_id, c.created_by_group_id,
			c.target_count, c.region, c.difficulty, c.is_featured, c.created_at, c.updated_at,
			COALESCE(cp.peaks_completed, 0) as peaks_completed,
			COALESCE(cp.total_peaks, (SELECT COUNT(*) FROM challenge_peaks WHERE challenge_id = c.id)) as total_peaks,
			cp.completed_at IS NOT NULL as is_completed
		FROM challenges c
		LEFT JOIN challenge_participants cp ON c.id = cp.challenge_id AND cp.user_id = $1
		WHERE c.created_by_user_id = $1 
		   OR cp.user_id = $1
		ORDER BY c.created_at DESC;
	`
	rows, err := dao.db.Query(query, userID)
	if err != nil {
		dao.l.Printf("Error getting challenges by user: %v", err)
		return nil, err
	}
	defer rows.Close()

	var challenges []models.ChallengeWithProgress
	for rows.Next() {
		var c models.ChallengeWithProgress
		err := rows.Scan(
			&c.ID, &c.Name, &c.Description, &c.ChallengeType, &c.CompetitionMode, &c.Visibility,
			&c.StartDate, &c.Deadline, &c.CreatedByUserID, &c.CreatedByGroupID,
			&c.TargetCount, &c.Region, &c.Difficulty, &c.IsFeatured, &c.CreatedAt, &c.UpdatedAt,
			&c.CompletedPeaks, &c.TotalPeaks, &c.IsCompleted,
		)
		if err != nil {
			dao.l.Printf("Error scanning challenge: %v", err)
			return nil, err
		}
		c.IsJoined = true
		challenges = append(challenges, c)
	}
	return challenges, nil
}

func (dao *ChallengeDao) GetFeaturedChallenges() ([]models.Challenge, error) {
	query := `
		SELECT
			id, name, description, challenge_type, competition_mode, visibility,
			start_date, deadline, created_by_user_id, created_by_group_id,
			target_count, region, difficulty, is_featured, created_at, updated_at
		FROM challenges
		WHERE is_featured = TRUE AND visibility = 'public'
		ORDER BY name;
	`
	rows, err := dao.db.Query(query)
	if err != nil {
		dao.l.Printf("Error getting featured challenges: %v", err)
		return nil, err
	}
	defer rows.Close()

	return dao.scanChallenges(rows)
}

func (dao *ChallengeDao) GetPublicChallenges(region *string, limit int, offset int) ([]models.Challenge, error) {
	query := `
		SELECT
			id, name, description, challenge_type, competition_mode, visibility,
			start_date, deadline, created_by_user_id, created_by_group_id,
			target_count, region, difficulty, is_featured, created_at, updated_at
		FROM challenges
		WHERE visibility = 'public'
		AND ($1::text IS NULL OR region = $1)
		ORDER BY is_featured DESC, name
		LIMIT $2 OFFSET $3;
	`
	rows, err := dao.db.Query(query, region, limit, offset)
	if err != nil {
		dao.l.Printf("Error getting public challenges: %v", err)
		return nil, err
	}
	defer rows.Close()

	return dao.scanChallenges(rows)
}

func (dao *ChallengeDao) SearchChallenges(queryStr string, limit int) ([]models.Challenge, error) {
	query := `
		SELECT
			id, name, description, challenge_type, competition_mode, visibility,
			start_date, deadline, created_by_user_id, created_by_group_id,
			target_count, region, difficulty, is_featured, created_at, updated_at
		FROM challenges
		WHERE visibility = 'public'
		AND (name ILIKE '%' || $1 || '%' OR region ILIKE '%' || $1 || '%')
		ORDER BY is_featured DESC, name
		LIMIT $2;
	`
	rows, err := dao.db.Query(query, queryStr, limit)
	if err != nil {
		dao.l.Printf("Error searching challenges: %v", err)
		return nil, err
	}
	defer rows.Close()

	return dao.scanChallenges(rows)
}

func (dao *ChallengeDao) scanChallenges(rows *sql.Rows) ([]models.Challenge, error) {
	var challenges []models.Challenge
	for rows.Next() {
		var c models.Challenge
		err := rows.Scan(
			&c.ID, &c.Name, &c.Description, &c.ChallengeType, &c.CompetitionMode, &c.Visibility,
			&c.StartDate, &c.Deadline, &c.CreatedByUserID, &c.CreatedByGroupID,
			&c.TargetCount, &c.Region, &c.Difficulty, &c.IsFeatured, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			dao.l.Printf("Error scanning challenge: %v", err)
			return nil, err
		}
		challenges = append(challenges, c)
	}
	return challenges, nil
}

// ==================== Challenge Peaks ====================

func (dao *ChallengeDao) AddChallengePeak(challengeID int64, peakID int64, sortOrder int) error {
	query := `
		INSERT INTO challenge_peaks (challenge_id, peak_id, sort_order)
		VALUES ($1, $2, $3)
		ON CONFLICT (challenge_id, peak_id) DO UPDATE SET sort_order = $3;
	`
	_, err := dao.db.Exec(query, challengeID, peakID, sortOrder)
	if err != nil {
		dao.l.Printf("Error adding challenge peak: %v", err)
		return err
	}
	return nil
}

func (dao *ChallengeDao) RemoveChallengePeak(challengeID int64, peakID int64) error {
	query := `DELETE FROM challenge_peaks WHERE challenge_id = $1 AND peak_id = $2;`
	_, err := dao.db.Exec(query, challengeID, peakID)
	if err != nil {
		dao.l.Printf("Error removing challenge peak: %v", err)
		return err
	}
	return nil
}

func (dao *ChallengeDao) GetChallengePeaks(challengeID int64) ([]models.ChallengePeakWithDetails, error) {
	query := `
		SELECT
			cp.id, cp.challenge_id, cp.peak_id, cp.sort_order, cp.created_at,
			p.name, p.alt_name, p.latitude, p.longitude, p.elevation_meters, p.region
		FROM challenge_peaks cp
		JOIN peaks p ON cp.peak_id = p.id
		WHERE cp.challenge_id = $1
		ORDER BY cp.sort_order, p.name;
	`
	rows, err := dao.db.Query(query, challengeID)
	if err != nil {
		dao.l.Printf("Error getting challenge peaks: %v", err)
		return nil, err
	}
	defer rows.Close()

	var peaks []models.ChallengePeakWithDetails
	for rows.Next() {
		var p models.ChallengePeakWithDetails
		err := rows.Scan(
			&p.ID, &p.ChallengeID, &p.PeakID, &p.SortOrder, &p.CreatedAt,
			&p.Name, &p.AltName, &p.Latitude, &p.Longitude, &p.Elevation, &p.Region,
		)
		if err != nil {
			dao.l.Printf("Error scanning challenge peak: %v", err)
			return nil, err
		}
		peaks = append(peaks, p)
	}
	return peaks, nil
}

func (dao *ChallengeDao) GetChallengePeaksWithUserProgress(challengeID int64, userID int64) ([]models.ChallengePeakWithDetails, error) {
	query := `
		SELECT
			cp.id, cp.challenge_id, cp.peak_id, cp.sort_order, cp.created_at,
			p.name, p.alt_name, p.latitude, p.longitude, p.elevation_meters, p.region,
			EXISTS(
				SELECT 1 FROM challenge_summit_log csl 
				WHERE csl.challenge_id = cp.challenge_id 
				AND csl.peak_id = cp.peak_id 
				AND csl.user_id = $2
			) as is_summited
		FROM challenge_peaks cp
		JOIN peaks p ON cp.peak_id = p.id
		WHERE cp.challenge_id = $1
		ORDER BY cp.sort_order, p.name;
	`
	rows, err := dao.db.Query(query, challengeID, userID)
	if err != nil {
		dao.l.Printf("Error getting challenge peaks with progress: %v", err)
		return nil, err
	}
	defer rows.Close()

	var peaks []models.ChallengePeakWithDetails
	for rows.Next() {
		var p models.ChallengePeakWithDetails
		err := rows.Scan(
			&p.ID, &p.ChallengeID, &p.PeakID, &p.SortOrder, &p.CreatedAt,
			&p.Name, &p.AltName, &p.Latitude, &p.Longitude, &p.Elevation, &p.Region,
			&p.IsSummited,
		)
		if err != nil {
			dao.l.Printf("Error scanning challenge peak: %v", err)
			return nil, err
		}
		peaks = append(peaks, p)
	}
	return peaks, nil
}

func (dao *ChallengeDao) SetChallengePeaks(challengeID int64, peakIDs []int64) error {
	tx, err := dao.db.Begin()
	if err != nil {
		return err
	}

	// Delete existing peaks
	_, err = tx.Exec(`DELETE FROM challenge_peaks WHERE challenge_id = $1`, challengeID)
	if err != nil {
		tx.Rollback()
		dao.l.Printf("Error deleting existing challenge peaks: %v", err)
		return err
	}

	// Insert new peaks
	for i, peakID := range peakIDs {
		_, err = tx.Exec(
			`INSERT INTO challenge_peaks (challenge_id, peak_id, sort_order) VALUES ($1, $2, $3)`,
			challengeID, peakID, i,
		)
		if err != nil {
			tx.Rollback()
			dao.l.Printf("Error inserting challenge peak: %v", err)
			return err
		}
	}

	return tx.Commit()
}

// ==================== Participants ====================

func (dao *ChallengeDao) JoinChallenge(challengeID int64, userID int64) error {
	query := `
		INSERT INTO challenge_participants (challenge_id, user_id, total_peaks)
		VALUES ($1, $2, (SELECT COUNT(*) FROM challenge_peaks WHERE challenge_id = $1))
		ON CONFLICT (challenge_id, user_id) DO NOTHING;
	`
	_, err := dao.db.Exec(query, challengeID, userID)
	if err != nil {
		dao.l.Printf("Error joining challenge: %v", err)
		return err
	}
	return nil
}

func (dao *ChallengeDao) LeaveChallenge(challengeID int64, userID int64) error {
	query := `DELETE FROM challenge_participants WHERE challenge_id = $1 AND user_id = $2;`
	_, err := dao.db.Exec(query, challengeID, userID)
	if err != nil {
		dao.l.Printf("Error leaving challenge: %v", err)
		return err
	}
	return nil
}

func (dao *ChallengeDao) GetChallengeParticipants(challengeID int64) ([]models.ChallengeParticipantWithUser, error) {
	query := `
		SELECT
			cp.id, cp.challenge_id, cp.user_id, cp.joined_at, cp.completed_at,
			cp.peaks_completed, cp.total_peaks,
			COALESCE(u.username, '') as user_name
		FROM challenge_participants cp
		JOIN users u ON cp.user_id = u.id
		WHERE cp.challenge_id = $1
		ORDER BY cp.peaks_completed DESC, cp.joined_at;
	`
	rows, err := dao.db.Query(query, challengeID)
	if err != nil {
		dao.l.Printf("Error getting challenge participants: %v", err)
		return nil, err
	}
	defer rows.Close()

	var participants []models.ChallengeParticipantWithUser
	for rows.Next() {
		var p models.ChallengeParticipantWithUser
		err := rows.Scan(
			&p.ID, &p.ChallengeID, &p.UserID, &p.JoinedAt, &p.CompletedAt,
			&p.PeaksCompleted, &p.TotalPeaks,
			&p.UserName,
		)
		if err != nil {
			dao.l.Printf("Error scanning participant: %v", err)
			return nil, err
		}
		participants = append(participants, p)
	}
	return participants, nil
}

func (dao *ChallengeDao) GetChallengeLeaderboard(challengeID int64) ([]models.LeaderboardEntry, error) {
	query := `
		SELECT
			cp.user_id, COALESCE(u.username, '') as user_name,
			cp.peaks_completed, cp.total_peaks, cp.joined_at, cp.completed_at
		FROM challenge_participants cp
		JOIN users u ON cp.user_id = u.id
		WHERE cp.challenge_id = $1
		ORDER BY cp.peaks_completed DESC, cp.completed_at ASC NULLS LAST, cp.joined_at ASC;
	`
	rows, err := dao.db.Query(query, challengeID)
	if err != nil {
		dao.l.Printf("Error getting challenge leaderboard: %v", err)
		return nil, err
	}
	defer rows.Close()

	var leaderboard []models.LeaderboardEntry
	rank := 0
	prevPeaks := -1
	actualRank := 0

	for rows.Next() {
		var entry models.LeaderboardEntry
		err := rows.Scan(
			&entry.UserID, &entry.UserName,
			&entry.PeaksCompleted, &entry.TotalPeaks, &entry.JoinedAt, &entry.CompletedAt,
		)
		if err != nil {
			dao.l.Printf("Error scanning leaderboard entry: %v", err)
			return nil, err
		}

		actualRank++
		// Handle tied rankings - same peaks_completed = same rank
		if entry.PeaksCompleted != prevPeaks {
			rank = actualRank
			prevPeaks = entry.PeaksCompleted
		}
		entry.Rank = rank

		// Calculate progress percentage
		if entry.TotalPeaks > 0 {
			entry.Progress = float64(entry.PeaksCompleted) / float64(entry.TotalPeaks) * 100
		}

		leaderboard = append(leaderboard, entry)
	}
	return leaderboard, nil
}

func (dao *ChallengeDao) UpdateParticipantProgress(challengeID int64, userID int64, peaksCompleted int, totalPeaks int) error {
	query := `
		UPDATE challenge_participants
		SET peaks_completed = $3, total_peaks = $4
		WHERE challenge_id = $1 AND user_id = $2;
	`
	_, err := dao.db.Exec(query, challengeID, userID, peaksCompleted, totalPeaks)
	if err != nil {
		dao.l.Printf("Error updating participant progress: %v", err)
		return err
	}
	return nil
}

func (dao *ChallengeDao) MarkParticipantCompleted(challengeID int64, userID int64) error {
	query := `
		UPDATE challenge_participants
		SET completed_at = NOW()
		WHERE challenge_id = $1 AND user_id = $2 AND completed_at IS NULL;
	`
	_, err := dao.db.Exec(query, challengeID, userID)
	if err != nil {
		dao.l.Printf("Error marking participant completed: %v", err)
		return err
	}
	return nil
}

func (dao *ChallengeDao) IsUserParticipant(challengeID int64, userID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM challenge_participants WHERE challenge_id = $1 AND user_id = $2);`
	var exists bool
	err := dao.db.QueryRow(query, challengeID, userID).Scan(&exists)
	if err != nil {
		dao.l.Printf("Error checking user participation: %v", err)
		return false, err
	}
	return exists, nil
}

// ==================== Groups ====================

func (dao *ChallengeDao) AddGroupToChallenge(challengeID int64, groupID int64, deadlineOverride *time.Time) error {
	query := `
		INSERT INTO challenge_groups (challenge_id, group_id, deadline_override)
		VALUES ($1, $2, $3)
		ON CONFLICT (challenge_id, group_id) DO UPDATE SET deadline_override = $3;
	`
	_, err := dao.db.Exec(query, challengeID, groupID, deadlineOverride)
	if err != nil {
		dao.l.Printf("Error adding group to challenge: %v", err)
		return err
	}
	return nil
}

func (dao *ChallengeDao) RemoveGroupFromChallenge(challengeID int64, groupID int64) error {
	query := `DELETE FROM challenge_groups WHERE challenge_id = $1 AND group_id = $2;`
	_, err := dao.db.Exec(query, challengeID, groupID)
	if err != nil {
		dao.l.Printf("Error removing group from challenge: %v", err)
		return err
	}
	return nil
}

func (dao *ChallengeDao) GetChallengeGroups(challengeID int64) ([]models.ChallengeGroupWithDetails, error) {
	query := `
		SELECT
			cg.id, cg.challenge_id, cg.group_id, cg.started_at, cg.completed_at, cg.deadline_override,
			g.name as group_name,
			(SELECT COUNT(*) FROM group_members gm WHERE gm.group_id = g.id) as member_count
		FROM challenge_groups cg
		JOIN groups g ON cg.group_id = g.id
		WHERE cg.challenge_id = $1
		ORDER BY cg.started_at;
	`
	rows, err := dao.db.Query(query, challengeID)
	if err != nil {
		dao.l.Printf("Error getting challenge groups: %v", err)
		return nil, err
	}
	defer rows.Close()

	var groups []models.ChallengeGroupWithDetails
	for rows.Next() {
		var g models.ChallengeGroupWithDetails
		err := rows.Scan(
			&g.ID, &g.ChallengeID, &g.GroupID, &g.StartedAt, &g.CompletedAt, &g.DeadlineOverride,
			&g.GroupName, &g.MemberCount,
		)
		if err != nil {
			dao.l.Printf("Error scanning challenge group: %v", err)
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}

func (dao *ChallengeDao) GetGroupChallenges(groupID int64) ([]models.Challenge, error) {
	query := `
		SELECT
			c.id, c.name, c.description, c.challenge_type, c.competition_mode, c.visibility,
			c.start_date, c.deadline, c.created_by_user_id, c.created_by_group_id,
			c.target_count, c.region, c.difficulty, c.is_featured, c.created_at, c.updated_at
		FROM challenges c
		JOIN challenge_groups cg ON c.id = cg.challenge_id
		WHERE cg.group_id = $1
		ORDER BY c.created_at DESC;
	`
	rows, err := dao.db.Query(query, groupID)
	if err != nil {
		dao.l.Printf("Error getting group challenges: %v", err)
		return nil, err
	}
	defer rows.Close()

	return dao.scanChallenges(rows)
}

// ==================== Summit Log ====================

func (dao *ChallengeDao) LogSummit(logEntry models.ChallengeSummitLog) error {
	query := `
		INSERT INTO challenge_summit_log (challenge_id, user_id, peak_id, activity_id, summited_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (challenge_id, user_id, peak_id) DO NOTHING;
	`
	_, err := dao.db.Exec(query, logEntry.ChallengeID, logEntry.UserID, logEntry.PeakID, logEntry.ActivityID, logEntry.SummitedAt)
	if err != nil {
		dao.l.Printf("Error logging summit: %v", err)
		return err
	}
	return nil
}

func (dao *ChallengeDao) GetChallengeSummitLog(challengeID int64, userID *int64) ([]models.ChallengeSummitLogWithDetails, error) {
	query := `
		SELECT
			csl.id, csl.challenge_id, csl.user_id, csl.peak_id, csl.activity_id, csl.summited_at, csl.created_at,
			p.name as peak_name, p.elevation_meters as peak_elevation
		FROM challenge_summit_log csl
		LEFT JOIN peaks p ON csl.peak_id = p.id
		WHERE csl.challenge_id = $1
		AND ($2::bigint IS NULL OR csl.user_id = $2)
		ORDER BY csl.summited_at DESC;
	`
	rows, err := dao.db.Query(query, challengeID, userID)
	if err != nil {
		dao.l.Printf("Error getting challenge summit log: %v", err)
		return nil, err
	}
	defer rows.Close()

	var logs []models.ChallengeSummitLogWithDetails
	for rows.Next() {
		var l models.ChallengeSummitLogWithDetails
		err := rows.Scan(
			&l.ID, &l.ChallengeID, &l.UserID, &l.PeakID, &l.ActivityID, &l.SummitedAt, &l.CreatedAt,
			&l.PeakName, &l.PeakElevation,
		)
		if err != nil {
			dao.l.Printf("Error scanning summit log: %v", err)
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}

func (dao *ChallengeDao) HasUserSummitedPeakForChallenge(challengeID int64, userID int64, peakID int64) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM challenge_summit_log 
			WHERE challenge_id = $1 AND user_id = $2 AND peak_id = $3
		);
	`
	var exists bool
	err := dao.db.QueryRow(query, challengeID, userID, peakID).Scan(&exists)
	if err != nil {
		dao.l.Printf("Error checking summit: %v", err)
		return false, err
	}
	return exists, nil
}
