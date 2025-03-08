package models

import "time"

type UserPeakJoin struct {
	PeakID     int64     // from up.peak_id
	UserID     int64     // from up.user_id
	ActivityID int64     // from up.activity_id
	SummitedAt time.Time // from up.summited_at
	UserName   int64     // from u.strava_athlete_id (as an example)
}
