# Challenge System Fixes - Phase 2

## Status: In Progress

## Overview
Fix remaining issues with challenge progress tracking, participant display, activity logs, and summit detection.

## Completed
- ✅ Team progress now shows cumulative for collaborative challenges (backend sums all participants)
- ✅ Default tab set to participants
- ✅ Added StravaAthleteID to ChallengeParticipantWithUser model and GetChallengeParticipants query

## Backend Tasks

### 1. ✅ Update LeaderboardEntry Model
- ✅ Add `StravaAthleteID int64` field
- ✅ Add `TotalDistance float64` and `TotalElevation float64` fields
- ✅ Update `GetChallengeLeaderboard` query to fetch these fields
- ✅ Update scan logic to populate new fields

### 2. ✅ Add Activities Endpoint for Challenges
- ✅ Create `GetChallengeActivities` method in challengeDao
  - Query activities table for all participants within challenge date range
  - Filter by challenge start/end dates
  - Return activities sorted by start_date DESC
- ✅ Add to ChallengeService
- ✅ Expose via API handler

### 3. ✅ Fix Summit Detection for Specific Peaks
- ✅ Check `ProcessActivityForChallenges` is being called during activity processing
- ✅ Verify summit detection workflow includes challenge summit logging
- ✅ Added challengeService to SummitService constructor
- ✅ Modified `CalculateSummitsForActivity` to call `ProcessActivityForChallenges` when summit detected
- ✅ Updated server.go to pass challengeService to SummitService
- NOTE: Existing activities need to be reprocessed to credit summits to challenges (run summit detection again)

## Frontend Tasks

### 4. ✅ Update TypeScript Models
- ✅ Add `stravaAthleteId: number` to `ChallengeParticipantWithUser` interface
- ✅ Add `totalDistance`, `totalElevation`, `totalSummitCount` to `LeaderboardEntry` interface

### 5. ✅ Update Participant Display (Collaborative Challenges)
- ✅ Remove profile picture images
- ✅ Show username OR strava ID (if username blank)
- ✅ Add link to Strava profile: `https://www.strava.com/athletes/{stravaAthleteId}`
- ✅ For distance/elevation: show actual distance/elevation contributed, NOT progress %
- ✅ Template changes in challenge-detail.component.html

### 6. ✅ Update Leaderboard Display (Competitive Challenges)
- ✅ Similar to #5 but for leaderboard section
- ✅ Show distance/elevation values based on goal type
- ✅ Remove profile pictures
- ✅ Add Strava links

### 7. ✅ Add Activity Log Tab
- ✅ Create method to fetch challenge activities via API
- ✅ Display activities in activity tab with:
  - Activity name
  - Participant name (username or Strava ID)
  - Distance/elevation based on goal type
  - Date
  - Link to Strava activity
  - Link to Strava profile
- ✅ Template: challenge-detail.component.html activity tab section
- ✅ Created ActivityWithUser model to include user info

### 8. ✅ Add Challenge Start Date Display
- ✅ Show startDate in challenge detail header
- ✅ Format similar to deadline display
- ✅ Template: challenge-detail.component.html meta-info section

### 9. Add Map Click-Through for Peaks
- TODO: Add click handler to peak cards in peaks tab
- TODO: Navigate to `/explore` with challenge filter
- TODO: May need to add route params or state passing

## Current Progress
- ✅ Basic display updates complete
- ✅ Team progress for collaborative challenges working
- ✅ Strava links and ID fallbacks implemented
- ✅ Start date display added
- ✅ Activity log implemented with participant names
- ✅ Summit detection now credits summits to challenges
- Still need: Map click-through (optional feature)

## Testing Checklist
- [ ] Distance collaborative: team progress sums all participants
- [ ] Distance competitive: leaderboard shows individual distances
- [ ] Elevation collaborative: team progress sums all participants
- [ ] Elevation competitive: leaderboard shows individual elevations
- [ ] Summit count: counts summits from log
- [ ] Specific peaks: auto-detects summits and marks peaks complete
- [ ] Activity log shows all qualifying activities
- [ ] Participant display shows Strava ID when username blank
- [ ] Strava profile links work
- [ ] Start date displays correctly

## Implementation Notes
- Summit detection may require backend workflow investigation
- Activity log needs new DTO/endpoint
- Consider performance of summing participants for collaborative challenges (currently acceptable)
